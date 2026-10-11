package main

import (
	"github.com/go-numerics/special"
	"github.com/go-plot/plot"
)

func init() { register("gamma", gamma) }

// Gamma on the real line: one curve per interval between its poles at
// 0, -1, -2, -3, so that no line is drawn across a pole.
func gamma() *plot.Figure {
	fig, axs := plot.Subplots(1, 1, plot.FigSize(6.4, 4))
	ax := axs[0][0]
	pieces := [][2]float64{{-4, -3}, {-3, -2}, {-2, -1}, {-1, 0}, {0, 5}}
	for k, p := range pieces {
		x := linspace(p[0]+0.005, p[1]-0.005, 300)
		if k == len(pieces)-1 {
			x = linspace(p[0]+0.005, p[1], 300)
		}
		y := make([]float64, len(x))
		for i, s := range x {
			y[i] = special.Gamma(s)
		}
		if k == len(pieces)-1 {
			ax.Plot(x, y, plot.Color("C0"), plot.Label(`$\Gamma(x)$`))
		} else {
			ax.Plot(x, y, plot.Color("C0"))
		}
	}
	x := linspace(0.05, 5, 200)
	lg := make([]float64, len(x))
	for i, s := range x {
		lg[i] = special.LogGamma(s)
	}
	ax.Plot(x, lg, plot.Color("C1"), plot.LineStyle("--"), plot.Label(`$\log\Gamma(x)$`))
	ax.Axhline(0, plot.Color("gray"), plot.LineWidth(0.8))
	ax.SetYLim(-8, 8)
	ax.SetXLim(-4, 5)
	ax.SetTitle("Gamma and LogGamma")
	ax.SetXLabel(`$x$`)
	ax.Legend(plot.Loc("lower right"))
	return fig
}
