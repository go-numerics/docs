package main

import (
	"fmt"
	"math"

	"github.com/go-ndarray/ndarray"
	"github.com/go-numerics/optimize"
	"github.com/go-plot/plot"
)

func init() { register("curvefit", curvefit) }

// The README's CurveFit example: 50 points of 2.5 exp(-1.3x) + 0.5 with a
// deterministic ripple, fitted from p0 = (1, 1, 1). The legend prints the
// parameters and standard errors CurveFit returned.
func curvefit() *plot.Figure {
	x, _ := ndarray.Linspace(0, 4, 50)
	y := x.Map(func(v float64) float64 { return 2.5*math.Exp(-1.3*v) + 0.5 + 0.01*math.Sin(37*v) })
	model := func(x float64, p []float64) float64 { return p[0]*math.Exp(-p[1]*x) + p[2] }
	r, err := optimize.CurveFit(model, x, y, []float64{1, 1, 1}, nil)
	if err != nil {
		panic(err)
	}
	perr := r.Perr()
	xs := linspace(0, 4, 200)
	fit := make([]float64, len(xs))
	for i, v := range xs {
		fit[i] = model(v, r.Popt)
	}
	fig, axs := plot.Subplots(1, 1, plot.FigSize(6.4, 4))
	ax := axs[0][0]
	ax.Plot(x, y, plot.Marker("o"), plot.LineStyle("None"), plot.MarkerSize(4), plot.Label("data"))
	ax.Plot(xs, fit, plot.Label(fmt.Sprintf(`$a e^{-bx} + c$: a = %.3f ± %.3f, b = %.3f ± %.3f, c = %.3f ± %.3f`,
		r.Popt[0], perr[0], r.Popt[1], perr[1], r.Popt[2], perr[2])))
	ax.SetTitle(fmt.Sprintf("CurveFit (Levenberg–Marquardt): %d evaluations", r.Nfev))
	ax.SetXLabel(`$x$`)
	ax.Legend(plot.FontSize(8))
	return fig
}
