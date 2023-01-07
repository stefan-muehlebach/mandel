package mandel

import (
	"bufio"
	"errors"
	"fmt"
    "math"
    "image/color"
    "log"
    "os"
	"regexp"
	"strconv"
    "strings"
)

//----------------------------------------------------------------------------
//
// Palette --
//
type Palette struct {
	colorList []color.RGBA
	palette []color.RGBA
}

func NewPalette() (*Palette) {
    var p *Palette
    p = new(Palette)
    return p
}

// Liest die Farbwerte der Palette aus einem File ein. Das Format des Files
// ist wie folgt definiert:
//
//     [PAL_NAME]
//     0xab 0x00 0xff
//     0xcd 0xff 0x00
//     0xaa 0xee 0xa0
//     ...
//
func (p *Palette) ReadFile(fileName, palName string) (error) {
    var fd *os.File
	var scanner *bufio.Scanner
	var line string
	var matches []string
    var err error
    var r, g, b uint64
	var regComm, regBlock, regData *regexp.Regexp
	var inBlock bool
	
	regComm  = regexp.MustCompile(`^ *(#.*)?$`)
	regBlock = regexp.MustCompile(`^ *\[([[:alnum:]]+)\] *$`)
	regData  = regexp.MustCompile(`^ *0x([0-9a-f]+) +0x([0-9a-f]+) +0x([0-9a-f]+) *$`)

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
				r, _ = strconv.ParseUint(matches[1], 16, 8)
				g, _ = strconv.ParseUint(matches[2], 16, 8)
				b, _ = strconv.ParseUint(matches[3], 16, 8)
				
				p.colorList = append(p.colorList,
						color.RGBA{uint8(r), uint8(g), uint8(b), 0xff})
			} else if regBlock.MatchString(line) {
				break
			} else {
				return errors.New(fmt.Sprintf("error on line: %s", line))
			}
		} else {
			if regBlock.MatchString(line) {
				matches = regBlock.FindStringSubmatch(line)
				if strings.Compare(matches[1], palName) == 0 {
					inBlock = true
				}
			}
		}
    }
	if ! inBlock {
		return errors.New(fmt.Sprintf("no palette with name '%s' found!", palName))
	}
	if len(p.colorList) < 2 {
		return errors.New("you must specify at least two colors!")
	}
    p.recalc()
    fd.Close()
	return nil
}

// Retourniert einen interpolierten Farbwert fuer den Iterationswert t.
// 0.0 <= t < 1.0.
//
func (p *Palette) GetColor(t float64) (color.RGBA) {
    var c0, c1, c2 color.RGBA
	
	if (t < 0.0) || (t > 1.0) {
		log.Fatal("t must be in the interval [0,1)!")
	}
	
	if t == 1.0 {
		return color.RGBA{0, 0, 0, 0xff}
	}

	i, d := math.Modf(t * float64(len(p.palette)-1))
	c0 = p.palette[int(i)]
	c1 = p.palette[int(i) + 1]

    r0 := float64(c0.R)
    g0 := float64(c0.G)
    b0 := float64(c0.B)
    r1 := float64(c1.R)
    g1 := float64(c1.G)
    b1 := float64(c1.B)
    c2.R = uint8(r0 + d * (r1-r0))
    c2.G = uint8(g0 + d * (g1-g0))
    c2.B = uint8(b0 + d * (b1-b0))
    c2.A = c0.A
    return c2
}

func (p *Palette) recalc() {
    var c0, c1 color.RGBA

	p.palette = make([]color.RGBA, 256 * (len(p.colorList)-1))
	for i := 0; i < len(p.colorList)-1; i++ {
		c0 = p.colorList[i]
		c1 = p.colorList[i+1]
		r0 := float64(c0.R)
        g0 := float64(c0.G)
        b0 := float64(c0.B)
        r1 := float64(c1.R)
        g1 := float64(c1.G)
        b1 := float64(c1.B)
        for j := 0; j < 256; j++ {
            d := float64(j) / float64(255)
            p.palette[i*256+j].R = uint8((1.0 - d) * r0 + d * r1)
            p.palette[i*256+j].G = uint8((1.0 - d) * g0 + d * g1)
            p.palette[i*256+j].B = uint8((1.0 - d) * b0 + d * b1)
            p.palette[i*256+j].A = 0xff
        }
    }
}

