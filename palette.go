package mandel

import (
    "bufio"
    _ "errors"
    "fmt"
    "image/color"
    _ "log"
    "math"
    "os"
    "regexp"
    "strings"
)

const (
    // palSize ist die Groesse der Farbpalette (Anzahl RGBA-Werte). Je groesser
    // dieser Wert, desto feiner die Aufteilung der Farben. Werte ueber 1000
    // machen jedoch in meinen Augen keinen Sinn mehr.
    palSize = 512

    // palFileName ist der Name der Datei, welche alle Farbpaletten enthaelt.
    // Diese Datei muss im aktuellen Verzeichnis zu finden sein.
    palFileName = "palette.ini"
)

var (
    // regxComm ist der regulaere Ausdruck, welche beim Lesen der
    // Konfigurationsdatei mit den Paletten fuer leere Zeilen und
    // ommentarzeilen
    // verwendet wird.
    //
    regxComm = regexp.MustCompile(`^ *(#.*)?$`)

    // regxSection dagegen ist der regulaere Ausdruck fuer die Erkennung der
    // Abschnitts-Titel der einzelnen Paletten
    //
    regxSection = regexp.MustCompile(`^ *\[([[:alnum:]]+)\] *$`)
)

//-----------------------------------------------------------------------------

type BaseColorType int

const (
    Red BaseColorType = iota
    Green
    Blue
    NumBaseColors
)

func (c BaseColorType) String() string {
    switch c {
    case Red:
        return "Red"
    case Green:
        return "Green"
    case Blue:
        return "Blue"
    default:
        return "(Unknown color)"
    }
}

//-----------------------------------------------------------------------------

// ColorList ist ein Hilfstyp (nicht exportiert), mit dem ein Slice von
// Farben verwaltet werden kann.
type ColorList []color.RGBA

// getColor dient der Ermittlung einer Farbe, welche im Array zwischen
// cl[i] und cl[i+1] liegt. Der Parameter t im Intervall [0,1) wird fuer
// eine lineare Interpolation zwischen den beiden Farben verwendet.
//
//func (cl ColorList) getColor(i int, t float64) color.RGBA {
//    return cl.getInterpColor(i, t, LinInterpFunc)
//}

// getInterpColor ist analog zu getColor, aber hier kann die Interpolations-
// funktion mit dem Parameter fnc bestimmt werden.
func (cl ColorList) InterpColor(i int, t float64) (col color.RGBA) {
    c0 := cl[i]
    c1 := cl[i+1]
    r0, g0, b0 := float64(c0.R), float64(c0.G), float64(c0.B)
    r1, g1, b1 := float64(c1.R), float64(c1.G), float64(c1.B)
    t0 := 1.0 - t
    col.R = uint8(t0*r0 + t*r1)
    col.G = uint8(t0*g0 + t*g1)
    col.B = uint8(t0*b0 + t*b1)
    col.A = 0xff
    return
}

type ColorFunc func(t float64) color.RGBA

// -----------------------------------------------------------------------------
//
// Ein neuer, schlankerer Versuch mit interpolierten Werten aus einem Array
// von float64-Werten.
type ValueList []float64

// Retourniert einen interpolierten Wert, der zwischen vl[i] und vl[i+1] liegt
// wobei mit t (in [0,1]) der 'Anteil' von vl[i] und vl[i+1] gewaehlt werden
// kann. fnc schliesslich ist die Interpolationsfunktion (mehr dazu: siehe
// weiter unten).
func (vl ValueList) InterpValue(i int, t float64, fnc InterpFunc) float64 {
    t = fnc(t)
    return (1.0-t)*vl[i] + t*vl[i+1]
}

func InterpValue(v1, v2, t float64, fnc InterpFunc) float64 {
    t = fnc(t)
    return (1.0-t)*v1 + t*v2
}

// -----------------------------------------------------------------------------
//
// Mit Hilfe dieses Funktionstyps InterpFunc koennen verschiedene
// Interpolationsfunktionen verwendet werden. Das Argument ist ein Wert
// in [0, 1], das Resultat ein Wert in [0, 1]. Es muss gelten:
//
//   - InterpFunc(0) = 0
//   - InterpFunc(1) = 1
//   - (a >= b) => (InterpFunc(a) >= InterpFunc(b))
type InterpFunc func(t float64) float64

// Die lineare Interpolation
func LinInterpFunc(t float64) float64 {
    return t
}

// Die kubische Interpolation, welche weichere Uebergaenge realisiert.
func CubicInterpFunc(t float64) float64 {
    return 3.0*t*t - 2.0*t*t*t
}

