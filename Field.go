//
// Copyright 2018 by Stefan Muehlebach
//

// Das Package mandel erlaubt die Berechnung von Ausschnitten der
// Mandelbrot-Menge.
//
// Ein Titel
//
// Und noch etwas Text.
//
// * Listen?
// * Listen?
//
// Weiterer Text...
//
package mandel

import (
    "math"
    "image"
    "image/color"
    "image/png"
    "log"
    "os"
    "strings"
)

//----------------------------------------------------------------------------
//
// Field --
//
// Enthaelt alle Angaben zu einem darstellbaren Bild der Mandelbrot-Menge.
//
type Field struct {
    cols, rows int
    minIter, maxIter float64
    pal *Palette
    f [][]float64
}

// Erstellt ein neues Feld mit cols Spalten und rows Zeilen.
//
func NewField(cols, rows int) (*Field) {
    var f *Field

    f = new(Field)
    f.cols = cols
    f.rows = rows
    f.f = make([][]float64, rows)
    for i:=0; i<rows; i++ {
        f.f[i] = make([]float64, cols)
    }

    return f
}

// Berechnet die Mandelbrot-Menge ueber dem Feld f mit der Ansicht v.
//
func (f *Field) CalcMandelbrot(v *View) {
    var dx, dy, xmin, ymax, cx, cy, zx, zy, zxtemp, iter, zn, nu float64
	var minIter, maxIter float64
    var row, col, it int

    minIter = float64(v.it)
    maxIter = float64(v.it)

    dx = v.w / float64(f.cols)
    dy = dx
    xmin = v.x - v.w / 2.0
    ymax = v.y + dy * float64(f.rows) / 2.0

    cy = ymax
    for row = 0; row < f.rows; row++ {
        cx = xmin
        for col = 0; col < f.cols; col++ {
            zx = 0.0
            zy = 0.0
            it = 0
            for (zx*zx + zy*zy < (1 << 16)) && (it < v.it) {
                zxtemp = zx*zx - zy*zy + cx
                zy = 2.0 * zx * zy + cy
                zx = zxtemp
                it++
            }
            if it < v.it {
                zn = math.Sqrt(zx*zx + zy*zy)
                nu = math.Log(math.Log(zn)) / math.Log(2.0)
                iter = float64(it) + 1.0 - nu;
            } else {
                iter = float64(it)
            }
            if iter < minIter {
                minIter = iter
            }
            // if iter > maxIter {
                // maxIter = iter
            // }
            f.f[row][col] = iter
            cx += dx
        }
        cy -= dy
    }
    f.minIter = minIter
    f.maxIter = maxIter
}

// Fuegt dem Feld eine Palette hinzu. Eine bereits hinterlegte Palette wird
// damit ueberschrieben.
//
func (f *Field) AddPalette(p *Palette) {
    f.pal = p
}

// Ermittelt den Farbwert eines bestimmten Punktes in der Menge.
//
func (f *Field) GetColor(col, row int) (color.RGBA) {
    v := f.f[row][col]/f.maxIter
    if v == 1.0 {
        return color.RGBA{0, 0, 0, 255}
    } else {
        return f.pal.GetColor(v)
    }
}

// Speichert die berechnete Menge als PNG in einem bestimmten Verzeichnis
// und unter einem bestimmten Filenamen ab.
//
func (f *Field) WritePNG(dirName, fileName string) {
    img := image.NewRGBA(image.Rect(0, 0, f.cols, f.rows))
    for row := 0; row < f.rows; row++ {
        for col := 0; col < f.cols; col++ {
			img.Set(col, row, f.pal.GetColor(f.f[row][col]/f.maxIter))
        }
    }

    s := []string{dirName, fileName}
    fh, err := os.Create(strings.Join (s, "/"))
    if err != nil {
        log.Fatal(err)
    }
    png.Encode(fh, img)
    fh.Close()
}
