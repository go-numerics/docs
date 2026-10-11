package main

import (
	"fmt"

	"github.com/go-numerics/special"
	"github.com/go-plot/plot"
)

func init() { register("bessel", bessel) }

// Bessel functions of the first kind, integer and real order, with Jv.
func bessel() *plot.Figure {
	x := linspace(0, 20, 400)
	fig, axs := plot.Subplots(1, 1, plot.FigSize(6.4, 4))
	ax := axs[0][0]
	for _, v := range []float64{0, 1, 2, 2.5} {
		y := make([]float64, len(x))
		for i, s := range x {
			y[i] = special.Jv(v, s)
		}
		ax.Plot(x, y, plot.Label(fmt.Sprintf(`$J_{%g}(x)$`, v)))
	}
	ax.Axhline(0, plot.Color("gray"), plot.LineWidth(0.8))
	ax.SetTitle("Jv: Bessel functions of the first kind")
	ax.SetXLabel(`$x$`)
	ax.Legend()
	return fig
}
