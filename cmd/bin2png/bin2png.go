package main

import (
	"flag"
	"fmt"
	_ "image"
	_ "image/color"
	"image/png"
	"io/fs"
	"log"
	_ "math/big"
	"os"
	"path"
	"runtime"
	_ "runtime/trace"
	_ "time"

	"github.com/stefan-muehlebach/mandel"
	"github.com/stefan-muehlebach/mandel/f64"
)

const (
	defPalName     = "Default"
	defPalLength   = -1
	defPalOffset   = 0.0
    defBinDir      = "data"
    defImgDir      = "images"
    binFilePattern = "*.bin"
	imgFilePattern = "img%05d.png"
)

var (
	palName        string
	palLength      int
	palOffset      float64
	binDir, imgDir string
)

func check(err error) {
    if err != nil {
        log.Fatalf(err.Error())
    }
}

func main() {
	var nWorkers int
	var palette mandel.Palette
	var field mandel.Field
	var outFile string
	var fh *os.File
	var err error
	var i int

	nWorkers = runtime.NumCPU()

	flag.StringVar(&palName, "palette", defPalName, "palette name")
	flag.IntVar(&palLength, "palLength", defPalLength,
		"palette length (default: equal to max iterations)")
	flag.Float64Var(&palOffset, "palOffset", defPalOffset,
		"offset (in %) of the first color of the palette")
	flag.StringVar(&binDir, "bindir", defBinDir, "input directory with binary files")
	flag.StringVar(&imgDir, "imgdir", defImgDir, "output directory for images")
	flag.IntVar(&nWorkers, "workers", 4, "number of go routines")
	flag.Parse()

	fmt.Printf("palette name    : %s\n", palName)
	fmt.Printf("palette length  : %d\n", palLength)
	fmt.Printf("palette offset  : %.1f%%\n", palOffset)
	fmt.Printf("input dir       : %s\n", binDir)
	fmt.Printf("output dir      : %s\n", imgDir)
	fmt.Printf("#workers        : %d\n", nWorkers)

	os.Mkdir(imgDir, 0755)
	fileSystem := os.DirFS(".")
	field = f64.NewField(1, 1, mandel.Samp1x1)
	palette, err = mandel.NewPalette(palName)
	check(err)
	if palLength < 0 {
		palette.LenIsMaxIter()
	} else {
		palette.LenIsNotMaxIter()
		palette.SetLength(palLength)
	}
	palette.SetOffset(palOffset / 100.0)
	field.AddPalette(palette)
	i = 0
	fs.WalkDir(fileSystem, binDir,
        func(inFile string, d fs.DirEntry, err error) error {
        		check(err)
        		if d.IsDir() {
        			return nil
        		}
        		fmt.Printf("processing '%s'\n", inFile)
        		err = field.Read(inFile)
        		check(err)
        		field.AdjPalette()
        		outFile = fmt.Sprintf(imgFilePattern, i)
        		fh, err = os.Create(path.Join(imgDir, outFile))
        		check(err)
        		png.Encode(fh, field)
        		fh.Close()
        		i++
        		return nil
	    })
}
