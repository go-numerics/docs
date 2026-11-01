---
title: "optimize"
weight: 10
description: "Minimisation, scalar and system roots, nonlinear least squares, curve fitting with covariance, and linear programming."
---

<!-- From go-numerics/optimize v0.1.0, README.md: the text is the library's, with its links pointed into this site and figures added. Change it there first. -->

Numerical optimisation and root finding in pure Go (CGO=0): minimisation,
scalar and system roots, nonlinear least squares, curve fitting with parameter
covariance, and linear programming.

The package covers the parts of `scipy.optimize` that scientific notebooks use
most. Each function documents the SciPy call it corresponds to. The test suite
replays SciPy 1.18.1 on standard problems, and it checks solutions, objective
values, status codes, messages, and iteration and evaluation counts.

```sh
go get github.com/go-numerics/optimize@v0.1.0
```

It works with [go-ndarray](https://github.com/go-ndarray/ndarray): curve-fit
data, covariances and LP constraint matrices are `*ndarray.Array`, and results
have `…Array()` accessors.

## What is there

| Go | SciPy | Algorithm |
|---|---|---|
| `Minimize(p, x0, NelderMead, opts)` | `minimize(method='Nelder-Mead')` | downhill simplex, adaptive option, bounds |
| `Minimize(…, Powell, …)` | `minimize(method='Powell')` | conjugate directions, bounds |
| `Minimize(…, CG, …)` | `minimize(method='CG')` | Polak–Ribière+, strong Wolfe |
| `Minimize(…, BFGS, …)` | `minimize(method='BFGS')` | inverse-Hessian BFGS, O(n²) update |
| `Minimize(…, NewtonCG, …)` | `minimize(method='Newton-CG')` | truncated Newton (Hess, HessP or FD) |
| `Minimize(…, LBFGSB, …)` | `minimize(method='L-BFGS-B')` | L-BFGS-B 3.0 with bounds |
| `Minimize(…, SLSQP, …)` with `Options.Constraints` | `minimize(method='SLSQP', constraints=…)`, MATLAB `fmincon` `'sqp'` | Kraft's SLSQP: bounds, equality and inequality constraints as `Constraint`, `NonlinearConstraint` or `LinearConstraint`; Lagrange multipliers in `Result.Multipliers` |
| `MinimizeScalar`, `Fminbound`, `Bracket` | `minimize_scalar` (brent, golden, bounded), `fminbound`, `bracket` | Brent, golden section, FMIN |
| `Brentq`, `Brenth`, `Ridder`, `Bisect`, `Newton`, `RootScalar` | `brentq`, `brenth`, `ridder`, `bisect`, `newton` (Newton/secant/Halley), `root_scalar` | |
| `Root`, `Fsolve` | `root(method='hybr')`, `fsolve` | MINPACK hybrd/hybrj (Powell hybrid) |
| `LeastSquares(…, TRF, …)` | `least_squares(method='trf')` | trust-region reflective, bounds, robust losses, `x_scale` |
| `LeastSquares(…, LM, …)`, `Leastsq` | `least_squares(method='lm')`, `leastsq` (with `cov_x`) | MINPACK lmder/lmdif |
| `CurveFit`, `CurveFitModel` | `curve_fit` (`sigma` 1-D/2-D, `absolute_sigma`, bounds, losses) | LM or TRF, covariance as SciPy |
| `Linprog` | `linprog` | revised simplex (default) or interior point; duals |
| `Quadprog` (**released in v0.5.0**; design: docs/design-qp.md in the optimize repo) | MATLAB `quadprog`, OSQP | the API, validation and result blocks of ½xᵀPx + qᵀx s.t. equalities, inequalities, two-sided rows and bounds, dense or sparse (go-sparse) P and A; solved today: positive definite P by the Goldfarb–Idnani dual active-set method (infeasibility certificates, dependent and redundant rows; `QPSolver` keeps the factors between solves and `UpdateQ`, `UpdateBounds`, `UpdateMatrixValues` and `WorkingSet` repair them instead of restarting, Goldfarb–Idnani Prop. 8.1–8.2: 0.03 of the cold pivots at a 5% perturbation of q, 0.30 at 50%, 0.80 at 200%, same answer); a singular P, P = 0 (LP) or sparse constraints go to the ADMM method of OSQP v1.0.0 (`QPADMM`, picked by `QPAuto`: go-sparse quasi-definite LDLᵀ with AMD, Ruiz scaling, adaptive ρ, termination with the duality gap, infeasibility and unboundedness certificates, non-convexity detection, `Polish`, and `QPSolver` updates and warm start for ADMM too). On the 38 problems of the oracle fixture OSQP can read, and on 12 generated sparse ones (random QP, equality QP, Lasso, portfolio), status, iteration count and ρ-update count equal OSQP v1.0.0's on every one. Checked against OSQP, Clarabel and quadprog on 39 problems (x to 1e-7 and 1e-8) and against the SLSQP `lsei` chain on 600 random ones. On 33 Maros–Mészáros problems (tiny and small tiers, dense or sparse fixtures, pinned commit and SHA-256) ADMM at 1e-9 with polishing reaches Clarabel's objective on 27 and ends with OSQP v1.0.0's own status on the other 6 (iteration limit on four, a dual-infeasibility certificate on two), and the active set reaches it on all 13 with a positive definite P. Against OSQP v1.0.0 (same iterations on 25 of 25 generated instances) and the `quadprog` package, one thread, medians of 5 interleaved rounds ([performance]({{< relref "performance" >}})): ADMM takes 0.5–0.65× OSQP's time on dense n ≥ 500 and 0.8–0.9× on the larger portfolios, but 1.04–1.49× in the solve phase and 2.2–4.9× in the setup on small sparse problems (go-sparse's AMD, symbolic analysis and LDLᵀ solves: open defect there); the active set takes 0.68–0.76× the `quadprog` package's time from n = 500, 1.09× at n = 200, 2.0× at n = 50. `QPAuto` picks the active set up to n + m = 2500 |
| `Brute` | `brute` | grid then Nelder–Mead polish |
| `DifferentialEvolution` (**on main, release planned**) | `differential_evolution` | the twelve strategies, immediate and deferred updating, `Workers` (identical results for every worker count with deferred updating), dithering, LHS / random / array init, `X0`, custom strategy, callback, `Constraints` (Lampinen's rule), `Integrality`, L-BFGS-B or trust-constr polish; follows SciPy's random stream exactly (`NewRNG(seed)` is `rng=seed`, `NewLegacyRNG(seed)` is `seed=seed`); no allocation per generation |
| `Basinhopping` (**on main, release planned**) | `basinhopping` | Metropolis acceptance (`T = 0` is monotonic), adaptive step size, custom `TakeStep`, both accept-test forms including `ForceAccept`, callback, `NIterSuccess`; any `Minimize` method as the local minimiser; follows SciPy's random stream (displacements and Metropolis draws on one source) |
| `Direct` (**on main, release planned**) | `direct` | DIRECT and DIRECT-L (`Original` selects Jones' version), a statement-by-statement port of the C code SciPy wraps (MIT licence, NOTICE in the repository); deterministic, same points, `Nfev`, `Nit` and status as SciPy; `Workers` evaluates the points of one rectangle concurrently with identical results |
| `CheckGrad`, `ApproxFprime` | `check_grad`, `approx_fprime` | forward differences |

Gradients are optional everywhere. Without one, finite differences use the
step rules of SciPy's `approx_derivative`, kept inside the bounds.
Non-convergence is reported in `Success`/`Status`/`Message` as SciPy does.
Invalid input returns an error.

## Examples

Each example below is a runnable, output-checked Go example in
`example_test.go`, and its output was checked against SciPy.

```go
// scipy.optimize.minimize(rosen, x0, jac=rosen_der, method='BFGS')
r, err := optimize.Minimize(optimize.Problem{Func: rosen, Grad: rosenDer},
	[]float64{-1.2, 1}, optimize.BFGS, nil)
// r.X = [1 1], r.Nit = 32 (SciPy: 32), "Optimization terminated successfully."

// scipy.optimize.minimize(f, x0, method='L-BFGS-B', bounds=[(-2, 0.5), (None, None)])
r, _ = optimize.Minimize(optimize.Problem{Func: f}, []float64{0, 0}, optimize.LBFGSB,
	&optimize.Options{Bounds: []optimize.Bound{{Lo: -2, Hi: 0.5}, optimize.Unbounded()}})

// scipy.optimize.minimize(f, x0, method='SLSQP', bounds=[(1, 5)]*4,
//     constraints=[{'type': 'ineq', 'fun': g}, {'type': 'eq', 'fun': h}])     (Hock–Schittkowski 71)
r, _ = optimize.Minimize(optimize.Problem{Func: f}, []float64{1, 5, 5, 1}, optimize.SLSQP,
	&optimize.Options{Bounds: bounds, Constraints: []optimize.ConstraintSpec{
		optimize.Constraint{Kind: optimize.Ineq, M: 1, Func: g},  // x1·x2·x3·x4 − 25 ≥ 0
		optimize.Constraint{Kind: optimize.Eq, M: 1, Func: h},    // Σ xᵢ² − 40 = 0
	}})
// r.Fun = 17.0140172, r.Nit = 5 (SciPy: 5), r.Multipliers

// The same problem by the interior-point method: minimize(..., method='trust-constr')
r, _ = optimize.Minimize(optimize.Problem{Func: f}, []float64{1, 5, 5, 1}, optimize.TrustConstr,
	&optimize.Options{Bounds: bounds, Constraints: cons})
// r.Fun = 17.01402, r.Success, r.Optimality, r.ConstrViolation

// popt, pcov = scipy.optimize.curve_fit(lambda x, a, b, c: a*np.exp(-b*x) + c, x, y)
fit, err := optimize.CurveFit(func(x float64, p []float64) float64 {
	return p[0]*math.Exp(-p[1]*x) + p[2]
}, xdata, ydata, []float64{1, 1, 1}, nil)
// fit.Popt, fit.PcovArray(), fit.Perr()  ->  a = 2.500 ± 0.004, b = 1.300 ± 0.005, c = 0.500 ± 0.002

// scipy.optimize.brentq(lambda x: x**3 - 2*x - 5, 2, 3)
root, _ := optimize.Brentq(func(x float64) float64 { return x*x*x - 2*x - 5 }, 2, 3, nil)
// root.Root = 2.094551481542…, 7 iterations (as SciPy)

// scipy.optimize.root(fun, [0, 0])        (fun writes f(x) into its second argument)
sol, _ := optimize.Root(optimize.SystemProblem{Func: fun}, []float64{0, 0}, nil)

// scipy.optimize.least_squares(fun, x0, bounds=([0, 0], [2, 0.5]), loss='soft_l1', f_scale=0.1)
ls, _ := optimize.LeastSquares(optimize.LSProblem{Func: resid, M: m}, []float64{1, 0.1}, optimize.TRF,
	&optimize.LSOptions{Bounds: []optimize.Bound{{0, 2}, {0, 0.5}}, Loss: optimize.LossSoftL1, FScale: 0.1})

// scipy.optimize.linprog(c, A_ub=A, b_ub=b, bounds=[(None, None), (-3, None)])
lp, _ := optimize.Linprog([]float64{-1, 4}, optimize.LPConstraints{AUb: A, BUb: b,
	Bounds: []optimize.Bound{optimize.Unbounded(), {Lo: -3, Hi: math.Inf(1)}}}, nil)
// lp.X = [10 -3], lp.Fun = -22, lp.IneqLin.Marginals (HiGHS sign convention)
```

Callbacks receive the point first and write their outputs into the slices that
follow: `Func(x)`, `Grad(x, g)`, `FuncGrad(x, g)` (SciPy's `jac=True`),
`Hess(x, h)`, `HessP(x, p, hp)`, residuals `Func(x, f)`, Jacobians `Jac(x, J)`.
Matrices are row-major `[]float64`.

The README's curve fit and robust least-squares examples, drawn with
[go-plot](https://github.com/go-plot/plot) from the results they return:

{{< plotfig "curvefit" "Fifty points of a decaying exponential and the curve CurveFit fitted to them" >}}

{{< plotfig "robust" "Six points with one outlier, fitted with the squared loss and with soft_l1" >}}

## How close to SciPy

All numbers below are against SciPy 1.18.1 fixtures generated by
`testdata/gen/*.py` (run with `uv run --with scipy --with numpy python …`) and
embedded in the tests. They are the same on arm64 and amd64, except where
noted.

| Area | Fixtures | Result |
|---|---|---|
| `Minimize` (NM, Powell, CG, BFGS, Newton-CG, L-BFGS-B) | 79 runs on Rosenbrock n=2/10/100, Beale, Himmelblau, Powell singular, Wood, with and without gradients, bounds, iteration limits | every converged SciPy run is matched with an objective at least as good, within the method's own tolerance; on the arithmetic-only problems, **36/38 runs take SciPy's exact path** (identical nit, nfev, njev); the 2 exceptions are explained in docs/design.md |
| L-BFGS-B | 44 runs, including 8 problems that reach rare branches (Cholesky failures, the 3.0 backtrack, Cauchy breakpoints) | identical nit/nfev/njev/status/message on every analytic-gradient run except Rosenbrock n=100 (513 vs 514 iterations, SciPy's Accelerate BLAS rounding); bit-identical to SciPy's C code built with sequential BLAS |
| `Minimize(SLSQP)` | 100 runs: 36 Hock–Schittkowski problems and 10 stress cases (infeasible, inconsistent and singular linearisations, fixed variables), with and without derivatives, `NonlinearConstraint`/`LinearConstraint` and option variants; SciPy's x86_64 wheel | 100/100 statuses; **88/100 on SciPy's exact path** (nit/nfev/njev, x to 1e-7, multipliers to 1e-6); the other 12 are runs on which SciPy's own arm64 and x86_64 wheels also disagree |
| `Minimize(TrustConstr)` | the same 100 runs | success/failure as SciPy on 100/100 (2 are errors in both: more equalities than variables); identical status on 93/100 on arm64, 94/100 on amd64 (the rest are gtol against xtol, both successes); objective within 1e-5 of SciPy or 1e-4 of the known optimum on every success; **65/100 on SciPy's exact path on arm64, 69/100 on amd64** — SciPy's own wheels disagree with each other on 35 runs, and the iterates of the 4 other runs that leave the path were traced to the last bit (design page) |
| `Root`/`Fsolve` | 23 runs (MGH problems with and without Jacobian, failure and option cases) | 23/23 identical counts, status and message; x bit-identical |
| `LeastSquares(LM)` / `Leastsq` | 39 / 44 runs on the MGH least-squares set | 36/39 and 41/44 identical counts; the mismatches are on problems that call exp/sin/cos, where Go's libm and the C library differ in the last ulp |
| `LeastSquares(TRF)` | 89 runs: MGH problems, bounds, all robust losses, `x_scale`, 2-/3-point differences | 86/89 identical nfev/njev/status (84/89 on amd64); every cost and active mask matches |
| `CurveFit` | 11 fits (exponentials, Gaussians, sum of two Gaussians, a 2-variable plane, 1-D and 2-D `sigma`, bounds, `absolute_sigma`, m = n) | popt within 1e-6, pcov within 1e-5 relative, identical nfev on every fit |
| Scalar roots | 71 runs on 11 hard functions (multiple roots, flat, steep, kinks, poles) × 7 methods | 68/71 on SciPy's exact path (same root, iterations, calls); the rest differ only through `cos`/`exp` in the last ulp |
| `MinimizeScalar`, `Bracket` | 36 + 6 runs | identical x/nit/nfev/message on arithmetic-only functions |
| `Linprog` | 46 problems vs HiGHS: Klee–Minty, cycling examples, degenerate, free and upper-bounded variables, redundant rows, infeasible, unbounded, random up to 60×40 | 46/46 statuses; objective error ≤ 1.8e-15 (simplex) and ≤ 2e-8 (interior point); solutions and duals inside SciPy's optimal sets |
| `Brute` | 4 grids (the SciPy documentation example, Himmelblau, Rosenbrock 3-D) | x and f to 1e-9 |
| `DifferentialEvolution` (on main) | 378 runs against SciPy 1.18.1 (linux/amd64) on 10 arithmetic-only objectives, both seeding forms (`NewRNG`: `rng=int`, PCG64; `NewLegacyRNG`: `seed=int`, MT19937): the 12 strategies, immediate and deferred updating, LHS / random / array init, `X0`, dither, tolerances, `MaxFev`, `MaxIter`, `PopSize`, linear / polynomial / equality / impossible constraints, integer variables | unpolished (296 runs): `nfev`, `nit`, message, x, fun, the final population and its energies bit for bit on every run, i.e. SciPy's exact random path; polished (82 runs): 56 bit-identical, the rest differ only through the polisher (L-BFGS-B or trust-constr, last digits; trust-constr values agree within 1e-3) |
| `Basinhopping` (on main) | 224 runs against SciPy 1.18.1 (linux/amd64): four polynomial landscapes with analytic gradients, BFGS and L-BFGS-B, both seeding forms, several temperatures, step intervals, `NIterSuccess`, both accept-test forms with `ForceAccept` | every run makes SciPy's decisions (same accept sequence at every step, same `nit`, failures and message, values within 1e-8); 95 are bit-identical in the whole trace, the others differ in the last digits of the local minimiser's final iterate (open defect in the BFGS / L-BFGS-B ports) |
| `Direct` (on main) | 770 runs against SciPy 1.18.1 (linux/amd64): 27 objectives, DIRECT and DIRECT-L, fixed and Jones-updated `eps`, `len_tol`/`vol_tol`/`f_min` stops, `maxfun` and `maxiter` limits, NaN and ±Inf regions | on the 618 arithmetic-only runs: status, message, `nfev`, `nit`, x and fun bit for bit, and the first 25 evaluated points; on the transcendental ones: same status and counts on all, fun bit-identical on most; the SciPy documentation example gives `nfev` 2011 and 207 |

How such close agreement is possible:

- The ports follow their references statement by statement.
- Every product that feeds an addition is rounded explicitly, because Go
  would otherwise fuse it into an FMA on arm64 and other architectures.
- The tests use FMA-free oracles.

See [Design]({{< relref "optimize/design" >}}).

**Negative controls.** One algorithmic step was broken at a time; each break
made the suite fail, and each was then reverted:

| Broken step | Failing tests |
|---|---|
| A term dropped from the BFGS update | 6 |
| A sign flipped in Moré–Thuente's `dcstep` | 57 |
| The wrong Brent golden constant | 4 |
| The Levenberg correction dropped in `lmpar` | 40 |
| A sign flipped in `dogleg` | 17 |
| TRF reflection disabled, or its radius growth changed | 20 / 48 of the 89 cases leave SciPy's path |
| Max instead of min in the simplex ratio test | 83 |
| The sign of the interior-point corrector flipped | iteration parity drops from 43 to 2 of 46 |
| L-BFGS-B's subspace backtrack skipped, its Cauchy f2 update dropped, or its breakpoint heap disabled | named fixtures fail for each |

## Performance

Medians of 7 rounds, with Go and SciPy interleaved on a loaded host.
[Performance]({{< relref "optimize/performance" >}}) has the method and caveats.

| Benchmark | Go | SciPy | Speed-up |
|---|---|---|---|
| Rosenbrock n=2, BFGS | 5.7 µs | 2.58 ms | 449× |
| Rosenbrock n=100, BFGS | 14.8 ms | 87.9 ms | 5.9× |
| Rosenbrock n=2, L-BFGS-B | 126 µs | 1.47 ms | 11.7× |
| Rosenbrock n=100, L-BFGS-B | 6.36 ms | 23.2 ms | 3.7× |
| `curve_fit`, 10⁴ points, LM | 4.13 ms | 2.78 ms | **0.67× (SciPy faster)** |
| `curve_fit`, 10⁴ points, TRF | 5.97 ms | 4.61 ms | **0.77× (SciPy faster)** |
| `brentq` | 0.27 µs | 8.6 µs | 32× |
| SLSQP, HS71 (analytic / finite differences) | 5.6 / 6.9 µs | 114 / 648 µs | 20× / 94× |
| SLSQP, n = 50 with 100 bounds | 3.35 ms | 4.59 ms | 1.4× |
| trust-constr, HS71 (analytic / finite differences) | 496 / 438 µs | 32.7 / 34.1 ms | 66× / 78× |
| trust-constr, n = 50 with 100 bounds | 27.8 ms | 69.9 ms | 2.5× |

Most of SciPy's time on small problems is Python overhead around the
callbacks. Even after subtracting SciPy's time spent in its own callbacks, the
Go algorithms remain 2–300× faster, depending on the problem.

On large vectorised curve fits SciPy is faster today. Its model call is one
NumPy expression, and its MINPACK/LAPACK kernels are compiled C. Our model
evaluation is at parity; the gap is in the solver's workspace handling.

## Not done

- **`minimize`:** COBYLA/COBYQA, TNC, dogleg/trust-ncg and the other
  trust-region methods; trust-constr's `keep_feasible`, sparse Jacobians and
  SR1 Hessian option; QP solvers (active set, ADMM) are planned next and will
  share the constraint types. For L-BFGS-B, the `hess_inv` operator is not
  provided.
- **`root` / `least_squares`:** `root` methods other than `hybr`; the
  `band`/`col_deriv` options; dogbox; the lsmr trust-region solver and sparse
  Jacobians.
- **`linprog`:** HiGHS-class performance, MILP, singleton-row presolve,
  crossover and sparse matrices.
- **Global optimisers:** shgo, dual_annealing, DE `init='sobol'` and
  `'halton'` (`Direct`, `DifferentialEvolution` and `Basinhopping` are on main,
  release planned).

## Licence

BSD 3-Clause; see [LICENSE](https://github.com/go-numerics/optimize/blob/v0.1.0/LICENSE).
