package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
)

// Blogpost is it's own template, 1 of 2, and has its own type

type blogpostContent struct {
	Author string
	Path   string
	Title  string // exported for the template lib
	// PubTime           time.Time // TODO ?

	template *template.Template
}

func (pm blogpostContent) String() string {
	return fmt.Sprintf(`
	P: %q
	T: %q
	A: %q
	`, pm.Path, pm.Title, pm.Author)
}

func slurpAndParseAllPosts(dirPath string) (map[string]blogpostContent, error) {

	// Read the blogposts from files
	log.Println("Reading blogposts from:", dirPath)
	filesP, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, newHTTPError(http.StatusInternalServerError, "slurpAndParseAllPosts: ReadDir", "", err)
	}

	blogposts := make(map[string]blogpostContent)

	// Loop the blogpost files
	for _, f := range filesP {

		log.Println(f.Name())

		md, html, err := readFile(filepath.Join(dirPath, f.Name()))
		if err != nil {
			return nil, newHTTPError(http.StatusInternalServerError, "slurpAndParseAllPosts: readFile", "", err)
		}

		var bpc blogpostContent

		for k, v := range md {
			switch k {
			case "Author":
				bpc.Author = v
			case "Date":
				// TODO FIXME PubTime
			case "Path":
				bpc.Path = v
			case "Title":
				bpc.Title = v
			default:
				log.Printf("Unknow metadata seen in blogpost: %q", k)
			}
		}

		// blogpost as a template
		tmpl, err := template.New("blogpost").Parse(html.String())
		if err != nil {
			return nil, newHTTPError(http.StatusInternalServerError, "slurpAndParseAllPosts: template.Parse", "", err)
		}

		// layouts
		bpc.template, err = tmpl.ParseFiles("templates/layout.html", "templates/blogpost.html")
		if err != nil {
			return nil, newHTTPError(http.StatusInternalServerError, "slurpAndParseAllPosts: template.ParseFiles", "", err)
		}

		blogposts[bpc.Path] = bpc
	}

	return blogposts, nil
}
