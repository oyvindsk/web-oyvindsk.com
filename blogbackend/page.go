package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
)

// Page is it's own template, 2 of 2, and has its own type

type pageContent struct {
	Author      string
	Description string
	Path        string
	Title       string // exported for the template lib
	// PubTime           time.Time // TODO ?

	template *template.Template
}

func (pm pageContent) String() string {
	return fmt.Sprintf(`
	P: %q
	T: %q
	A: %q
	D: %q
	`, pm.Path, pm.Title, pm.Author, pm.Description)
}

func slurpAndParseAllPages(dirPath string) (map[string]pageContent, error) {

	// Read the pages from files
	log.Println("Reading pages from:", dirPath)
	filesP, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, newHTTPError(http.StatusInternalServerError, "slurpAndParseAllPages: ReadDir", "", err)
	}

	pages := make(map[string]pageContent)

	// Loop the page files
	for _, f := range filesP {

		log.Println(f.Name())

		md, html, err := readFile(filepath.Join(dirPath, f.Name()))
		if err != nil {
			return nil, newHTTPError(http.StatusInternalServerError, "slurpAndParseAllPages: readFile", "", err)
		}

		var pc pageContent

		for k, v := range md {
			switch k {
			case "Author":
				pc.Author = v
			case "Description":
				pc.Description = v
			case "Date":
				// TODO FIXME PubTime
			case "Path":
				pc.Path = v
			case "Title":
				pc.Title = v
			default:
				log.Printf("Unknow metadata seen in page: %q", k)
			}
		}

		// page as a template
		tmpl, err := template.New("").Parse(html.String())
		if err != nil {
			return nil, newHTTPError(http.StatusInternalServerError, "slurpAndParseAllPages: template.Parse", "", err)
		}

		// layouts
		pc.template, err = tmpl.ParseFiles("templates/layout.html")
		if err != nil {
			return nil, newHTTPError(http.StatusInternalServerError, "slurpAndParseAllPages: template.ParseFiles", "", err)
		}

		pages[pc.Path] = pc
	}

	return pages, nil
}
