// Dieses Package habe ich erstellt, um mit Go Bildserien der Mandelbrotmenge
// zu erstellen.
package mandel

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"os"
	"path/filepath"
	"regexp"
)

var (
	// confDir ist der Name des Verzeichnisses, in welchem alle Konfigurations
	// Dateien erwartet werden. Dies ist aktuell fest eingestellt und zeigt
	// auf '~/.config/mandel'
	confDir = filepath.Join(".config", "mandel")
)

// OpenConfFile ist eine interne Funktion, mit welcher Konfigurationsdateien
// im Verzeichnis [confDir] standardisiert angesprochen werden koennen.
func OpenConfFile(fileName string) (*os.File, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	absPath := filepath.Join(homeDir, confDir, fileName)
	fh, err := os.Open(absPath)
	return fh, err
}

type View interface {
	Values() (x, y, w float64, maxIt int)
	SetValues(x, y, w float64, maxIt int)
}

type Path interface {
	Read(pathName string) error
	AddView(x, y, w float64, maxIt int)
	NumViews() int
	GetView(t float64) View
}

type Field interface {
	CalcMandelbrot(v View)
	AddPalette(pal Palette)
	AdjPalette()
	Write(fileName string) error
	Read(fileName string) error
	ColorModel() color.Model
	Bounds() image.Rectangle
	At(x, y int) color.Color
}

type Palette interface {
	SetLength(len int)
	Length() int
	LenIsMaxIter()
	LenIsNotMaxIter()
	IsLenMaxIter() bool
	SetOffset(offset float64)
	Offset() float64
	GetColor(f float64) color.RGBA

	GetRegexp() *regexp.Regexp
	ProcessLine(line string) error
	Update()
	Ready() bool
}

// Der Typ SampleMode gibt an, von wie vielen Punkten der Mittelwert fuer ein
// Pixel in der Anzeige berechnet werden soll. Es stehen aktuell die Groessen
// 1x2, 2x2, 4x4 und 8x8 zur Verfuegung.
type SampleMode int

const (
	Samp1x1 SampleMode = 1
	Samp2x2 SampleMode = 2
	Samp4x4 SampleMode = 4
	Samp8x8 SampleMode = 8
)

// Mit String wird eine brauchbare textuelle Darstellung der Sample-Groesse
// erstellt. Wird u.A. bei der Ausgabe der Flags verwendet, mit denen das
// Programm aufgerufen wird (ist auch Teil des [flag.Value] Interfaces).
func (sm SampleMode) String() string {
	switch sm {
	case Samp1x1:
		return fmt.Sprintf("1x1 sample size")
	case Samp2x2:
		return fmt.Sprintf("2x2 sample size")
	case Samp4x4:
		return fmt.Sprintf("4x4 sample size")
	case Samp8x8:
		return fmt.Sprintf("8x8 sample size")
	default:
		return "Unknown sample size"
	}
}

func (sm *SampleMode) Set(s string) error {
	switch s {
	case "1x1":
		*sm = Samp1x1
	case "2x2":
		*sm = Samp2x2
	case "4x4":
		*sm = Samp4x4
	case "8x8":
		*sm = Samp8x8
	default:
		return errors.New("Unknow sample mode: " + s)
	}
	return nil
}
