package big

import (
    "math/big"
)

// View ist eine Ansicht der komplexen Zahlenebene mit einem bestimmten
// Anzeigebereich, einem Mittelpunkt und einer bestimmten Anzahl Iterationen.
type View struct {
    x, y, w  *big.Float
    it, samp int
}

// Erstellt eine neue, leere View ohne definierten Bereich oder Mittelpunkt.
func NewView() *View {
    v := &View{}
    v.x = big.NewFloat(0.0).SetPrec(100)
    v.y = big.NewFloat(0.0).SetPrec(100)
    v.w = big.NewFloat(0.0).SetPrec(100)
    return v
}

// Definiert die Parameter der View.
func (v *View) SetValues(x, y, w *big.Float, it, samp int) {
    v.x.Set(x)
    v.y.Set(y)
    v.w.Set(w)
    v.it = it
    v.samp = samp
}

func (v *View) GetValues() (x, y, w *big.Float, it, samp int) {
    return v.x, v.y, v.w, v.it, v.samp
}
