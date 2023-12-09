package f64

import (
    "bufio"
    "errors"
    "fmt"
    "math"
    "os"
    "regexp"
    "strconv"
    "strings"
    . "github.com/stefan-muehlebach/mandel"
)

const (
    pathFileName = "path.ini"
)

// Der Datentyp Path definiert eine Kamerafahrt ueber der komplexen
// Zahlenebene. Eine solche Fahrt besteht aus mehreren Views, welche festlegen,
// wie gross der gezeigte Ausschnitt der komplexen Ebene ist und wieviele
// Iterationen bei der Berechnung maximal verwendet werden sollen.
type f64Path struct {
    viewList []View
}

// Erstellt einen neuen Pfad, der noch keine Ansichten hat.
func NewPath() *f64Path {
    p := &f64Path{}
    p.viewList = make([]View, 0)
    return p
}

// Damit lassen sich die Stuetzstellen der Kamerafahrt aus einem File
// einlesen. Format des Files:
//
// Neues Format:
//
//    xm0 ym0 w0 it0    (x/y des Mittelpunktes, Breite des Bildes, max Iter)
//    xm1 ym1 w1 it1
//    ...
func (p *f64Path) Read(pathName string) (error) {
    var fd *os.File
    var scanner *bufio.Scanner
    var line string
    var matches []string
    var err error
    var x, y, w float64
    var it int64
    var regComm, regBlock, regData *regexp.Regexp
    var inBlock bool

    regComm = regexp.MustCompile(`^ *(#.*)?$`)
    regBlock = regexp.MustCompile(`^ *\[([[:alnum:]]+)\] *$`)
    regData = regexp.MustCompile(`^ *([+-]?[0-9]+\.[0-9]+) +([+-]?[0-9]+\.[0-9]+) +([+-]?[0-9]+\.[0-9]+(?:[Ee]-?[0-9]+)?) +([0-9]+) *$`)

    fd, err = OpenConfFile(pathFileName)
    if err != nil {
        return err
    }
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
                x, _ = strconv.ParseFloat(matches[1], 64)
                y, _ = strconv.ParseFloat(matches[2], 64)
                w, _ = strconv.ParseFloat(matches[3], 64)
                it, _ = strconv.ParseInt(matches[4], 10, 32)
                p.AddView(x, y, w, int(it))
            } else if regBlock.MatchString(line) {
                break
            } else {
                return errors.New(fmt.Sprintf("error on line: %s", line))
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
        return errors.New(fmt.Sprintf("no path with name '%s' found!", pathName))
    }
    fd.Close()
    return nil
}

// Fuegt der Kamerafahrt eine neue Ansicht oder Stuetzstelle hinzu. Die neue
// Ansicht wird immer am Ende der bestehenden Kamerafahrt angehaengt.
func (p *f64Path) AddView(x, y, w float64, it int) {
    v := NewView()
    v.SetValues(x, y, w, it)
    p.viewList = append(p.viewList, v)
}

// Mit NumViews wird die Anzahl der Ansichten in diesem Pfad ermittelt.
func (p *f64Path) NumViews() int {
    return len(p.viewList)
}

// Berechnet eine neue View auf dem Pfad zwischen der ersten und der letzten
// View. v muss eine initialisierte View sein. Der Parameter t ist ein Wert
// zwischen 0.0 und 1.0 und gibt die Position auf dem Pfad an.
// func (p *f64Path) GetViewOld(t float64, v *View) {
//     var x, y float64
//     var v0, v1 *View

//     if (t == 1.0) || (p.NumViews() == 1) {
//         i :=  p.NumViews() - 1
//         v.x = p.viewList[i].x
//         v.y = p.viewList[i].y
//         v.w = p.viewList[i].w
//         v.it = p.viewList[i].it
//     } else {
//         // TO DO: das ist noch das extrem komplizierte Interpolationsverfahren,
//         // welches mit grosser W'keit einfacher implementiert werden koennte
//         // (siehe Interpolation bei der Palette).
//         i := int(t * float64(len(p.viewList)-1))
//         v0 = p.viewList[i]
//         v1 = p.viewList[i+1]
//         t = t*float64(len(p.viewList)-1) - float64(i)
//         tt := 0.5 * (1.0 - math.Cos(t*math.Pi))
//         fw := math.Pow(v1.w/v0.w, tt)
//         w := v0.w * fw
//         it := v0.it + int(tt*float64(v1.it-v0.it))
//         if v0.w >= v1.w {
//             k := 1.0 - math.Pow(1.0-math.Pow(1.0-tt, 1.3), 1.0/1.3)
//             x = v1.x - fw*(v1.x-v0.x)*k
//             y = v1.y - fw*(v1.y-v0.y)*k
//         } else {
//             k := 1.0 - math.Pow(1.0-math.Pow(tt, 1.3), 1.0/1.3)
//             fw = fw / (v1.w / v0.w)
//             x = v0.x - fw*(v0.x-v1.x)*k
//             y = v0.y - fw*(v0.y-v1.y)*k
//         }
//         v.x = x
//         v.y = y
//         v.w = w
//         v.it = it
//     }
// }

func (p *f64Path) GetView(t float64) (v View) {
    // var x, y float64
    // var v0, v1 View

    if (t == 1.0) || (p.NumViews() == 1) {
        i := p.NumViews() - 1
        v = p.viewList[i]
    } else {
        v = NewView()
        // TO DO: das ist noch das extrem komplizierte Interpolationsverfahren,
        // welches mit grosser W'keit einfacher implementiert werden koennte
        // (siehe Interpolation bei der Palette).
        i := int(t * float64(p.NumViews()-1))
        // v0 = p.viewList[i]
        x0, y0, w0, it0 := p.viewList[i].Values()
        // v1 = p.viewList[i+1]
        x1, y1, w1, it1 := p.viewList[i+1].Values()
        
        t = t*float64(p.NumViews()-1) - float64(i)
        tt := 0.5 * (1.0 - math.Cos(t*math.Pi))
        fw := math.Pow(w1/w0, tt)
        w := w0 * fw
        it := it0 + int(tt*float64(it1-it0))
        x, y := 0.0, 0.0
        if w0 >= w1 {
            k := 1.0 - math.Pow(1.0-math.Pow(1.0-tt, 1.3), 1.0/1.3)
            x = x1 - fw*(x1-x0)*k
            y = y1 - fw*(y1-y0)*k
        } else {
            k := 1.0 - math.Pow(1.0-math.Pow(tt, 1.3), 1.0/1.3)
            fw = fw / (w1 / w0)
            x = x0 - fw*(x0-x1)*k
            y = y0 - fw*(y0-y1)*k
        }
        v.SetValues(x, y, w, it)
    }
    return
}
