package mandel_test

import (
    "mandel"
)

func ExampleNewView () {
    v := mandel.NewView ()
}

func ExampleField_CalcMandelbrot () {
    f := mandel.NewField (640, 480)
    v := mandel.NewView ()
    v.SetValues (3.5, 2.625, -1.0, 0.0, 1024)
    f.CalcMandelbrot (v)
    f.WritePNG ("images", "mandel.png")
}

