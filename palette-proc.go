package mandel

import (
	"image/color"
	"math"
	"regexp"
	"strconv"
)

var (
	procPalRegexp = regexp.MustCompile(`^ *([rgb]) *: *([-+0-9\.]+) +([-+0-9\.]+) +([-+0-9\.]+) +([-+0-9\.]+) *$`)
	colorMap      = map[byte]BaseColorType{'r': Red, 'g': Green, 'b': Blue}
)

type ProcParamType int

const (
	ParA ProcParamType = iota
	ParB
	ParC
	ParD
	NumProcParams
)

type ProcPalette struct {
	basePalette
	v    [][]float64
	rexp *regexp.Regexp
}

func NewProcPalette() *ProcPalette {
	p := &ProcPalette{}
	p.Init()
	p.v = make([][]float64, NumBaseColors)
	p.rexp = procPalRegexp
	return p
}

func (p *ProcPalette) GetRegexp() *regexp.Regexp {
	return p.rexp
}

func (p *ProcPalette) ProcessLine(line string) error {
	var matches []string

	matches = p.rexp.FindStringSubmatch(line)
	colIdx := colorMap[matches[1][0]]
	a, _ := strconv.ParseFloat(matches[2], 32)
	b, _ := strconv.ParseFloat(matches[3], 32)
	c, _ := strconv.ParseFloat(matches[4], 32)
	d, _ := strconv.ParseFloat(matches[5], 32)
	p.SetParamList(colIdx, []float64{a, b, c, d})
	return nil
}

func (p *ProcPalette) Update() {
	for i := 0; i < len(p.colorList); i++ {
		f := float64(i) / float64(len(p.colorList)-1)
		r := 255.0 * p.Value(Red, f)
		g := 255.0 * p.Value(Green, f)
		b := 255.0 * p.Value(Blue, f)
		p.colorList[i] = color.RGBA{uint8(r), uint8(g), uint8(b), 0xff}
	}
}

func (p *ProcPalette) Ready() bool {
	return true
}

func (p *ProcPalette) SetParamList(col BaseColorType, v []float64) {
	p.v[col] = v
}

func (p *ProcPalette) ParamList(col BaseColorType) []float64 {
	return p.v[col]
}

func (p *ProcPalette) SetParam(col BaseColorType, par ProcParamType, v float64) {
	p.v[col][par] = v
}

func (p *ProcPalette) Param(col BaseColorType, par ProcParamType) float64 {
	return p.v[col][par]
}

func (p *ProcPalette) Value(col BaseColorType, f float64) float64 {
	v := p.v[col]
	return v[0] + v[1]*math.Cos(2*math.Pi*(v[2]*f+v[3]))
}

