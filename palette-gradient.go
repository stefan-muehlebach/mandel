package mandel

import (
    "container/list"
    "image/color"
    "regexp"
    "strconv"
)

// Dieser Typ realisiert eine Palette mit interpolierten Farbverläufen.
// Jede Farbe kann einzeln eingestellt werden kann.
var (
    gradientPalRegexp = regexp.MustCompile(`^ *([0-9\.]+) *: *([0-9\.]+|-) +([0-9\.]+|-) +([0-9\.]+|-) *$`)
)

type GradientPalette struct {
    basePalette
    pointList []*list.List
    intFunc   InterpFunc
    rexp      *regexp.Regexp
}

type GradPoint struct {
    Pos, Val float64
}

func NewGradientPalette() *GradientPalette {
    p := &GradientPalette{}
    p.Init()
    p.pointList = make([]*list.List, NumBaseColors)
    for i := Red; i < NumBaseColors; i++ {
        p.pointList[i] = list.New()
    }
    p.intFunc = LinInterpFunc
    p.intFunc = CubicInterpFunc
    p.rexp = gradientPalRegexp
    return p
}

func (p *GradientPalette) GetRegexp() *regexp.Regexp {
    return p.rexp
}

func (p *GradientPalette) ProcessLine(line string) error {
    var matches []string
    var t, v float64

    matches = p.rexp.FindStringSubmatch(line)
    t, _ = strconv.ParseFloat(matches[1], 32)
    for i := Red; i < NumBaseColors; i++ {
        if matches[i+2] == "-" {
            continue
        }
        v, _ = strconv.ParseFloat(matches[i+2], 32)
        gp := &GradPoint{t, v}
        if err := p.AddGradPoint(i, gp); err != nil {
            return err
        }
    }
    return nil
}

func (p *GradientPalette) Update() {
    var colVal [3]uint8
    var e *list.Element
    var gradPt, nextGradPt *GradPoint

    for i := 0; i < len(p.colorList); i++ {
        f := float64(i) / float64(len(p.colorList)-1)

        for j := Red; j < NumBaseColors; j++ {
            l := p.pointList[j]
            for e = l.Back(); e != nil; e = e.Prev() {
                if f >= e.Value.(*GradPoint).Pos {
                    break
                }
            }
            gradPt = e.Value.(*GradPoint)
            if f == gradPt.Pos {
                colVal[j] = uint8(255.0 * gradPt.Val)
            } else {
                nextGradPt = e.Next().Value.(*GradPoint)
                t := (f - gradPt.Pos) / (nextGradPt.Pos - gradPt.Pos)
                colVal[j] = uint8(255.0 * InterpValue(gradPt.Val, nextGradPt.Val, t, p.intFunc))
            }
        }
        p.colorList[i] = color.RGBA{colVal[0], colVal[1], colVal[2], 0xff}
    }
}

func (p *GradientPalette) Ready() bool {
    for i := Red; i < NumBaseColors; i++ {
        l := p.pointList[i]
        e := l.Front()
        if e.Value.(*GradPoint).Pos != 0.0 {
            return false
        }
        e = l.Back()
        if e.Value.(*GradPoint).Pos != 1.0 {
            return false
        }
    }
    return true
}

func (p *GradientPalette) AddGradPoint(col BaseColorType, gp *GradPoint) error {
    var e *list.Element
    l := p.pointList[col]

    if l.Len() == 0 {
        l.PushFront(gp)
        return nil
    }
    for e = l.Front(); e != nil; e = e.Next() {
        if gp.Pos < e.Value.(*GradPoint).Pos {
            break
        }
    }
    if e == nil {
        l.PushBack(gp)
    } else {
        l.InsertBefore(gp, e)
    }
    return nil
}

func (p *GradientPalette) DelGradPoint(col BaseColorType, gp *GradPoint) {
    l := p.pointList[col]
    for e := l.Front(); e != nil; e = e.Next() {
        if e.Value == gp {
            l.Remove(e)
            return
        }
    }
}

func (p *GradientPalette) SetGradPoint(col BaseColorType, j int, gp *GradPoint) {
    l := p.pointList[col]
    e := l.Front()
    for ; j > 0; j-- {
        e = e.Next()
    }
    e.Value = gp
}

/*
func (p *GradientPalette) FirstGradPoint(col BaseColorType) (*GradPoint) {
    l := p.pointList[col]
    e := l.Front()
    return e.Value.(*GradPoint)
}

func (p *GradientPalette) LastGradPoint(col BaseColorType) (*GradPoint) {
    l := p.pointList[col]
    e := l.Back()
    return e.Value.(*GradPoint)
}
*/

func (p *GradientPalette) GradPoint(col BaseColorType, j int) *GradPoint {
    l := p.pointList[col]
    e := l.Front()
    for ; j > 0; j-- {
        e = e.Next()
    }
    return e.Value.(*GradPoint)
}

func (p *GradientPalette) GradPointList(col BaseColorType) []*GradPoint {
    l := p.pointList[col]
    gpl := make([]*GradPoint, 0)
    for e := l.Front(); e != nil; e = e.Next() {
        gpl = append(gpl, e.Value.(*GradPoint))
    }
    return gpl
}

func (p *GradientPalette) NumGradPoints(col BaseColorType) (int) {
    return p.pointList[col].Len()
}

