package main

import (
	"log"
    "flag"
    "fmt"
    _ "image"
    _ "image/color"
    "image/png"
    _ "log"
    _ "math/big"
    "os"
    "path/filepath"
    "runtime"
    _ "runtime/trace"
    "time"

    "github.com/stefan-muehlebach/mandel"
    "github.com/stefan-muehlebach/mandel/f64"
)

const (
    defWriteBin   = false
    defNumCols    = 320
    defNumRows    = 240
    defPathName   = "Default"
    defPalName    = "Default"
    defPalLength  = -1
    defPalOffset  = 0.0
    defBinDir     = "data"
    defImgDir     = "images"
    defNumImages  = 128
    defSampleMode = mandel.Samp1x1

    imgFilePattern = "img%05d.png"
    binFilePattern = "img%05d.bin"
)

var (
    writeBin       bool
    cols, rows     int
    palName        string
    palLength      int
    palOffset      float64
    imgDir, binDir string
    outDir         string
    pathName       string
    numImages      int
    totalImages    int
    sampleMode     mandel.SampleMode = defSampleMode
 ) 

func check(err error) {
    if err != nil {
        log.Fatalf(err.Error())
    }
}

// Diese Funktion wird von mehreren Go-Routinen ausgefuert. Auf diesem Level
// findet die Parallelisierung statt. Gesteuert werden die Routinen ueber die
// Channels ch (Input-Channel fuer die Auftraege) und done (Output-Channel
// fuer die Meldung, dass diese Go-Routine korrekt terminiert).
func Worker(id int, ch <-chan int, done chan<- bool, path mandel.Path) {
    var t1, t2, t3 time.Time
    var t float64
    var i int
    var field mandel.Field
    var palette mandel.Palette
    var view mandel.View
    var outFile string
    var fh *os.File
    var err error

    field = f64.NewField(cols, rows, sampleMode)
    palette, err = mandel.NewPalette(palName)
    check(err)
    if palLength < 0 {
        palette.LenIsMaxIter()
    } else {
        palette.LenIsNotMaxIter()
        palette.SetLength(palLength)
    }
    palette.SetOffset(palOffset/100.0)
    field.AddPalette(palette)

    // view = &f64.View{}

    for {
        i = <-ch
        if i < 0 {
            break
        }
        t1 = time.Now()
        t = float64(i) / float64(totalImages)
        view = path.GetView(t)
        field.CalcMandelbrot(view)
        t2 = time.Now()
        if !writeBin {
            outFile = fmt.Sprintf(imgFilePattern, i)
            fh, err = os.Create(filepath.Join(imgDir, outFile))
            check(err)
            field.AdjPalette()
            png.Encode(fh, field)
            fh.Close()
        } else {
            outFile = fmt.Sprintf(binFilePattern, i)
            err = field.Write(filepath.Join(binDir, outFile))
            check(err)
        }
        t3 = time.Now()
        fmt.Printf("GO[%d]: image: %05d; calc: %v, file creation: %v\n",
                id, i, t2.Sub(t1), t3.Sub(t2))
    }
    done <- true
}

func main() {

    var nWorkers int
    var path mandel.Path
    var ch chan int
    var done chan bool
    var err error

    nWorkers = runtime.NumCPU()

    flag.BoolVar(&writeBin, "bin", defWriteBin, "write data as binary files instead of images")
    flag.IntVar(&cols, "cols", defNumCols, "number of columns")
    flag.IntVar(&rows, "rows", defNumRows, "number of rows")
    flag.StringVar(&palName, "palette", defPalName, "palette name")
    flag.IntVar(&palLength, "palLength", defPalLength, "palette length (default: equal to max iterations)")
    flag.Float64Var(&palOffset, "palOffset", 0.0, "offset (in %) of the first color of the palette")
    flag.StringVar(&pathName, "path", defPathName, "path name")
    flag.StringVar(&imgDir, "imgdir", defImgDir, "output directory for images")
    flag.StringVar(&binDir, "bindir", defBinDir, "output directory for binary files")
    flag.Var(&sampleMode, "sampleMode", "mode of subpixel sampling")
    flag.IntVar(&nWorkers, "workers", 4, "number of go routines")
    flag.IntVar(&numImages, "images", defNumImages, "number of images between two views")
    flag.Parse()

    if writeBin {
        outDir = binDir
    } else {
        outDir = imgDir
    }

    fmt.Printf("write bin files : %v\n", writeBin)
    fmt.Printf("image size      : %d x %d\n", cols, rows)
    fmt.Printf("palette name    : %s\n", palName)
    fmt.Printf("palette length  : %d\n", palLength)
    fmt.Printf("palette offset  : %.1f %%\n", palOffset)
    fmt.Printf("path name       : %s\n", pathName)
    fmt.Printf("output dir      : %s\n", outDir)
    fmt.Printf("#workers        : %d\n", nWorkers)
    fmt.Printf("#images/view    : %d\n", numImages)
    fmt.Printf("sample mode     : %v\n", sampleMode)

    os.Mkdir(outDir, 0755)

    path = f64.NewPath()
    err = path.Read(pathName)
    check(err)

    // if pth.NumViews() == 1 {
    //     field := f64.NewField(cols, rows, sampleMode)
    //     palette, err = mandel.ReadPalette(palName)
    //     utils.Check(err)
    //     if palLength < 0 {
    //         palette.LenIsMaxIter()
    //     } else {
    //         palette.LenIsNotMaxIter()
    //         palette.SetLength(palLength)
    //     }
    //     palette.SetOffset(palOffset)
    //     field.AddPalette(palette)
    //     view := f64.NewView()
    //     t1 := time.Now()
    //     pth.GetView(0.0, view)
    //     field.CalcMandelbrot(view)
    //     t2 := time.Now()
    //     imageFile := fmt.Sprintf("mandel%s.png", pathName)
    //     fh, err := os.Create(imageFile)
    //     if err != nil {
    //         log.Fatal(err)
    //     }
    //     field.AdjPalette()
    //     png.Encode(fh, field)
    //     fh.Close()
    //     t3 := time.Now()
    //     fmt.Printf("single thread: image: '%s', calculation: %v, image creation: %v\n",
    //             imageFile, t2.Sub(t1), t3.Sub(t2))
    //     return
    // }

    ch   = make(chan int)
    done = make(chan bool)
    totalImages = numImages*(path.NumViews()-1) + 1

    t1 := time.Now()
    for i := 0; i < nWorkers; i++ {
        go Worker(i, ch, done, path)
    }
    for i := 0; i < totalImages; i++ {
        ch <- i
    }
    for i := 0; i < nWorkers; i++ {
        ch <- -1
    }
    for i := 0; i < nWorkers; i++ {
        <-done
    }
    d := time.Since(t1)
    fmt.Printf("Calculation took %v\n", d)
}

