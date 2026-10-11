package main

import (
	"fmt"

	"github.com/go-numerics/integrate"
	"github.com/go-plot/plot"
)

func init() { register("vanderpol", vanderpol) }

// The README's stiff example: van der Pol with mu = 1000 over [0, 3000],
// solved with Radau. An explicit method needs millions of steps here.
func vanderpol() *plot.Figure {
	const mu = 1000
	f := func(t float64, y, dydt []float64) {
		dydt[0] = y[1]
		dydt[1] = mu*(1-y[0]*y[0])*y[1] - y[0]
	}
	sol, err := integrate.SolveIVP(f, [2]float64{0, 3000}, []float64{2, 0},
		&integrate.IVPOptions{Method: integrate.Radau})
	if err != nil {
		panic(err)
	}
	y0 := make([]float64, len(sol.T))
	for k, y := range sol.Y {
		y0[k] = y[0]
	}
	fig, axs := plot.Subplots(1, 1, plot.FigSize(6.4, 3.6))
	ax := axs[0][0]
	ax.Plot(sol.T, y0, plot.Marker("o"), plot.MarkerSize(2))
	ax.SetTitle(fmt.Sprintf(`Radau, $\mu = 1000$: %d steps, %d LU decompositions`, sol.NSteps, sol.NLU))
	ax.SetXLabel(`time $t$`)
	ax.SetYLabel(`$y_0$`)
	return fig
}
