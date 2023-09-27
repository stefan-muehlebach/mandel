package f64

import (
	"log"
    "testing"
)

const (
    dt = 0.003
)

var (
    path *Path
    view *View
    view2 View
    err error
)

func init() {
    path, err = ReadPath("Default")
    if err != nil {
        log.Fatalf("could not read path: %v", err)
    }
    view = NewView()
}

func BenchmarkGetView(b *testing.B) {
    for i:=0; i<b.N; i++ {
        for t:=0.0; t<=1.0; t+=dt {
            path.GetView(t, view)
        }
    }
}

func BenchmarkGetView2(b *testing.B) {
    for i:=0; i<b.N; i++ {
        for t:=0.0; t<=1.0; t+=dt {
            view2 = path.GetView2(t)
        }
    }
}

