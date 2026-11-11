package main

import "testing"

func TestReadme(t *testing.T) {
	page := "---\ntitle: \"x\"\n---\n\n<!-- From go-numerics/special v0.5.8, README.md: the text. -->\nold\n"
	readme := "<p>logo</p>\n\n# special\n\n[![b](u)](v)\n\nText [a](docs/perf.md) and [l](LICENSE).\n"
	got := Readme(page, "special", "v0.8.0", readme, map[string]string{"docs/perf.md": "performance"})
	want := "---\ntitle: \"x\"\n---\n\n<!-- From go-numerics/special v0.8.0, README.md: the text. -->\n" +
		"Text [a]({{< relref \"special/performance\" >}}) and [l](https://github.com/go-numerics/special/blob/v0.8.0/LICENSE).\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestRewrite(t *testing.T) {
	pages := map[string]string{"docs/perf.md": "performance", "docs/accuracy-table.md": "accuracy-table"}
	cases := map[string]string{
		"[a](perf.md)":                  `[a]({{< relref "special/performance" >}})`,
		"[a](accuracy-table.md#x)":      `[a]({{< relref "special/accuracy-table#x" >}})`,
		"[a](gamma-fastpath.md)":        "[a](https://github.com/go-numerics/special/blob/v1/docs/gamma-fastpath.md)",
		"[a](../LICENSE)":               "[a](https://github.com/go-numerics/special/blob/v1/LICENSE)",
		"[a](https://example.org/x.md)": "[a](https://example.org/x.md)",
		"[a](#frag)":                    "[a](#frag)",
	}
	for in, want := range cases {
		if got := Rewrite(in, "special", "v1", pages); got != want {
			t.Errorf("%q: got %q, want %q", in, got, want)
		}
	}
}
