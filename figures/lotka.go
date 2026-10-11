package main

import (
	"fmt"

	"github.com/go-numerics/integrate"
	"github.com/go-plot/plot"
)

func init() { register("lotka", lotka) }

// The README's SolveIVP example: Lotka–Volterra with RK45, the default
// method, and the default tolerances. The markers are the solver's own
// accepted steps, sol.T: no interpolation.
func lotka() *plot.Figure {
	f := func(t float64, y, dydt []float64) {
		dydt[0] = 1.5*y[0] - y[0]*y[1]
		dydt[1] = -3*y[1] + y[0]*y[1]
	}
	sol, err := integrate.SolveIVP(f, [2]float64{0, 15}, []float64{10, 5}, nil)
	if err != nil {
		panic(err)
	}
	prey := make([]float64, len(sol.T))
	predators := make([]float64, len(sol.T))
	for k, y := range sol.Y {
		prey[k], predators[k] = y[0], y[1]
	}
	fig, axs := plot.Subplots(1, 1, plot.FigSize(6.4, 4))
	ax := axs[0][0]
	ax.Plot(sol.T, prey, plot.Marker("o"), plot.MarkerSize(3), plot.Label("prey"))
	ax.Plot(sol.T, predators, plot.Marker("s"), plot.MarkerSize(3), plot.Label("predators"))
	ax.SetTitle(fmt.Sprintf("SolveIVP, RK45: %d steps, %d evaluations", sol.NSteps, sol.NFev))
	ax.SetXLabel(`time $t$`)
	ax.Legend()
	return fig
}
