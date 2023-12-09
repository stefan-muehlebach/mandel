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
package f64

import (
    "encoding/gob"
    "image"
    "image/color"
    _ "image/png"
    _ "log"
    "math"
    _ "math/rand"
    "os"
    _ "strings"

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
type f64Field struct {
    Cols, Rows int
    MaxIter    float64
    pal        Palette
    F          [][]float64
    sm         SampleMode
    // iterHist []float64
}

// Erstellt ein neues Feld mit cols Spalten und rows Zeilen.
func NewField(cols, rows int, sm SampleMode) *f64Field {
    f := &f64Field{}
    f.Cols = cols
    f.Rows = rows
    f.sm = sm
    f.F = make([][]float64, f.Rows)
    for i := 0; i < f.Rows; i++ {
        f.F[i] = make([]float64, f.Cols)
    }
    return f
}

// Berechnet die Mandelbrot-Menge ueber dem Feld f mit der Ansicht v.
func (f *f64Field) CalcMandelbrot(v View) {
    var dx, dy, xmin, ymax, h, cx, cy, iter float64
    // var total float64
    var row, col int
    
    x, y, w, it := v.Values()

    f.MaxIter = float64(it)
    // if len(f.iterHist) != v.it {
    //     f.iterHist = make([]float64, v.it)
    // } else {
    //     for i, _ := range f.iterHist {
    //         f.iterHist[i] = 0.0
    //     }
    // }

    dx = w / float64(f.Cols)
    dy = dx
    h = dy * float64(f.Rows)
    xmin = x - w/2.0
    ymax = y + h/2.0

    cy = ymax
    for row = 0; row < f.Rows; row++ {
        cx = xmin
        for col = 0; col < f.Cols; col++ {
            if f.skipCalculation(cx, cy) {
                iter = f.MaxIter
            } else {
                iter = f.calcCell(cx, cy, dx, dy, it)
            }
            if iter == f.MaxIter {
                f.F[row][col] = -1.0
            } else {
                f.F[row][col] = iter
            }
            cx += dx
        }
        cy -= dy
    }

    // total = 0.0
    // for _, v := range f.iterHist {
    //     total += v
    // }
    // for i, v := range f.iterHist {
    //     f.iterHist[i] = v / total
    // }
}

func (f *f64Field) calcCell(cx, cy, dx, dy float64, maxIter int) (iter float64) {
    var rx, ry float64
    var cellRow, cellCol int

    if f.sm == Samp1x1 {
        return f.calcPixel(cx, cy, maxIter)
    }

    iter = 0.0
    dx /= float64(f.sm)
    dy /= float64(f.sm)
    ry = cy
    for cellRow = 0; cellRow < int(f.sm); cellRow++ {
        rx = cx
        for cellCol = 0; cellCol < int(f.sm); cellCol++ {
            iter += f.calcPixel(rx, ry, maxIter)
            rx += dx
        }
        ry -= dy
    }
    return iter / (float64(f.sm) * float64(f.sm))
}

func (f *f64Field) calcPixel(cx, cy float64, maxIter int) (iter float64) {
    var zx, zy, zx2, zy2 float64
    var zn, nu float64
    var it int

    zx, zy = 0.0, 0.0
    zx2, zy2 = 0.0, 0.0
    for it = 0; (it < maxIter) && (zx2+zy2 <= escRadius2); it++ {
        zy = 2.0*zx*zy + cy
        zx = zx2 - zy2 + cx
        zx2 = zx * zx
        zy2 = zy * zy
    }
    iter = float64(it)
    if it < maxIter {
        zn = math.Log(zx2+zy2) / 2.0
        nu = math.Log(zn*math.Log2E) * math.Log2E
        iter += 1.0 - nu
    }
    return
}

// Bestimmte Bereiche der komplexen Ebene gehoeren garantiert zur
// Mandelbrotmenge und brauchen nicht noch berechnet zu werden. Diese
// Funktion prueft, ob die komplexe Zahl cx + i*cy bestimmt in der Menge liegt
// und liefert true, falls dies der Fall ist und false, falls nicht.
func (f *f64Field) skipCalculation(cx, cy float64) bool {
    if (cx > -1.2 && cx <= -1.1 && cy > -0.1 && cy < 0.1) ||
            (cx > -1.1 && cx <= -0.9 && cy > -0.2 && cy < 0.2) ||
            (cx > -0.9 && cx <= -0.8 && cy > -0.1 && cy < 0.1) ||
            (cx > -0.69 && cx <= -0.61 && cy > -0.2 && cy < 0.2) ||
            (cx > -0.61 && cx <= -0.5 && cy > -0.37 && cy < 0.37) ||
            (cx > -0.5 && cx <= -0.39 && cy > -0.48 && cy < 0.48) ||
            (cx > -0.39 && cx <= 0.14 && cy > -0.55 && cy < 0.55) ||
            (cx > 0.14 && cx < 0.29 && cy > -0.42 && cy < -0.07) ||
            (cx > 0.14 && cx < 0.29 && cy > 0.07 && cy < 0.42) {
        return true
    } else {
        return false
    }
}

// Fuegt dem Feld eine Palette hinzu. Eine bereits hinterlegte Palette wird
// durch die neue Palette ersetzt.
func (f *f64Field) AddPalette(p Palette) {
    f.pal = p
}

// Passt die Groesse der Palette der maximalen Anzahl von Iterationen an.
// Diese Methode MUSS vor der Ausgabe des Bildes als PNG aufgerufen werden,
// ansonsten wird das Bild etwas 'bi color'... ;-)
func (f *f64Field) AdjPalette() {
    if f.pal.IsLenMaxIter() {
        f.pal.SetLength(int(f.MaxIter))
    }
}

// Methoden fuer das Speichern, resp. Lesen der binaeren Daten eines Feldes.
// TO DO: ev gibt es dafuer auch eine Go-konvformere Loesung, analog der
// unten gezeigten Loesung fuer das Speichern der Daten als PNG- oder JPG-File.
func (f *f64Field) Write(fileName string) (error) {
    fh, err := os.Create(fileName)
    if err != nil {
        return err
    }
    defer fh.Close()
    enc := gob.NewEncoder(fh)
    err = enc.Encode(f)
    return err
}

func (f *f64Field) Read(fileName string) (error) {
    fh, err := os.Open(fileName)
    if err != nil {
        return err
    }
    defer fh.Close()
    dec := gob.NewDecoder(fh)
    err = dec.Decode(f)
    return err
}

// Methoden des image.Image Interfaces. Auf diese Weise wird die Speicherung
// der Felddaten als PNG oder JPG realisiert.
func (f *f64Field) ColorModel() color.Model {
    return color.RGBAModel
}

func (f *f64Field) Bounds() image.Rectangle {
    return image.Rect(0, 0, f.Cols, f.Rows)
}

func (f *f64Field) At(x, y int) color.Color {
    return f.pal.GetColor(f.F[y][x])
}

