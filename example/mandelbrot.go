package main

import (
    "example/mandel"
    "fmt"
    "time"
    "os"
    "flag"
	"log"
    _ "runtime/trace"
	_ "math/big"
)

func Worker(id int, ch chan int, done chan bool, path *mandel.Path,
        pal *mandel.Palette) {
    var t1, t2 time.Time
    var field *mandel.Field
    var view *mandel.View
    var fileName string

    field = mandel.NewField(cols, rows)
    field.AddPalette(pal)

    view = mandel.NewView()

    for {
        i := <- ch
        if i < 0 {
            break
        }
        t := float64(i) / float64(nImages-1)
        path.GetView(t, view)
		//fmt.Printf ("view width: %.50f\n", view.GetWidth ())
		//fmt.Printf ("view x    : %.50f\n", view.GetX ())
		//fmt.Printf ("view y    : %.50f\n", view.GetY ())
        t1 = time.Now()
        field.CalcMandelbrot(view)
        t2 = time.Now()
        fmt.Printf("goroutine[%d], image: %d, dt: %v\n", id, i, t2.Sub(t1))
        fileName = fmt.Sprintf("mandel%03d.png", i)
        field.WritePNG(imageDir, fileName)
    }
    done <- true
}

var nImages int
var cols, rows int
var imageDir string

func main() {
    var paletteFile, paletteName, pathFile, pathName string
    var nWorkers int

    var palette *mandel.Palette
    var path *mandel.Path
    var ch chan int
    var done chan bool

    //var fh *os.File

    //fh, _ = os.Create ("mandel_main.trace")
    //trace.Start (fh)

    flag.IntVar(&cols, "cols", 640, "number of columns")
    flag.IntVar(&rows, "rows", 480, "number of rows")
    flag.StringVar(&paletteFile, "palFile", "palette.ini", "palette file")
	flag.StringVar(&paletteName, "palName", "DEFAULT", "palette name")
    flag.StringVar(&pathFile, "pathFile", "path.ini", "path file")
    flag.StringVar(&pathName, "pathName", "DEFAULT", "path name")
    flag.StringVar(&imageDir, "outdir", "images", "image directory")
    flag.IntVar(&nWorkers, "workers", 4, "number of go routines")
    flag.IntVar(&nImages, "images", 128, "number of images to create")
    flag.Parse()
	
	fmt.Printf("image size  : %d x %d\n", cols, rows)
	fmt.Printf("palette file: %s\n", paletteFile)
	fmt.Printf("palette name: %s\n", paletteName)
	fmt.Printf("path file   : %s\n", pathFile)
	fmt.Printf("path name   : %s\n", pathName)
	fmt.Printf("image dir   : %s\n", imageDir)
	fmt.Printf("# workers   : %d\n", nWorkers)
	fmt.Printf("# images    : %d\n", nImages)

    os.Mkdir(imageDir, 0755)

    palette = mandel.NewPalette()
	err := palette.ReadFile(paletteFile, paletteName)
	if err != nil {
		log.Fatal(err)
	}
	
    path = mandel.NewPath()
    err = path.ReadFile(pathFile, pathName)
	if err != nil {
		log.Fatal(err)
	}

    ch = make(chan int)
    done = make(chan bool)

    for i:=0; i<nWorkers; i++ {
        go Worker(i, ch, done, path, palette)
    }
    for i:=0; i<nImages; i++ {
        ch <- i
    }
    for i:=0; i<nWorkers; i++ {
        ch <- -1
    }
    for i:=0; i<nWorkers; i++ {
        <- done
    }

    //trace.Stop ()
}

