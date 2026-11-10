---
title: "Performance"
weight: 30
description: "Timings against SciPy 1.18.1, including the curve fits where SciPy is faster, and how they were taken."
---

<!-- From go-numerics/optimize v0.1.0, docs/perf.md: the text is the library's, with its links pointed into this site. Change it there first. -->

## Current timings (refbench, v0.8.0, 2026-11-10)

Measured with optimize's `bench/` (go-scikit/refbench): one solver call per
row on both sides, one core, SciPy 1.18.1, interleaved rounds per row. The
1-minute load stayed between 5.09 and 6.59, under refbench's gate of 8. Same
data (SHA-256) and the same `nfev`/`nit` as SciPy are checked per row. The
full table is optimize's `docs/perf-table.md`.

| Benchmark | Go | SciPy | Speed-up |
|---|---:|---:|---:|
| Rosenbrock n=2, BFGS, analytic gradient | 3.70 µs | 1.24 ms | 334× |
| Rosenbrock n=100, BFGS, analytic gradient (O(n²) update, tolerance-validated) | 5.74 ms | 48.1 ms | **8.4×** |
| Rosenbrock n=100, BFGS, finite differences | 9.73 ms | 572 ms | 59× |
| Rosenbrock n=2, L-BFGS-B, analytic gradient | 74.0 µs | 745 µs | 10.1× |
| Rosenbrock n=100, L-BFGS-B, analytic gradient | 4.81 ms | 12.4 ms | 2.6× |
| Rosenbrock n=100, CG, analytic gradient | 1.29 ms | 53.4 ms | 41× |
| Rosenbrock n=100, Newton-CG, Hessian-vector product | 765 µs | 26.2 ms | 34× |
| `curve_fit`, 10⁴ points, LM (model per point / `Model.Vec`) | 1.34 / 0.90 ms | 1.41 / 1.43 ms | 1.06× / 1.59× |
| `curve_fit`, 10⁴ points, TRF (model per point / `Model.Vec`) | 1.91 / 1.75 ms | 2.37 / 2.44 ms | 1.24× / 1.39× (tie, wide spread) |
| `curve_fit`, 10⁴ points, TRF with bounds | 2.57 / 1.88 ms | 3.61 / 3.66 ms | 1.40× / 1.95× |
| `brentq`, x³−2x−5 on [2, 3] | 0.13 µs | 3.9 µs | 29× |

