package main

import (
	"math"

	"github.com/go-numerics/optimize"
	"github.com/go-plot/plot"
)

func init() { register("robust", robust) }

// The README's LeastSquares example, fitted twice: with the ordinary squared
// loss, which the outlier at t = 5 drags upwards, and with soft_l1, which
// keeps it from dominating. Same bounds, same start.
func robust() *plot.Figure {
	t := []float64{0, 1, 2, 3, 4, 5}
	y := []float64{1.0, 0.62, 0.36, 0.23, 0.13, 3.0} // the last point is an outlier
	p := optimize.LSProblem{M: len(t), Func: func(x, f []float64) {
		for i := range t {
			f[i] = x[0]*math.Exp(-x[1]*t[i]) - y[i]
		}
	}}
	bounds := []optimize.Bound{{Lo: 0, Hi: 2}, {Lo: 0, Hi: 0.5}}
	linear, err := optimize.LeastSquares(p, []float64{1, 0.1}, optimize.TRF,
		&optimize.LSOptions{Bounds: bounds})
	if err != nil {
		panic(err)
	}
	soft, err := optimize.LeastSquares(p, []float64{1, 0.1}, optimize.TRF,
		&optimize.LSOptions{Bounds: bounds, Loss: optimize.LossSoftL1, FScale: 0.1})
	if err != nil {
		panic(err)
	}
	xs := linspace(0, 5, 200)
	curve := func(x []float64) []float64 {
		v := make([]float64, len(xs))
		for i, s := range xs {
			v[i] = x[0] * math.Exp(-x[1]*s)
		}
		return v
	}
	fig, axs := plot.Subplots(1, 1, plot.FigSize(6.4, 4))
	ax := axs[0][0]
	ax.Plot(t, y, plot.Marker("o"), plot.LineStyle("None"), plot.Label("data"))
	ax.Plot(xs, curve(linear.X), plot.LineStyle("--"), plot.Label("squared loss"))
	ax.Plot(xs, curve(soft.X), plot.Label("soft_l1 loss, f_scale = 0.1"))
	ax.SetTitle("LeastSquares (TRF, bounds): one outlier, two losses")
	ax.SetXLabel(`$t$`)
	ax.Legend()
	return fig
}
