package f64_cmplx

// View ist eine Ansicht der komplexen Zahlenebene mit einem bestimmten
// Anzeigebereich, einem Mittelpunkt und einer bestimmten Anzahl Iterationen.
//
type View struct {
	z complex128
    w float64
    it int
}

// Erstellt eine neue, leere View ohne definierten Bereich oder Mittelpunkt.
//
func NewView() (*View) {
    return &View{}
}

// Definiert die Parameter der View.
//
func (v *View) SetValues(z complex128, w float64, it int) {
    v.z, v.w, v.it = z, w, it
}

func (v *View) GetValues() (z complex128, w float64, it int) {
    return v.z, v.w, v.it
}
