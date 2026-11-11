// docsregen rewrites the library-derived pages of the go-numerics docs site
// from a checkout of the library: front matter kept, "From ... vX" comment
// moved to the new tag, first heading dropped, relative links pointed into
// the site (relref) or at the tagged file on GitHub.
//
//	docsregen <docs repo> <library> <tag> <library checkout>
package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

var commentRE = regexp.MustCompile(`<!-- From go-numerics/(\w+) (v[0-9.]+), (docs/[^:\s]+):[^>]*-->`)
var linkRE = regexp.MustCompile(`\]\(([^)\s]+)\)`)

// Rewrite points the relative links of src into the site or at GitHub.
func Rewrite(src, lib, tag string, pages map[string]string) string {
	return linkRE.ReplaceAllStringFunc(src, func(m string) string {
		t := linkRE.FindStringSubmatch(m)[1]
		if strings.Contains(t, "://") || strings.HasPrefix(t, "#") || strings.HasPrefix(t, "mailto:") || strings.HasPrefix(t, "{{") {
			return m
		}
		file, anchor, _ := strings.Cut(t, "#")
		if anchor != "" {
			anchor = "#" + anchor
		}
		if stem, ok := pages[path.Clean(path.Join("docs", file))]; ok {
			return `]({{< relref "` + lib + "/" + stem + anchor + `" >}})`
		}
		return "](https://github.com/go-numerics/" + lib + "/blob/" + tag + "/" + path.Clean(path.Join("docs", file)) + anchor + ")"
	})
}

// Readme renders a README as the section's _index page: the front matter and
// comment of page, the logo, title and badges dropped, links relative to the
// repository root.
func Readme(page, lib, tag, readme string, pages map[string]string) string {
	m := commentRE2.FindStringSubmatchIndex(page)
	head := page[:m[1]]
	head = verRE.ReplaceAllString(head, "${1}"+tag+"${2}")
	lines := strings.SplitAfter(readme, "\n")
	k := 0
	for k < len(lines) {
		l := strings.TrimSpace(lines[k])
		if l == "" || strings.HasPrefix(l, "<p") || strings.HasPrefix(l, "# ") || strings.HasPrefix(l, "[![") {
			k++
			continue
		}
		break
	}
	body := strings.Join(lines[k:], "")
	body = linkRE.ReplaceAllStringFunc(body, func(m string) string {
		t := linkRE.FindStringSubmatch(m)[1]
		if strings.Contains(t, "://") || strings.HasPrefix(t, "#") || strings.HasPrefix(t, "mailto:") || strings.HasPrefix(t, "{{") {
			return m
		}
		file, anchor, _ := strings.Cut(t, "#")
		if anchor != "" {
			anchor = "#" + anchor
		}
		file = path.Clean(file)
		if stem, ok := pages[file]; ok {
			return `]({{< relref "` + lib + "/" + stem + anchor + `" >}})`
		}
		return "](https://github.com/go-numerics/" + lib + "/blob/" + tag + "/" + file + anchor + ")"
	})
	return head + "\n" + body
}

var verRE = regexp.MustCompile(`(From go-numerics/\w+ )v[0-9.]+(, README\.md)`)

var commentRE2 = regexp.MustCompile(`<!-- From go-numerics/\w+ v[0-9.]+, README.md[^>]*-->`)

func main() {
	if len(os.Args) == 6 && os.Args[1] == "readme" {
		// docsregen readme <docs repo> <library> <tag> <README file>
		docs, lib, tag := os.Args[2], os.Args[3], os.Args[4]
		files, _ := filepath.Glob(filepath.Join(docs, "content/en", lib, "*.md"))
		pages := map[string]string{}
		for _, f := range files {
			b, _ := os.ReadFile(f)
			if m := commentRE.FindStringSubmatch(string(b)); m != nil {
				pages[m[3]] = strings.TrimSuffix(filepath.Base(f), ".md")
			}
		}
		page, err := os.ReadFile(filepath.Join(docs, "content/en", lib, "_index.md"))
		rd, err2 := os.ReadFile(os.Args[5])
		if err != nil || err2 != nil {
			fmt.Fprintln(os.Stderr, err, err2)
			os.Exit(1)
		}
		fmt.Print(Readme(string(page), lib, tag, string(rd), pages))
		return
	}
	if len(os.Args) != 5 {
		fmt.Fprintln(os.Stderr, "usage: docsregen <docs repo> <library> <tag> <library checkout>")
		os.Exit(2)
	}
	docs, lib, tag, src := os.Args[1], os.Args[2], os.Args[3], os.Args[4]
	files, _ := filepath.Glob(filepath.Join(docs, "content/en", lib, "*.md"))
	pages := map[string]string{}
	texts := map[string]string{}
	for _, f := range files {
		b, err := os.ReadFile(f)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		texts[f] = string(b)
		if m := commentRE.FindStringSubmatch(string(b)); m != nil {
			pages[m[3]] = strings.TrimSuffix(filepath.Base(f), ".md")
		}
	}
	for _, f := range files {
		m := commentRE.FindStringSubmatch(texts[f])
		if m == nil {
			fmt.Println("skipped (no docs/ source):", f)
			continue
		}
		b, err := os.ReadFile(filepath.Join(src, m[3]))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		body := string(b)
		if i := strings.Index(body, "\n"); strings.HasPrefix(body, "# ") && i > 0 {
			body = strings.TrimPrefix(body[i+1:], "\n")
		}
		text := texts[f]
		end := strings.Index(text, m[0])
		comment := strings.Replace(m[0], " "+m[2]+",", " "+tag+",", 1)
		out := text[:end] + comment + "\n\n" + Rewrite(body, lib, tag, pages)
		if err := os.WriteFile(f, []byte(out), 0o644); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println("wrote", f)
	}
}
