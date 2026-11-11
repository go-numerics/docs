---
title: "Complex functions"
weight: 20
description: "special/cmplx: a drop-in superset of math/cmplx with C11 Annex G special values, each part measured in ulps against mpmath."
---

<!-- From go-numerics/special v0.8.0, README.md "Complex functions": change it there first. -->

`github.com/go-numerics/special/cmplx` (new in v0.4.0) is a drop-in
superset of the standard library's `math/cmplx`. Changing the import path
is enough: the names and signatures are the same, with C11 Annex G
special values.

- **New functions:** `Mul` and `Div` (the Annex G product and quotient),
  `Proj`, `Log2`, `Log1p` and `Expm1`.
- **Unchanged from `math/cmplx`:** `Abs`, `Phase`, `Conj`, `Polar`,
  `Rect`, `Inf`, `NaN`, `IsInf` and `IsNaN` are one call to a correctly
  specified real function, and are re-exported as they are.
- **Rewritten:** every other function. The algorithms are FreeBSD msun's
  (as NumPy's npymath carries them), carried in double-double where
  float64 loses accuracy.
- **Underlying real functions:** `sin`, `cos` and `tan` use the package's
  Payne–Hanek reduction, and `atan2` is table-driven double-double.

```
go get github.com/go-numerics/special@v0.8.0
```

```go
import "github.com/go-numerics/special/cmplx"

cmplx.Tanh(400 + 1i)   // (1+0i); math/cmplx: NaN
cmplx.Log(1 + 1e-10i)  // (5e-21+1e-10i); math/cmplx: real part 0
cmplx.Pow(1i, 2)       // (-1+0i), exactly; math/cmplx: (-1+1.2e-16i)
```

## Accuracy

Each part is measured in ulps against mpmath at 192 and 320 bits.
The fixtures cover the tiny, subnormal, huge and near-overflow regions,
the branch points, and both sides of every branch cut. Worst part per
function:

| function | special/cmplx | math/cmplx (Go 1.27.3) |
|---|--:|--:|
| Mul, Div | 0.5 ulp | 1.2e16 ulp (builtin `*`, `/`) |
| Exp, Cosh, Sinh, Cos, Sin | 1.29 ulp | 439 ulp, Inf near Re z = 710 |
| Tanh, Tan, Cot | 2.55 ulp | NaN for \|Re z\| > 355 |
| Log, Log10, Log2, Sqrt | 1.03 ulp (0.5 outside subnormals) | 8.9e15 ulp near \|z\| = 1 |
| Log1p, Expm1 | 0.53 ulp | 2.5e27 ulp (naive formula) |
| Atanh, Atan | 1.44 ulp | 1.4e16 ulp |
| Asinh, Asin, Acos, Acosh | 2.82 ulp | 1e9 ulp, wrong quadrant above 1e154 |
| Pow | 0.51 ulp; up to 2937 ulp where \|y log x\| reaches 3e5 | 3.4e20 ulp |

Two areas are open:

- [special#26](https://github.com/go-numerics/special/issues/26):
  Asinh, Asin, Acos and Acosh are above the 2-ulp target.
- [special#27](https://github.com/go-numerics/special/issues/27): the
  error of Pow grows with |y log x|, though it stays within the design's
  (2 + |y log x|) ulp bound.

Every region is listed in
[accuracy-complex.md](https://github.com/go-numerics/special/blob/main/docs/accuracy-complex.md),
and the design is in
[design-complex.md](https://github.com/go-numerics/special/blob/main/docs/design-complex.md).
Timing against NumPy's complex ufuncs has not been measured yet.

## What math/cmplx gets wrong, and why

`math/cmplx` evaluates each textbook formula directly. That fails in the
finite range in three ways:

- **Overflow of intermediates:** `cosh(2x)` in `Tanh`, `e^x` in `Exp`.
- **Cancellation near the unit circle and the branch points:**
  `log(hypot)`, and `log(iz + sqrt(1-z²))` in the inverse functions.
- **The real functions underneath:** `math.Sincos` and `math.Tan` keep
  the reduced argument to an absolute accuracy, so they lose up to
  hundreds of ulps near their zeros and poles. `math.Atan2(-5e-324,
  -1e130)` returns +π.
