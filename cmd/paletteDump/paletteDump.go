package main

import (
	_ "flag"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"

	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/color"
	"github.com/stefan-muehlebach/gg/fonts"
	"github.com/stefan-muehlebach/mandel"
)

const (
	NumColumns     = 2
	ColorBarWidth  = 1024
	ColorBarHeight = 100
	TextHeight     = 30
	Padding        = 10
	StripeWidth    = ColorBarWidth + 2*Padding
	StripeHeight   = ColorBarHeight + TextHeight + 2*Padding
)

func main() {
	var palType string

	palNameList, err := mandel.PaletteNames()
	if err != nil {
		log.Fatalf("couldn't read palette names: %v", err)
	}
	numPals := len(palNameList)
	numRows := numPals / NumColumns
	if numPals%NumColumns != 0 {
		numRows += 1
	}
	imgWidth := NumColumns * StripeWidth
	imgHeight := numRows * StripeHeight
	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))
	gc := gg.NewContextForRGBA(img)

	gc.SetFillColor(color.WhiteSmoke)
	gc.Clear()
	gc.SetFontFace(fonts.NewFace(fonts.GoRegular, 18.0))
	gc.SetStrokeColor(color.Black)
	for i, palName := range palNameList {
		col := i / numRows
		row := i % numRows
		x0 := float64(col * StripeWidth)
		y0 := float64(row * StripeHeight)

		gc.SetStrokeWidth(1.0)
		gc.DrawRectangle(x0, y0, StripeWidth, StripeHeight)
		gc.FillStroke()

		fmt.Printf("  [%2d]: %s\n", i, palName)
		pal, err := mandel.NewPalette(palName)
		if err != nil {
			log.Fatalf("couldn't create palette: %v", err)
		}
		switch pal.(type) {
		case *mandel.GradientPalette:
			palType = "Gradient Palette"
		case *mandel.ProcPalette:
			palType = "Procedure Palette"
		}
		pal.SetLength(ColorBarWidth)
		pal.LenIsMaxIter()
		pal.SetOffset(0.0)

		for x := 0; x < ColorBarWidth; x++ {
			color := pal.GetColor(float64(x))
			for y := 0; y < ColorBarHeight; y++ {
				img.Set(int(x0+Padding)+x, int(y0+TextHeight+Padding)+y, color)
			}
		}
		gc.SetStrokeColor(color.DarkSlateGrey)
		gc.SetStrokeWidth(3.0)
		gc.DrawRectangle(x0+Padding, y0+TextHeight+Padding, ColorBarWidth, ColorBarHeight)
		gc.Stroke()

		gc.DrawStringAnchored(palName, x0+Padding, y0+Padding+TextHeight-Padding, 0.0, 0.0)
		gc.DrawStringAnchored(palType, x0+StripeWidth-Padding,
			y0+Padding+TextHeight-Padding, 1.0, 0.0)
	}
	fileName := "paletteDump.png"
	fh, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("couldn't create file: %v", err)
	}
	png.Encode(fh, img)
	fh.Close()
}
