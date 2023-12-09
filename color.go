package mandel

import (
    _ "github.com/stefan-muehlebach/gg/color"
)

type RGBA struct {
    R, G, B, A uint32
}

func (c RGBA) RGBA() (r, g, b, a uint32) {
    return uint32(c.R), uint32(c.G), uint32(c.B), uint32(c.A)
}

func (c RGBA) AddRGBA(c2 RGBA) {
    c.R += c2.R
    c.G += c2.G
    c.B += c2.B
    c.A += c2.A
}

func (c RGBA) DivRGBA(d int) {
    c.R /= uint32(d)
    c.G /= uint32(d)
    c.B /= uint32(d)
    c.A /= uint32(d)
}

