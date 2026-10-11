// Command figures draws every figure the documentation shows, with the
// go-numerics libraries at the versions go.mod names computing the data and
// go-plot drawing it, into ../static/images/figures/<name>.svg.
// Each figure is one file here, and the pages include that file verbatim, so
// the code shown is the code that drew the picture beside it.
//
//	cd figures && go run .
//
// The CI runs it and fails if any image differs from the committed one.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/go-plot/plot"
)

var figures = map[string]func() *plot.Figure{}

func register(name string, f func() *plot.Figure) { figures[name] = f }

func main() {
	out := filepath.Join("..", "static", "images", "figures")
	if len(os.Args) > 1 {
		out = os.Args[1]
	}
	if err := os.MkdirAll(out, 0o755); err != nil {
		panic(err)
	}
	names := make([]string, 0, len(figures))
	for n := range figures {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		f, err := os.Create(filepath.Join(out, n+".svg"))
		if err != nil {
			panic(err)
		}
		if err := figures[n]().SaveSVG(f); err != nil {
			panic(err)
		}
		if err := f.Close(); err != nil {
			panic(err)
		}
		fmt.Println(n + ".svg")
	}
}

// linspace returns n evenly spaced values from lo to hi, as numpy.linspace.
func linspace(lo, hi float64, n int) []float64 {
	v := make([]float64, n)
	for i := range v {
		v[i] = lo + (hi-lo)*float64(i)/float64(n-1)
	}
	return v
}
