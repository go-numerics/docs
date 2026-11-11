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

	"github.com/spf13/cobra"
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

// osExit is replaced by tests.
var osExit = os.Exit

// pageStems maps each library docs/ source named by a page's comment to the
// page's stem, and returns the pages' files and texts.
func pageStems(docs, lib string) (files []string, texts, pages map[string]string, err error) {
	files, _ = filepath.Glob(filepath.Join(docs, "content/en", lib, "*.md"))
	pages = map[string]string{}
	texts = map[string]string{}
	for _, f := range files {
		b, err := os.ReadFile(f)
		if err != nil {
			return nil, nil, nil, err
		}
		texts[f] = string(b)
		if m := commentRE.FindStringSubmatch(string(b)); m != nil {
			pages[m[3]] = strings.TrimSuffix(filepath.Base(f), ".md")
		}
	}
	return files, texts, pages, nil
}

// RootCmd builds the command: docsregen <docs repo> <library> <tag> <library
// checkout>, or docsregen readme <docs repo> <library> <tag> <README file>.
func RootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:           "docsregen <docs repo> <library> <tag> <library checkout>",
		Short:         "Regenerate the library-derived pages of the go-numerics docs site",
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			docs, lib, tag, src := args[0], args[1], args[2], args[3]
			out := cmd.OutOrStdout()
			files, texts, pages, err := pageStems(docs, lib)
			if err != nil {
				return err
			}
			for _, f := range files {
				m := commentRE.FindStringSubmatch(texts[f])
				if m == nil {
					fmt.Fprintln(out, "skipped (no docs/ source):", f)
					continue
				}
				b, err := os.ReadFile(filepath.Join(src, m[3]))
				if err != nil {
					return err
				}
				body := string(b)
				if i := strings.Index(body, "\n"); strings.HasPrefix(body, "# ") && i > 0 {
					body = strings.TrimPrefix(body[i+1:], "\n")
				}
				text := texts[f]
				end := strings.Index(text, m[0])
				comment := strings.Replace(m[0], " "+m[2]+",", " "+tag+",", 1)
				res := text[:end] + comment + "\n\n" + Rewrite(body, lib, tag, pages)
				if err := os.WriteFile(f, []byte(res), 0o644); err != nil {
					return err
				}
				fmt.Fprintln(out, "wrote", f)
			}
			return nil
		},
	}
	readme := &cobra.Command{
		Use:          "readme <docs repo> <library> <tag> <README file>",
		Short:        "Print the section's _index page rendered from a README",
		SilenceUsage: true,
		Args:         cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			docs, lib, tag := args[0], args[1], args[2]
			_, _, pages, err := pageStems(docs, lib)
			if err != nil {
				return err
			}
			page, err := os.ReadFile(filepath.Join(docs, "content/en", lib, "_index.md"))
			if err != nil {
				return err
			}
			rd, err := os.ReadFile(args[3])
			if err != nil {
				return err
			}
			fmt.Fprint(cmd.OutOrStdout(), Readme(string(page), lib, tag, string(rd), pages))
			return nil
		},
	}
	root.AddCommand(readme)
	return root
}

func main() {
	cmd := RootCmd()
	cmd.SetErr(os.Stderr)
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		if strings.Contains(err.Error(), "arg") || strings.Contains(err.Error(), "command") {
			fmt.Fprintln(os.Stderr, cmd.UsageString())
			osExit(2)
			return
		}
		osExit(1)
	}
}
