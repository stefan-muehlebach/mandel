package f64

// View ist eine Ansicht der komplexen Zahlenebene mit einem bestimmten
// Anzeigebereich, einem Mittelpunkt und einer bestimmten Anzahl Iterationen
// sowie Anzahl Samples.
type f64View struct {
    x, y, w float64
    it      int
}

// Erstellt eine neue, leere View ohne definierten Bereich oder Mittelpunkt.
func NewView() *f64View {
    return &f64View{}
}

// Definiert die Parameter der View.
func (v *f64View) SetValues(x, y, w float64, it int) {
    v.x, v.y, v.w, v.it = x, y, w, it
}

func (v *f64View) Values() (x, y, w float64, it int) {
    return v.x, v.y, v.w, v.it
}
