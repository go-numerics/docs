---
title: "go-numerics documentation"
linkTitle: "Home"
type: docs
cascade:
  type: docs
description: "Numerical methods in pure Go (CGO=0): optimisation and root finding, quadrature and ODE solvers, special functions — on go-ndarray."
---

**go-numerics is numerical computing in pure Go.** Three libraries port
reference algorithms statement by statement — MINPACK, L-BFGS-B, QUADPACK, the
Runge–Kutta, Radau and BDF solvers, Cephes and Boost.Math — and test them
against recorded answers: SciPy 1.18.1 runs for `optimize` and `integrate`,
with the netlib Fortran as a second oracle, and
[mpmath](https://mpmath.org) at high precision for `special`, with SciPy
measured on the same points. Where the Go code is slower than the reference,
the pages say so.

<img class="figure" src="images/figures/lotka.svg" alt="Prey and predator populations from SolveIVP with RK45, one marker per accepted step">

Every figure on these pages was computed by the library it illustrates and
drawn by [go-plot](https://github.com/go-plot/plot), from the code shown
beside it.

## The libraries

| Library | Version | What it does |
| --- | --- | --- |
| [optimize]({{< relref "optimize" >}}) | v0.3.0 | minimisation (Nelder–Mead, Powell, CG, BFGS, Newton-CG, L-BFGS-B), constrained minimisation (SLSQP with bounds, equality and inequality constraints), scalar and system roots, nonlinear least squares, curve fitting with covariance, linear programming |
| [integrate]({{< relref "integrate" >}}) | v0.1.0 | adaptive quadrature (QUADPACK), fixed and sample rules, multiple integrals, initial value problems (RK45, RK23, DOP853, Radau, BDF, stiffness switching) |
| [special]({{< relref "special" >}}) | v0.1.0 | gamma, beta, error and normal functions, incomplete gamma and beta and their inverses, Bessel, zeta, elliptic integrals; an element-wise layer over go-ndarray |

## Install

```sh
go get github.com/go-numerics/optimize@v0.1.0
go get github.com/go-numerics/integrate@v0.1.0
go get github.com/go-numerics/special@v0.1.0
```

Each module needs Go 1.27.2 or later and is pure Go (`CGO_ENABLED=0`). Each
one's CI runs its tests on amd64 and arm64 and under emulation on riscv64,
loong64, ppc64le and s390x; `integrate` and `special` add 386 and 32-bit arm.

## What each library's pages hold

| Library | Pages |
| --- | --- |
| optimize | [overview]({{< relref "optimize" >}}) · [design]({{< relref "optimize/design" >}}) · [performance]({{< relref "optimize/performance" >}}) |
| integrate | [overview]({{< relref "integrate" >}}) · [design]({{< relref "integrate/design" >}}) · [verification]({{< relref "integrate/verification" >}}) · [quadrature table]({{< relref "integrate/quad-table" >}}) · [ODE table]({{< relref "integrate/ode-table" >}}) · [performance]({{< relref "integrate/performance" >}}) |
| special | [overview]({{< relref "special" >}}) · [accuracy]({{< relref "special/accuracy" >}}) · [accuracy table]({{< relref "special/accuracy-table" >}}) · [performance]({{< relref "special/performance" >}}) · [performance table]({{< relref "special/performance-table" >}}) |

## Honest weak points

Where these libraries are slower than the reference, their own measurements
say so, and so do these pages:

- `optimize`: `curve_fit` on 10⁴ points is slower than SciPy's (0.67× with
  LM, 0.77× with TRF); the gap is in the solver's workspace handling
  ([performance]({{< relref "optimize/performance" >}})).
- `special`: single-threaded, a Go loop is about 3× slower than SciPy's C for
  the gamma and error-function families and 10–15× for the Bessel functions of
  real order and the incomplete beta function; the `nd` functions, spread over
  the cores, beat SciPy for 21 of 31 functions measured
  ([performance]({{< relref "special/performance" >}})).

## Built on

[go-ndarray](https://github.com/go-ndarray/ndarray) (arrays: curve-fit data,
covariances, LP matrices, ODE solutions and the element-wise `special/nd`
layer). The figures are drawn by [go-plot](https://github.com/go-plot/plot),
and a [go-notebook](https://github.com/go-notebook/notebook) cell can use all
of it.
