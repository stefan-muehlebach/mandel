//
// Copyright 2018 by Stefan Muehlebach
//

// Das Package mandel erlaubt die Berechnung von Ausschnitten der
// Mandelbrot-Menge.
//
// # Ein Titel
//
// Und noch etwas Text.
//
// * Listen?
// * Listen?
//
// Weiterer Text...
package big

import (
    "image"
    "image/color"
    "image/png"
    "log"
    "math"
    "math/big"
    "math/rand"
    "os"
    "strings"

    . "github.com/stefan-muehlebach/mandel"
)

const (
    escRadius  = 256.0
    escRadius2 = escRadius * escRadius
)

// ----------------------------------------------------------------------------
//
// Field --
//
// Enthaelt alle Angaben zu einem darstellbaren Bild der Mandelbrot-Menge.
type Field struct {
    cols, rows int
    maxIter    float64
    pal        *Palette
    f          [][]float64
    iterHist   []float64
}

// Erstellt ein neues Feld mit cols Spalten und rows Zeilen.
func NewField(cols, rows int) *Field {
    var f *Field

    f = new(Field)
    f.cols = cols
    f.rows = rows
    f.f = make([][]float64, rows)
    for i := 0; i < rows; i++ {
        f.f[i] = make([]float64, cols)
    }
    return f
}

// Berechnet die Mandelbrot-Menge ueber dem Feld f mit der Ansicht v.
func (f *Field) CalcMandelbrot(v *View) {
    var dx, dy, xmin, ymax, cx, cy, zx, zy, zx2, zy2, rad, zn *big.Float
    var sx, sy *big.Float
    var escRad, two, half *big.Float
    var iter, znf, nu float64
    var row, col, it int
    var iterate func(x, y *big.Float) float64
    var superSampled func(x, y *big.Float, samp int) float64

    iterate = func(cx, cy *big.Float) (iter float64) {
        zx.SetFloat64(0.0)
        zy.SetFloat64(0.0)
        zx2.SetFloat64(0.0)
        zy2.SetFloat64(0.0)
        for it = 0; (it < v.it) && (escRad.Cmp(rad.Add(zx2, zy2)) > 0); it++ {
            zy.Add(zy.Mul(zy.Mul(zx, zy), two), cy)
            zx.Add(zx.Sub(zx2, zy2), cx)
            zx2.Mul(zx, zx)
            zy2.Mul(zy, zy)
        }
        iter = float64(it)
        if it < v.it {
            zn.Sqrt(rad)
            znf, _ = zn.Float64()
            nu = math.Log(math.Log(znf)*math.Log2E) * math.Log2E
            iter += 1.0 - nu
        }
        return iter
    }

    superSampled = func(cx, cy *big.Float, samp int) (iter float64) {
        iter = iterate(cx, cy)
        for s := 0; s < samp; s++ {
            sx.Add(cx, sx.Mul(dx, big.NewFloat(rand.Float64()-0.5)))
            sy.Add(cy, sy.Mul(dy, big.NewFloat(rand.Float64()-0.5)))
            iter += iterate(sx, sy)
        }
        iter /= float64(samp + 1)
        return iter
    }

    dx = big.NewFloat(0.0).SetPrec(100)
    dy = big.NewFloat(0.0).SetPrec(100)
    xmin = big.NewFloat(0.0).SetPrec(100)
    ymax = big.NewFloat(0.0).SetPrec(100)
    cx = big.NewFloat(0.0).SetPrec(100)
    cy = big.NewFloat(0.0).SetPrec(100)
    sx = big.NewFloat(0.0).SetPrec(100)
    sy = big.NewFloat(0.0).SetPrec(100)
    zx = big.NewFloat(0.0).SetPrec(100)
    zy = big.NewFloat(0.0).SetPrec(100)
    zx2 = big.NewFloat(0.0).SetPrec(100)
    zy2 = big.NewFloat(0.0).SetPrec(100)
    zn = big.NewFloat(0.0).SetPrec(100)

    rad = big.NewFloat(0.0).SetPrec(100)
    escRad = big.NewFloat(escRadius2).SetPrec(100)
    two = big.NewFloat(2.0).SetPrec(100)
    half = big.NewFloat(0.5).SetPrec(100)

    f.maxIter = float64(v.it)

    dx.Quo(v.w, big.NewFloat(float64(f.cols)))
    dy.Set(dx)

    xmin.Sub(v.x, xmin.Mul(v.w, half))
    ymax.Add(v.y, ymax.Mul(dy, big.NewFloat(float64(f.rows)/2.0)))

    cy.Set(ymax)
    for row = 0; row < f.rows; row++ {
        cx.Set(xmin)
        for col = 0; col < f.cols; col++ {
            if v.samp > 0 {
                iter = superSampled(cx, cy, v.samp)
            } else {
                iter = iterate(cx, cy)
            }
            if iter == f.maxIter {
                f.f[row][col] = -1.0
            } else {
                f.f[row][col] = iter
            }
            cx.Add(cx, dx)
        }
        cy.Sub(cy, dy)
    }
}

// Fuegt dem Feld eine Palette hinzu. Eine bereits hinterlegte Palette wird
// damit ueberschrieben.
func (f *Field) AddPalette(p *Palette) {
    f.pal = p
}

// Passt die Groesse der Palette der maximalen Anzahl von Iterationen an.
// Diese Methode MUSS vor der Ausgabe des Bildes als PNG aufgerufen werden,
// ansonsten wird das Bild etwas 'bi color'... ;-)
func (f *Field) AdjPalette() {
    if f.pal.IsLenMaxIter() {
        f.pal.SetLength(int(f.maxIter))
    }
}

// Methoden des image.Image Interfaces
func (f *Field) ColorModel() color.Model {
    return color.RGBAModel
}

func (f *Field) Bounds() image.Rectangle {
    return image.Rect(0, 0, f.cols, f.rows)
}

func (f *Field) At(x, y int) color.Color {
    return f.pal.GetColor(f.f[y][x])
}

// ----------------------------------------------------------------------------
//
// Diese Methode(n) (Draw und WritePNG) werden eigentlich nicht mehr
// benoetigt, da 'Field' nun 'image.Image' implementiert.
func (f *Field) Draw(img *image.RGBA) {
    if f.pal.IsLenMaxIter() {
        f.pal.SetLength(int(f.maxIter))
    }
    for row := 0; row < f.rows; row++ {
        for col := 0; col < f.cols; col++ {
            it := f.f[row][col]
            img.Set(col, row, f.pal.GetColor(it))
        }
    }
}

// Speichert die berechnete Menge als PNG in einem bestimmten Verzeichnis
// und unter einem bestimmten Filenamen ab.
func (f *Field) WritePNG(dirName, fileName string) {
    var img *image.RGBA
    img = image.NewRGBA(image.Rect(0, 0, f.cols, f.rows))
    f.Draw(img)
    s := []string{dirName, fileName}
    fh, err := os.Create(strings.Join(s, "/"))
    if err != nil {
        log.Fatal(err)
    }
    png.Encode(fh, img)
    fh.Close()
}
