package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/a-h/templ"
	"github.com/gomarkdown/markdown"
	mdHTML "github.com/gomarkdown/markdown/html"
	g "maragu.dev/gomponents"
	gc "maragu.dev/gomponents/components"
	gh "maragu.dev/gomponents/html"
)

func renderHTML5Page(innerHTML []byte) string {
	var b strings.Builder

	_ = gc.HTML5(gc.HTML5Props{
		Title:    "My HTML5 Page",
		Language: "en",
		Head: []g.Node{
			// gh.Script(gh.Src("https://cdn.tailwindcss.com?plugins=typography")),
			// 			<link
			//   rel="stylesheet"
			//   href="https://cdn.jsdelivr.net/npm/@picocss/pico@2/css/pico.min.css"
			// >
			gh.Link(gh.Rel("stylesheet"), gh.Href("https://cdn.jsdelivr.net/npm/@picocss/pico@2/css/pico.classless.min.css")),

			g.Raw(`

			<style>
			--pico-background-color: rgb(19, 22.5, 30.5);
			--pico-text-selection-color: rgba(245, 107, 61, 0.1875);
			--pico-primary: #f56b3d;
			--pico-primary-background: #d24317;
			--pico-primary-underline: rgba(245, 107, 61, 0.5);
			--pico-primary-hover: #f8a283;
			--pico-primary-hover-background: #e74b1a;
			--pico-primary-focus: rgba(245, 107, 61, 0.375);
			--pico-primary-inverse: #fff;
			
			</style>
			`),
		},
		Body: []g.Node{
			//    <main class="container">
			gh.Main(
				gh.Class("container"),

				g.Raw(string(innerHTML)),
			),
		},
	}).Render(&b)

	return b.String()
}

func main() {
	// Read the contents of input.md
	// md, err := os.ReadFile("input.md")
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Error reading input.md: %v\n", err)
	// 	os.Exit(1)
	// }

	// html := convertMarkdownToHTML(md)
	// fmt.Printf("\n%s\n", html)

	component := Button("John")

	http.Handle("/", templ.Handler(component))

	fmt.Println("Listening on :8080")
	http.ListenAndServe(":8080", nil)

	// hello("world").Render(context.Background(), os.Stdout)

}

func convertMarkdownToHTML(md []byte) string {

	// parsed := markdown.Parse(md)

	// https://blog.kowalczyk.info/article/cxn3/advanced-markdown-processing-in-go.html

	// Serve the page at http://localhost:8080
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Read the contents of input.md
		md, err := os.ReadFile("input.md")
		if err != nil {
			http.Error(w, "Error reading input.md", http.StatusInternalServerError)
			return
		}

		// Render the article HTML5 page
		htmlFlags := mdHTML.CommonFlags | mdHTML.HrefTargetBlank
		opts := mdHTML.RendererOptions{Flags: htmlFlags}
		renderer := mdHTML.NewRenderer(opts)
		articleHTML := markdown.ToHTML(md, nil, renderer)

		pageHTML := renderHTML5Page(articleHTML)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(pageHTML))
	})

	fmt.Println("Listening on http://localhost:8080 ...")
	http.ListenAndServe(":8080", nil)
	html := markdown.ToHTML(md, nil, renderer)
	return string(html)
}
