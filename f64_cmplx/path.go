package f64_cmplx

import (
	"bufio"
	"errors"
    "fmt"
    "math"
    "os"
	"regexp"
	"strconv"
    "strings"
)

// Path
//
// Definiert eine Kamerafahrt ueber der komplexen Zahlenebene.
//
type Path struct {
    viewList []*View
}

// Erstellt eine neue (leere) Kamerafahrt.
//
func NewPath() (*Path) {
    var p *Path

    p = new(Path)
    p.viewList = make([]*View, 0)

    return p
}

// Damit lassen sich die Stuetzstellen der Kamerafahrt aus einem File
// einlesen. Format des Files:
//
// Neues Format:
//
//     xm0 ym0 w0 it0    (x/y des Mittelpunktes, Breite des Bildes, max Iter)
//     xm1 ym1 w1 it1
//     ...
//
func (p *Path) ReadFile(fileName, pathName string) (error) {
    var fd *os.File
	var scanner *bufio.Scanner
	var line string
	var matches []string
    var err error
    var x, y, w float64
    var it int64
	var regComm, regBlock, regData *regexp.Regexp
	var inBlock bool
	
	regComm  = regexp.MustCompile(`^ *(#.*)?$`)
	regBlock = regexp.MustCompile(`^ *\[([[:alnum:]]+)\] *$`)
	regData  = regexp.MustCompile(`^ *([+-]?[0-9]+\.[0-9]+) +([+-]?[0-9]+\.[0-9]+) +([+-]?[0-9]+\.[0-9]+(?:[Ee]-?[0-9]+)?) +([0-9]+) *$`)

    fd, err = os.Open(fileName)
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
				x, _  = strconv.ParseFloat(matches[1], 64)
				y, _  = strconv.ParseFloat(matches[2], 64)
				w, _  = strconv.ParseFloat(matches[3], 64)
				it, _ = strconv.ParseInt(matches[4], 10, 32)
				p.AddView(complex(x, y), w, int(it))
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
	if ! inBlock {
		return errors.New(fmt.Sprintf("no path with name '%s' found!", pathName))
	}
    fd.Close()
	return nil
}

// Fuegt der Kamerafahrt eine neue Ansicht oder Stuetzstelle hinzu. Die neue
// Ansicht wird immer am Ende der bestehenden Kamerafahrt angehaengt.
//
func (p *Path) AddView(z complex128, w float64, it int) {
	var v *View
	
	v = &View{z, w, it}
	p.viewList = append(p.viewList, v)
}

func (p *Path) NumViews() (int) {
    return len(p.viewList)
}

// Berechnet eine neue View auf dem Pfad zwischen der ersten und der letzten
// View. v muss eine initialisierte View sein. Der Parameter t ist ein Wert
// zwischen 0.0 und 1.0 und gibt die Position auf dem Pfad an.
//
func (p *Path) GetView(t float64, v *View) {
	var z complex128
    var v0, v1 *View

	if (t == 1.0) || (len(p.viewList) == 1) {
		i := len(p.viewList)-1
		v.z = p.viewList[i].z
		v.w = p.viewList[i].w
		v.it = p.viewList[i].it
	} else {
		i := int(t * float64(len(p.viewList) - 1))
		v0  = p.viewList[i]
		v1  = p.viewList[i+1]
		t   = t * float64(len(p.viewList) - 1) - float64(i)
		tt := 0.5 * (1.0 - math.Cos(t * math.Pi))
		fw := math.Pow(v1.w / v0.w, tt)
		w  := v0.w * fw
		it := v0.it + int(tt * float64(v1.it - v0.it))
		if v0.w >= v1.w {
			k := 1.0 - math.Pow(1.0 - math.Pow(1.0 - tt, 1.3), 1.0/1.3)
			z  = v1.z - complex(fw*k, 0.0) * (v1.z - v0.z)
		} else {
			k := 1.0 - math.Pow(1.0 - math.Pow(tt, 1.3), 1.0/1.3)
			fw = fw / (v1.w / v0.w)
			z  = v0.z - complex(fw*k, 0.0) * (v0.z - v1.z)
		}
		v.z = z
		v.w = w
		v.it = it
	}
}
