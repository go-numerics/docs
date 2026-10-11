package main

import (
	"github.com/go-numerics/special"
	"github.com/go-plot/plot"
)

func init() { register("ndtr", ndtr) }

// The lower tail of the normal distribution and the complementary error
// function, on a log scale, down to the subnormal range of float64.
func ndtr() *plot.Figure {
	x := linspace(-38, 0, 400)
	p := make([]float64, len(x))
	for i, s := range x {
		p[i] = special.Ndtr(s)
	}
	z := linspace(0, 26.5, 400)
	q := make([]float64, len(z))
	for i, s := range z {
		q[i] = special.Erfc(s)
	}
	fig, axs := plot.Subplots(1, 2, plot.FigSize(8, 3.6))
	axs[0][0].Plot(x, p)
	axs[0][0].SetYScale("log")
	axs[0][0].SetYLim(1e-320, 10)
	axs[0][0].SetTitle(`Ndtr($x$)`)
	axs[0][0].SetXLabel(`$x$`)
	axs[0][1].Plot(z, q, plot.Color("C1"))
	axs[0][1].SetYScale("log")
	axs[0][1].SetYLim(1e-320, 10)
	axs[0][1].SetTitle(`Erfc($x$)`)
	axs[0][1].SetXLabel(`$x$`)
	fig.TightLayout()
	return fig
}