- **No row is slower than SciPy.** BFGS at n = 100 was (optimize#53) while
  its update followed SciPy's bit for bit: two dense n×n products per
  iteration, which SciPy runs through BLAS `dgemm`. Since v0.8.0 the default
  update is the O(n²) rank-two form (5.74 ms against 48.1 ms), tolerance-validated
  and not bit-identical to SciPy (542 iterations against 544).
  `Options.SciPyExactUpdate` keeps the bit-exact path, about 90 ms at n = 100.
- **`curve_fit`:** see the profile below. With the model written over arrays
  (`Model.Vec`, go-simd/floats' vector `Exp`), the fits are 1.6–2.0× faster
  than SciPy.

The table below is the earlier script's run of 2026-10-10, kept for the
analysis that goes with it.

## Method (2026-10-10 run)

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
- **BFGS at n = 100:** in this run our update was O(n²) where SciPy does two
  dense n×n matrix products per iteration. That gave the 5× left once SciPy's
  callbacks are removed. A bit-exact update replaced it for a while and was
  slower than SciPy; since v0.8.0 the O(n²) update is the default again (see
  the current timings above).
- **L-BFGS-B at n = 100:** SciPy's core is compiled C with BLAS, and ours is
  Go. With callbacks counted we are 3.7× faster. Without them the cores are
  within 2× of each other, ours ahead.
- **Finite differences:** each forward-difference gradient costs n Python
  calls in SciPy and n inlined Go calls here. At n = 100 that is 73×.
- **`curve_fit` with 10⁴ points: SciPy was faster in the table above.**
  - SciPy's model call is one vectorised NumPy expression per evaluation, so
    it pays no per-point Python overhead. Its MINPACK and LAPACK kernels are
    compiled C.
  - These two rows were measured before the MINPACK kernel work of optimize
    PR #7, after which the 2026-10-17 run gave 1.07× (LM) and 1.27× (TRF);
    they have not been re-timed since, because the host's load has stayed
    above refbench's limit of 8.
  - Both sides take the same path: same data (SHA-256), same `p0` and
    tolerances, the same forward-difference Jacobian and the same counts
    (21 evaluations for LM, 6 for TRF).
  - A profile (2026-11-09) puts 65% of LM and 54% of TRF in the model
    sweeps: `CurveFit` calls a Go closure once per point (about 3.2 ns per
    point with `math.Exp`), where SciPy calls its NumPy expression once per
    sweep.
  - `Model.Vec` (v0.7.0) takes the model over arrays, called once per sweep,
    so it can use go-simd/floats' vector `Exp` (about 1.35 ns per element).
    With the same arithmetic per point it gives the fit of the per-point
    model to the last bit; with the vector `Exp` (0.511 ULP) the fit agrees
    with SciPy within the fixture tolerance and takes the same evaluations.
    Timed on 2026-11-09: 1.59× (LM) and 1.74× (TRF) SciPy's speed, against
    1.06× and 1.24× with the model per point.
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

## trust-constr

Same harness and rounds (`tools/slsqpbench`), 2026-10-31, medians of 7
interleaved rounds, one thread each, Apple M4 Max, 2026-10-31, v0.4.1; the
1-minute load was 7.1 before and 7.3 after, and 6.4–7.7 after each round,
below the tool's limit of 8. An earlier run that began at 7.3 and ended at
21.9 was discarded; `slsqpbench` now discards such runs itself.

| Benchmark | Go | SciPy 1.18.1 | Speed-up |
|---|---|---|---|
| HS71, analytic derivatives | 496 µs | 32.7 ms | 66.0× |
| HS71, finite differences | 438 µs | 34.1 ms | 77.8× |
| HS38 (bounds) | 277 µs | 32.1 ms | 116× |
| HS43 (3 nonlinear inequalities) | 62.4 µs | 7.15 ms | 115× |
| Quad50 (n = 50, 100 bounds, 1 linear) | 27.8 ms | 69.9 ms | 2.5× |

SciPy's trust-constr is Python around NumPy/LAPACK calls; at n = 50 each
iteration factorises the 101×151 barrier Jacobian, by pivoted QR with
dlaqp2's norm downdating and go-simd reflectors here.

## Quadprog

`tools/qpbench` times `Quadprog` against OSQP v1.0.0 (the C core of the `osqp`
1.1.3 wheel) and the `quadprog` 0.1.14 package on generated instances of the
`osqp_benchmarks` classes (random QP, equality QP, portfolio, Lasso, Huber, SVM,
control) and on dense positive definite problems. Apple M4 Max, 2026-11-01,
medians of 5 interleaved rounds, one thread, 1-minute load 2.5 before and 2.4
after, OSQP's own setup and solve timers, ε = 1e-3. Iteration counts equal
OSQP's on all 25 instances; ratio = ours / OSQP.

| Instance | n | m | setup ours | OSQP | ratio | solve ours | OSQP | ratio |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| eq-qp-n50 | 50 | 25 | 212 µs | 45.2 µs | 4.69 | 96.5 µs | 93.1 µs | 1.04 |
| random-qp-n100 | 100 | 1000 | 9.58 ms | 4.28 ms | 2.24 | 16.1 ms | 14.7 ms | 1.09 |
| lasso-n10 | 1020 | 1020 | 898 µs | 275 µs | 3.26 | 1.33 ms | 1 ms | 1.33 |
| huber-n5 | 1505 | 1500 | 834 µs | 262 µs | 3.18 | 2.07 ms | 1.39 ms | 1.49 |
| portfolio-k20 | 2020 | 2021 | 3.6 ms | 1.14 ms | 3.16 | 20.6 ms | 26.5 ms | **0.78** |
| dense-n500 | 500 | 1050 | 145 ms | 222 ms | **0.65** | 90 ms | 166 ms | **0.54** |
| dense-n1000 | 1000 | 2100 | 832 ms | 1.67 s | **0.50** | 511 ms | 918 ms | **0.56** |

On the sparse problems the extra time is go-sparse's: AMD with the symbolic
analysis is 17–80% of our setup and the LDLᵀ solves 29–73% of our solve phase,
more than OSQP's whole solve phase on `huber-n5` (1.0 ms of 2.07 ms against
1.39 ms). That is an open defect in go-sparse, written down in the repository's
`docs/perf.md` with the full 25-row table, and not worked around here.

The active set (Goldfarb–Idnani) against the `quadprog` package, setup + solve:

| n | m | ours | quadprog | ratio |
| --- | --- | --- | --- | --- |
| 10 | 21 | 7.98 µs | 1.75 µs | 4.56 |
| 50 | 105 | 270 µs | 134 µs | 2.01 |
| 200 | 420 | 12.2 ms | 11.3 ms | 1.09 |
| 500 | 1050 | 177 ms | 231 ms | **0.76** |
| 1000 | 2100 | 1.46 s | 2.15 s | **0.68** |

`QPAuto` chooses the active set up to n + m = 2500: against our own ADMM at
1e-3 the ratio is 0.75 at n + m = 1550 and 1.09 at 3100, so it crosses 1 near
2600. Reproduce with `cd tools/qpbench && GOWORK=off go run . --repo ../.. --rounds 5`.

## Reproducing

```sh
benchmarks/run.sh 7          # needs uv (https://docs.astral.sh/uv/) for SciPy
GOWORK=off go test -run '^$' -bench . -benchtime=1s   # Go side only
```
