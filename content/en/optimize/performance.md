---
title: "Performance"
weight: 30
description: "Timings against SciPy 1.18.1, including the curve fits where SciPy is faster, and how they were taken."
---

<!-- From go-numerics/optimize v0.1.0, docs/perf.md: the text is the library's, with its links pointed into this site. Change it there first. -->

## Method

The benchmarks pair `bench_test.go` with `benchmarks/bench_scipy.py`, which
runs the same problems from the same starting points. `benchmarks/run.sh`
alternates the two: each round runs every Go benchmark once
(`-benchtime=300ms`), then the whole SciPy script (each case a median of 3
batches of at least 0.2 s). Interleaving spreads the load changes of a shared
host evenly over both sides. `benchmarks/summarize.py` reports the median
over rounds. The raw numbers are in `benchmarks/results.json`.

Setup of the run below:

- 7 rounds on 2026-10-10 on an Apple Silicon (arm64) Mac, Go 1.27.2, SciPy
  1.18.1, NumPy 2.5.4, CPython 3.12.
- The host is shared: the load average was 13–19 during the run (60+ earlier
  that day). Absolute times are therefore pessimistic for both sides, and
  ratios move by about ±20% from run to run.
- **The SciPy objectives are SciPy's own vectorised `rosen`/`rosen_der`,
  and the curve-fit model is a vectorised NumPy expression.** That is the
  fastest form a notebook would write, not a straw man.

## Results (median of 7 interleaved rounds)

| Benchmark | Go | SciPy | SciPy time in Python callbacks | Speed-up | Speed-up on SciPy's time minus its callbacks |
|---|---|---|---|---|---|
| Rosenbrock n=2, BFGS, analytic gradient | 5.7 µs | 2.58 ms | 0.75 ms | **449×** | 319× |
| Rosenbrock n=100, BFGS, analytic gradient | 14.8 ms | 87.9 ms | 12.4 ms | **5.9×** | 5.1× |
| Rosenbrock n=2, L-BFGS-B, analytic gradient | 126 µs | 1.47 ms | 0.77 ms | **11.7×** | 5.6× |
| Rosenbrock n=100, L-BFGS-B, analytic gradient | 6.36 ms | 23.2 ms | 11.0 ms | **3.7×** | 1.9× |
| Rosenbrock n=2, BFGS, finite differences | 18.0 µs | 4.35 ms | 1.19 ms | **242×** | 176× |
| Rosenbrock n=100, L-BFGS-B, finite differences (stops at 15000 evaluations) | 3.31 ms | 242 ms | 133 ms | **73×** | 33× |
| Rosenbrock n=2, Nelder–Mead | 16.2 µs | 2.83 ms | 1.44 ms | **175×** | 86× |
| `curve_fit`, 10⁴ points, a·e^(−bx)+c, LM | 4.13 ms | 2.78 ms | (vectorised) | **0.67×** | — |
| `curve_fit`, 10⁴ points, TRF | 5.97 ms | 4.61 ms | (vectorised) | **0.77×** | — |
| `brentq`, x³−2x−5 on [2, 3] | 0.27 µs | 8.6 µs | — | **32×** | — |

How to read the columns:

- **"SciPy time in Python callbacks"** is the time to make SciPy's own
  `nfev` objective calls and `njev` gradient calls, outside the optimiser.
- **The last column** divides SciPy's time minus those callbacks by ours.
  It is the honest comparison of the algorithms' own code, because Go's
  callbacks cost almost nothing.

## What the numbers say

- **Small problems (n = 2):** SciPy's time is almost all fixed Python
  overhead (argument handling, `ScalarFunction` wrappers, NumPy calls on
  2-element arrays) plus about 10 µs per callback. Go spends nanoseconds
  on each. The gap is 100–450× and stays above 85× even after crediting
  SciPy with free callbacks.
- **BFGS at n = 100:** our update is O(n²) where SciPy does two dense n×n
  matrix products per iteration. That is the 5× left once SciPy's callbacks
  are removed.
- **L-BFGS-B at n = 100:** SciPy's core is compiled C with BLAS, and ours is
  Go. With callbacks counted we are 3.7× faster. Without them the cores are
  within 2× of each other, ours ahead.
- **Finite differences:** each forward-difference gradient costs n Python
  calls in SciPy and n inlined Go calls here. At n = 100 that is 73×.
- **`curve_fit` with 10⁴ points: SciPy is faster here, and we say so.**
  - SciPy's model call is one vectorised NumPy expression per evaluation, so
    it pays no per-point Python overhead. Its MINPACK and LAPACK kernels are
    compiled C.
  - Our per-point model evaluation is at parity with NumPy: 58 µs against
    56 µs per sweep of 10⁴ exponentials, measured separately.
  - The remaining gap is the solver itself on m×n = 10⁴×3 work arrays. The
    MINPACK port keeps MINPACK's reference operation order and its
    three-accumulator norm for bit parity, and it allocates its m-sized
    workspace on every call.
  - `CurveFitOptions.Workers` (parallel model evaluation) did not change the
    picture: the model accounts for only about a third of the time.
  - Closing the gap means a fit-specific LM that reuses its buffers and uses
    blocked norms. That is not done.
- **`brentq`:** 32×, almost entirely because the 8 callbacks are Go calls.

## SLSQP

Measured with `tools/slsqpbench` (Go): interleaved rounds of
`BenchmarkSLSQP*` (`-test.cpu 1`) and SciPy's `minimize(method='SLSQP')` on
the same problems with one BLAS thread; it refuses to start above a 1-minute
load of 8. Apple M4 Max, SciPy 1.18.1, 2026-10-30, medians of 7 rounds, load
4.6 before and 5.4 after.

| Benchmark | Go | SciPy 1.18.1 | Speed-up |
|---|---|---|---|
| HS71, analytic derivatives | 5.59 µs | 114 µs | 20.4× |
| HS71, finite differences | 6.90 µs | 648 µs | 93.9× |
| HS38 (bounds, 76 iterations) | 37.0 µs | 936 µs | 25.3× |
| HS43 (3 nonlinear inequalities) | 6.83 µs | 163 µs | 23.8× |
| Quad50 (n = 50, 100 bounds, 1 linear) | 3.35 ms | 4.59 ms | 1.4× |

At n = 50 both run the same O(m·n²) subproblem per iteration in compiled
code. `trsm` dominated the Go profile and the row was at parity (1.0×)
until it was unrolled to four columns per pass, bit-identical to the
reference loop.

## Reproducing

```sh
benchmarks/run.sh 7          # needs uv (https://docs.astral.sh/uv/) for SciPy
GOWORK=off go test -run '^$' -bench . -benchtime=1s   # Go side only
```
