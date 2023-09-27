package big

import (
    "bufio"
    "errors"
    "fmt"
    "math"
    "math/big"
    "os"
    "regexp"
    "strconv"
    "strings"
)

const (
    BigFloatPrec = 100
    pathFileName = "path.ini"
)

// Path
//
// Definiert eine Kamerafahrt ueber der komplexen Zahlenebene.
type Path struct {
    viewList []*View
    samp     int
}

// Erstellt eine neue (leere) Kamerafahrt.
func NewPath() *Path {
    var p *Path

    p = new(Path)
    p.viewList = make([]*View, 0)
    p.samp = 0

    return p
}

func (p *Path) SetSamples(samp int) {
    p.samp = samp
}

// Damit lassen sich die Stuetzstellen der Kamerafahrt aus einem File
// einlesen. Format des Files:
//
// Neues Format:
//
//    xm0 ym0 w0 it0    (x/y des Mittelpunktes, Breite des Bildes, max Iter)
//    xm1 ym1 w1 it1
//    ...
func ReadPath(pathName string) (*Path, error) {
    var p *Path
    var fd *os.File
    var scanner *bufio.Scanner
    var line string
    var matches []string
    var err error
    var x, y, w *big.Float
    var it int64
    var regComm, regBlock, regData *regexp.Regexp
    var inBlock bool

    x = big.NewFloat(0.0).SetPrec(100)
    y = big.NewFloat(0.0).SetPrec(100)
    w = big.NewFloat(0.0).SetPrec(100)

    regComm = regexp.MustCompile(`^ *(#.*)?$`)
    regBlock = regexp.MustCompile(`^ *\[([[:alnum:]]+)\] *$`)
    regData = regexp.MustCompile(`^ *([+-]?[0-9]+\.[0-9]+) +([+-]?[0-9]+\.[0-9]+) +([+-]?[0-9]+\.[0-9]+(?:[Ee]-?[0-9]+)?) +([0-9]+) *$`)

    fd, err = os.Open(pathFileName)
    if err != nil {
        return nil, err
    }
    p = NewPath()
    inBlock = false
    scanner = bufio.NewScanner(fd)
    for scanner.Scan() {
        line = scanner.Text()
        if regComm.MatchString(line) {
            continue
        }
        if inBlock {
            if regData.MatchString(line) {
                matches = regData.FindStringSubmatch(line)
                x.SetString(matches[1])
                y.SetString(matches[2])
                w.SetString(matches[3])
                it, _ = strconv.ParseInt(matches[4], 10, 32)
                p.AddView(x, y, w, int(it))
            } else if regBlock.MatchString(line) {
                break
            } else {
                return nil, errors.New(fmt.Sprintf("error on line: %s", line))
            }
        } else {
            if regBlock.MatchString(line) {
                matches = regBlock.FindStringSubmatch(line)
                if strings.Compare(matches[1], pathName) == 0 {
                    inBlock = true
                }
            }
        }
    }
    if !inBlock {
        return nil, errors.New(fmt.Sprintf("no path with name '%s' found!", pathName))
    }
    fd.Close()
    return p, nil
}

// Fuegt der Kamerafahrt eine neue Ansicht oder Stuetzstelle hinzu. Die neue
// Ansicht wird immer am Ende der bestehenden Kamerafahrt angehaengt.
func (p *Path) AddView(x, y, w *big.Float, it int) {
    var v *View

    v = NewView()
    v.SetValues(x, y, w, it, p.samp)
    p.viewList = append(p.viewList, v)
}

func (p *Path) NumViews() int {
    return len(p.viewList)
}

// Berechnet eine neue View auf dem Pfad zwischen der ersten und der letzten
// View. v muss eine initialisierte View sein. Der Parameter t ist ein Wert
// zwischen 0.0 und 1.0 und gibt die Position auf dem Pfad an.
func (p *Path) GetView(t float64, v *View) {
    var x, y, w *big.Float
    var v0, v1 *View

    x = big.NewFloat(0.0).SetPrec(100)
    y = big.NewFloat(0.0).SetPrec(100)
    w = big.NewFloat(0.0).SetPrec(100)

    if (t == 1.0) || (len(p.viewList) == 1) {
        i := len(p.viewList) - 1
        v.x.Set(p.viewList[i].x)
        v.y.Set(p.viewList[i].y)
        v.w.Set(p.viewList[i].w)
        v.it = p.viewList[i].it
    } else {
        i := int(t * float64(len(p.viewList)-1))
        v0 = p.viewList[i]
        v1 = p.viewList[i+1]
        t = t*float64(len(p.viewList)-1) - float64(i)
        tt := 0.5 * (1.0 - math.Cos(t*math.Pi))
        q := big.NewFloat(0.0).SetPrec(100)
        q.Quo(v1.w, v0.w)
        t0, _ := q.Float64()
        fw := math.Pow(t0, tt)
        w.Mul(v0.w, big.NewFloat(fw))
        it := v0.it + int(tt*float64(v1.it-v0.it))
        if v0.w.Cmp(v1.w) >= 0 {
            k := 1.0 - math.Pow(1.0-math.Pow(1.0-tt, 1.3), 1.0/1.3)
            t1 := fw * k
            x.Sub(v1.x, v0.x)
            x.Mul(x, big.NewFloat(t1))
            x.Sub(v1.x, x)
            y.Sub(v1.y, v0.y)
            y.Mul(y, big.NewFloat(t1))
            y.Sub(v1.y, y)
        } else {
            k := 1.0 - math.Pow(1.0-math.Pow(tt, 1.3), 1.0/1.3)
            fw = fw / t0
            t1 := fw * k
            x.Sub(v0.x, v1.x)
            x.Mul(x, big.NewFloat(t1))
            x.Sub(v0.x, x)
            y.Sub(v0.y, v1.y)
            y.Mul(y, big.NewFloat(t1))
            y.Sub(v0.y, y)
        }
        v.x.Set(x)
        v.y.Set(y)
        v.w.Set(w)
        v.it = it
    }
    v.samp = p.samp
}