// Retourniert einen Slice mit den Namen aller verfuegbaren Paletten.
func PaletteNames() ([]string, error) {
    var fd *os.File
    var scanner *bufio.Scanner
    var line string
    var matches []string
    var err error
    var palNames []string

    palNames = make([]string, 0)
    fd, err = OpenConfFile(palFileName)
    if err != nil {
        return nil, err
    }
    defer fd.Close()
    scanner = bufio.NewScanner(fd)
    for scanner.Scan() {
        line = scanner.Text()
        if regxSection.MatchString(line) {
            matches = regxSection.FindStringSubmatch(line)
            palNames = append(palNames, matches[1])
        } else {
            continue
        }
    }
    return palNames, nil
}

// Dieser Datentyp repraesentiert eine Palette, wobei die Implementation
// offen, resp. beliebig ist. Dieser Datentyp enthaelt noch einige Felder,
// welche fuer Farbpaletten typisch sind, die in Zusammenhang mit den Mandel-
// brotmengen verwendet werden.
type basePalette struct {
    colorList    ColorList
    len          int
    lenIsMaxIter bool
    offset       float64
}

// Initialisiert die Felder des Basistyps einer Palette.
func (p *basePalette) Init() {
    p.colorList = make([]color.RGBA, palSize)
    p.lenIsMaxIter = true
    p.offset = 0.0
}

// Erstellt eine neue Palette aufgrund des Paletten-Namens in palName.
func NewPalette(palName string) (Palette, error) {
    var p Palette
    var fd *os.File
    var scanner *bufio.Scanner
    var line string
    var matches []string
    var err error
    var inSection bool

    fd, err = OpenConfFile(palFileName)
    if err != nil {
        return nil, err
    }
    defer fd.Close()
    inSection = false
    scanner = bufio.NewScanner(fd)
    for scanner.Scan() {
        line = scanner.Text()
        if regxComm.MatchString(line) {
            continue
        }
        if inSection {
            if p == nil {
                if gradientPalRegexp.MatchString(line) {
                    p = NewGradientPalette()
                } else if procPalRegexp.MatchString(line) {
                    p = NewProcPalette()
                } else {
                    return nil, fmt.Errorf("error on line: %s", line)
                }
            }
            if p.GetRegexp().MatchString(line) {
                if err := p.ProcessLine(line); err != nil {
                    return nil, fmt.Errorf("error on line: '%s':\n%v", line, err)
                }
            } else if regxSection.MatchString(line) {
                break
            } else {
                return nil, fmt.Errorf("error on line: %s", line)
            }
        } else {
            if regxSection.MatchString(line) {
                matches = regxSection.FindStringSubmatch(line)
                if strings.Compare(matches[1], palName) == 0 {
                    inSection = true
                }
            }
        }
    }
    if !inSection {
        return nil, fmt.Errorf("no palette '%s' found!", palName)
    }
    if !p.Ready() {
        return nil, fmt.Errorf("some values are missing for this palette")
    }
    p.Update()
    return p, nil
}

// Setzt die Laenge der Palette auf den Wert len.
func (p *basePalette) SetLength(len int) {
    p.len = len
}

// Retourniert die Laenge der Palette.
func (p *basePalette) Length() int {
    return p.len
}

// Hinterlegt, dass die Laenge der Palette gleich der maximalen Anzahl
// Iterationen sein soll.
func (p *basePalette) LenIsMaxIter() {
    p.lenIsMaxIter = true
}

// Hinterlegt, dass die Laenge der Palette nichts mit der maximalen Anzahl
// Iterationen zu tun haben soll.
func (p *basePalette) LenIsNotMaxIter() {
    p.lenIsMaxIter = false
}

// Holt die Information aus der Palette, ob die Palettenlaenge mit der max.
// Anzahl Iterationen verbunden sein soll.
func (p *basePalette) IsLenMaxIter() bool {
    return p.lenIsMaxIter
}

// Setzt den Offset fuer die erste Farbe der Palette.
func (p *basePalette) SetOffset(offset float64) {
    p.offset = offset
}

// Retourniert den Offset fuer die erste Farbe der Palette.
func (p *basePalette) Offset() float64 {
    return p.offset
}

// Mit GetColor kann eine Farbe aus der Farbpalette ermittelt werden.
// f ist eine beliebige Zahl (>= 0.0), welche zusammen mit der hinterlegten,
// fiktiven Palettenlaenge p.len und dem definierten Offset p.offset verwendet
// wird, den Index der gesuchten Farbe zu bestimmen. Ist f < 0.0, dann wird
// als Farbe IMMER Schwarz retourniert.
func (p *basePalette) GetColor(f float64) color.RGBA {
    var m, d, i, t float64

    if f < 0.0 {
        return color.RGBA{0, 0, 0, 0xff}
    }
    f += p.offset * float64(p.len)
    m = math.Mod(f, float64(p.len))
    d = float64(palSize-1) * (m / float64(p.len))
    i, t = math.Modf(d)
    return p.colorList.InterpColor(int(i), t)
}
