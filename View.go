package mandel

//----------------------------------------------------------------------------
//
// View --
//
// View ist eine Ansicht der komplexen Zahlenebene mit einem bestimmten
// Anzeigebereich, einem Mittelpunkt und einer bestimmten Anzahl Iterationen.
//
type View struct {
    x, y, w float64
    it int
}

// Erstellt eine neue View ohne definierten Bereich oder Mittelpunkt.
//
func NewView() (*View) {
    return &View{}
}

// Definiert die Parameter der View.
//
func (v *View) SetValues(x, y, w float64, it int) {
    v.x, v.y, v.w, v.it = x, y, w, it
}

