# docs

[![docs](https://github.com/go-numerics/docs/actions/workflows/docs.yml/badge.svg)](https://github.com/go-numerics/docs/actions/workflows/docs.yml)
[![License](https://img.shields.io/badge/license-BSD--3--Clause-7445AD?style=flat-square)](LICENSE)

**The go-numerics documentation site**, published at
<https://go-numerics.github.io/docs/>.

A [Hugo](https://gohugo.io/) site with the [tannevaled/hextra](https://github.com/tannevaled/hextra)
fork of the [Hextra](https://github.com/imfing/hextra) theme, built the same way
as the [go-plot documentation](https://github.com/go-plot/docs), with
go-numerics' branding. It documents
[`go-numerics/optimize`](https://github.com/go-numerics/optimize),
[`go-numerics/integrate`](https://github.com/go-numerics/integrate) and
[`go-numerics/special`](https://github.com/go-numerics/special), one section
each.

## Versions

`params.libraries` in `hugo.yaml` names the release of each library the pages
describe. A version of the documentation is **the newest minor version among
them**, or `main` while none is tagged:

| Source | Published under |
| --- | --- |
| branch `main`, with `params.libraries` in `hugo.yaml` | `https://go-numerics.github.io/docs/<MAJOR.MINOR>/`, or `/docs/main/` |
| the newest version (else `main`) | also `https://go-numerics.github.io/docs/latest/` |

`https://go-numerics.github.io/docs/` redirects to `latest/`. Each push to
`main` replaces the directory of the version it describes, and nothing else.
The pages describe **optimize v0.1.0, integrate v0.1.0 and special v0.1.0**,
published under `0.1/`.

To document a new release, set its entry in `params.libraries`, require the
same version in `figures/go.mod`, redraw the figures, and update the pages.

## Pages

A library's pages are its own documentation: its `README.md` becomes the
section's home and each `docs/*.md` a page, with the text unchanged and the
links pointed into this site. A comment at the top of each page names the
file and the release it came from. Change the text in the library first, then
here.

## Figures

Every figure on the site is computed by a go-numerics library and drawn by
[go-plot](https://github.com/go-plot/plot). `figures/` is a Go module, one file
per figure; `go run .` there writes `static/images/figures/<name>.svg`. The
`plotfig` shortcode (`layouts/_shortcodes/plotfig.html`) shows an image and
then **reads the file that drew it**, so the code on the page is the code that
ran. The `figures` job of the CI draws them all again and fails if an image
differs from the committed one, or if `figures/go.mod` requires another
version of a library than the one the pages describe.

```sh
cd figures && go run .
```

The figures are drawn and compared on **arm64**. Go fuses multiply-adds
there and not on amd64, so the libraries' last digits differ between the two,
and so do numbers printed in three of the titles (Quad's error estimate,
Radau's step count, the far tail of the log plots). Redraw them on arm64 —
any Apple-silicon Mac, or `GOARCH=arm64` under emulation.

## Layout

| Path | What |
| --- | --- |
| `content/en/` | the pages; `_index.md` is the home, one directory per library, `weight` orders the sidebar |
| `i18n/en.yaml` | the site's strings, beside the theme's |
| `hugo.yaml` | `baseURL`, the library versions described, the navbar (version › theme › search › GitHub), the brand mounts |
| `figures/` | the Go module that draws every figure, one file each |
| `layouts/_shortcodes/plotfig.html` | a figure and the code that drew it |
| `static/images/figures/` | the figures, as `figures/` writes them |
| `layouts/_partials/custom/content-begin.html` | the outdated-version banner (reads `versions.json`) |
| `layouts/_partials/favicons.html` | the brand's favicons |
| `assets/css/custom.css` | the brand violet as Hextra's primary colour, the table headers, a white card under figures |
| `themes/hextra` | submodule: [tannevaled/hextra](https://github.com/tannevaled/hextra), pinned to a release tag |
| `branding` | submodule: [go-numerics/brand](https://github.com/go-numerics/brand) (the mark, its PNGs and ICO) |
| `scripts/publish-version.sh` | puts one build into a `gh-pages` checkout and rewrites `versions.json`, `latest` and the root redirect |

## Working locally

```sh
git clone --recurse-submodules https://github.com/go-numerics/docs.git
cd docs
mkdir -p data && go run github.com/tannevaled/hextra/cmd/pagehistory@v0.21.0 > data/pagehistory.json
hugo server --baseURL http://localhost:1313/docs/0.1/
```

then open <http://localhost:1313/docs/0.1/>. Hugo 0.146 or later; no Python,
no Node. Go 1.27.2 or later only to redraw the figures.

## Publication

`.github/workflows/docs.yml`:

| Job | When | What |
| --- | --- | --- |
| `build` | pull requests, `main` | builds with `--baseURL …/docs/<version>/`, checks internal links and anchors (lychee, offline), uploads the site |
| `links` | pull requests, `main` | fetches every external URL the pages name |
| `figures` | pull requests, `main` | draws every figure again and compares it with the committed image |
| `deploy` | `main` | `scripts/publish-version.sh` into `gh-pages` (started empty on the first deploy), then a commit and a push |

GitHub Pages serves the `gh-pages` branch.

## Licence

BSD-3-Clause.
