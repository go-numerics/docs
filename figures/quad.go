package main

import (
	"fmt"
	"math"

	"github.com/go-numerics/integrate"
	"github.com/go-plot/plot"
)

func init() { register("quad", quad) }

// The README's Quad example: an integrand with a singularity at 0, which
// QUADPACK's extrapolation handles. The title prints what Quad returned.
func quad() *plot.Figure {
	f := func(x float64) float64 { return math.Log(x) / math.Sqrt(x) }
	r, err := integrate.Quad(f, 0, 1, nil)
	if err != nil {
		panic(err)
	}
	x := linspace(0.002, 1, 300)
	y := make([]float64, len(x))
	for i, v := range x {
		y[i] = f(v)
	}
	fig, axs := plot.Subplots(1, 1, plot.FigSize(6.4, 4))
	ax := axs[0][0]
	ax.Plot(x, y, plot.Label(`$\log x / \sqrt{x}$`))
	ax.FillBetween(x, y, nil, plot.Alpha(0.3))
	ax.SetYLim(-12, 0.5)
	ax.SetTitle(fmt.Sprintf("Quad: %.12f, error estimate %.1e, %d evaluations", r.Value, r.AbsErr, r.NEval))
	ax.SetXLabel(`$x$`)
	ax.Legend(plot.Loc("lower right"))
	return fig
}
