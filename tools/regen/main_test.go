package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

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

func run(t *testing.T, args ...string) (string, error) {
	t.Helper()
	var out, errb bytes.Buffer
	cmd := RootCmd()
	cmd.SetArgs(args)
	cmd.SetOut(&out)
	cmd.SetErr(&errb)
	err := cmd.Execute()
	return out.String(), err
}

func site(t *testing.T) string {
	t.Helper()
	docs := t.TempDir()
	dir := filepath.Join(docs, "content/en/special")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	page := "---\ntitle: p\n---\n\n<!-- From go-numerics/special v0.5.8, docs/perf.md: x -->\nold\n"
	if err := os.WriteFile(filepath.Join(dir, "performance.md"), []byte(page), 0o644); err != nil {
		t.Fatal(err)
	}
	idx := "---\ntitle: i\n---\n\n<!-- From go-numerics/special v0.5.8, README.md: x -->\nold\n"
	if err := os.WriteFile(filepath.Join(dir, "_index.md"), []byte(idx), 0o644); err != nil {
		t.Fatal(err)
	}
	return docs
}

func TestCLIPages(t *testing.T) {
	docs := site(t)
	src := t.TempDir()
	os.MkdirAll(filepath.Join(src, "docs"), 0o755)
	os.WriteFile(filepath.Join(src, "docs/perf.md"), []byte("# Perf\n\nSee [l](../LICENSE).\n"), 0o644)
	out, err := run(t, docs, "special", "v0.8.0", src)
	if err != nil || !strings.Contains(out, "wrote ") {
		t.Fatalf("out %q err %v", out, err)
	}
	b, _ := os.ReadFile(filepath.Join(docs, "content/en/special/performance.md"))
	if !strings.Contains(string(b), "docs/perf.md: x -->") || !strings.Contains(string(b), "special v0.8.0,") ||
		!strings.Contains(string(b), "blob/v0.8.0/LICENSE") {
		t.Errorf("page %q", b)
	}
}

func TestCLIReadme(t *testing.T) {
	docs := site(t)
	rd := filepath.Join(t.TempDir(), "README.md")
	os.WriteFile(rd, []byte("# s\n\nText [a](docs/perf.md).\n"), 0o644)
	out, err := run(t, "readme", docs, "special", "v0.8.0", rd)
	if err != nil || !strings.Contains(out, `relref "special/performance"`) {
		t.Fatalf("out %q err %v", out, err)
	}
}

func TestCLIUsage(t *testing.T) {
	if _, err := run(t, "only", "three", "args"); err == nil {
		t.Error("want an error for 3 args")
	}
	if _, err := run(t, "readme", "x", "y", "z", filepath.Join(t.TempDir(), "missing")); err == nil {
		t.Error("want an error for a missing README")
	}
}

func TestMainExit(t *testing.T) {
	old, oldArgs := osExit, os.Args
	defer func() { osExit, os.Args = old, oldArgs }()
	code := -1
	osExit = func(c int) { code = c }
	os.Args = []string{"docsregen", "a"}
	main()
	if code != 2 {
		t.Errorf("exit %d, want 2", code)
	}
}
