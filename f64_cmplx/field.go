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
package f64_cmplx

import (
    "image"
    "image/color"
    "image/png"
    "log"
    "math"
    "math/cmplx"
    "os"
    "strings"
    . "github.com/stefan-muehlebach/mandel"
)

const (
    escRad  = 2.0
    escRad2 = escRad * escRad
)

//----------------------------------------------------------------------------
//
// Field --
//
// Enthaelt alle Angaben zu einem darstellbaren Bild der Mandelbrot-Menge.
//
type Field struct {
    cols, rows int
    maxIter float64
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
    var dRe, dIm, reDiff, imDiff, iter, zn, nu float64
    var row, col, it int
	var z0, z, c complex128

    f.maxIter = float64(v.it)
	
	dRe = v.w / float64(f.cols)
	dIm = dRe
	
	z0 = v.z - complex(float64(f.cols/2)*dRe, -float64(f.rows/2)*dIm)
	
	imDiff = 0.0
    for row = 0; row < f.rows; row++ {
        reDiff = 0.0
		for col = 0; col < f.cols; col++ {
			c  = z0 + complex(reDiff, imDiff)
            z  = 0.0 + 0.0i
            it = 0
            for (real(z)*real(z) + imag(z)*imag(z) <= escRad2) && (it < v.it) {
				z = z*z + c
                it++
            }
            if it < v.it {
                zn = math.Log(cmplx.Abs(z))
                nu = math.Log(zn * math.Log2E) * math.Log2E
                iter = float64(it) + 1.0 - nu;
            } else {
                iter = -1.0
            }
            f.f[row][col] = iter
			reDiff += dRe
        }
		imDiff -= dIm
    }
}

// Fuegt dem Feld eine Palette hinzu. Eine bereits hinterlegte Palette wird
// damit ueberschrieben.
//
func (f *Field) AddPalette(p *Palette) {
    f.pal = p
}

// Passt die Groesse der Palette der maximalen Anzahl von Iterationen an.
// Diese Methode MUSS vor der Ausgabe des Bildes als PNG aufgerufen werden,
// ansonsten wird das Bild etwas 'bi color'... ;-)
//
func (f *Field) AdjPalette() {
    if f.pal.IsLenMaxIter() {
        f.pal.SetLength(int(f.maxIter))
    }
}

// Methoden des image.Image Interfaces
//
func (f *Field) ColorModel() (color.Model) {
	return color.RGBAModel
}

func (f *Field) Bounds() (image.Rectangle) {
	return image.Rect(0, 0, f.cols, f.rows)
}

func (f *Field) At(x, y int) (color.Color) {
	return f.pal.GetColor(f.f[y][x])
}

//----------------------------------------------------------------------------
//
// Diese Methode(n) (Draw und WritePNG) werden eigentlich nicht mehr
// benoetigt, da 'Field' nun 'image.Image' implementiert.
//
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
//
func (f *Field) WritePNG(dirName, fileName string) {
    var img *image.RGBA
    img = image.NewRGBA(image.Rect(0, 0, f.cols, f.rows))
    f.Draw(img)
    s := []string{dirName, fileName}
    fh, err := os.Create(strings.Join (s, "/"))
    if err != nil {
        log.Fatal(err)
    }
    png.Encode(fh, img)
    fh.Close()
}

