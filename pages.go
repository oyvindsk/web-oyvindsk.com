package main

import (
	"bufio"
	"bytes"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

// The Page structure defines all the data for a page, including metadata / front-matter and resulting HTML
type Page struct {
	Title   string
	PubTime time.Time
	Content template.HTML // the result of parsing and executing the templates, for easy serving
	Path    string        // the path - last part of the url - that is the address of this page. Might want to make this an [] ?
}

// Read all pages and its front-matter into memory - Few pages, shouldn't be a problem
func readPages(dir string) (pages map[string]Page, err error) {

	pages = make(map[string]Page)

	// Read the page from file
	log.Println("Reading pages from:", dir)
	files, err := ioutil.ReadDir(dir)
	checkAndDie("Listing pages", err)

	// Filter out the files we actually want to read  - From http://0value.com/build-a-blog-engine-in-Go
	for i := 0; i < len(files); {
		if files[i].IsDir() || filepath.Ext(files[i].Name()) != ".html" {
			files[i], files = files[len(files)-1], files[:len(files)-1] // replace the one we don't want with the last one, and shorten our slice by 1
		} else {
			i++
		}
	}

	// Loop the page files
	for _, f := range files {

		log.Println("\t", f.Name())

		// open the file
		file, err := os.Open(dir + "/" + f.Name()) // fixme?
		defer file.Close()
		if err != nil {
			checkAndWarn("Open file", err)
			continue
		}

		// Use a scanner to read the page file line-by-line
		scanner := bufio.NewScanner(file)

		// parse the front matter
		frontMatter, err := parseFrontmatter(scanner)
		if err != nil {
			checkAndWarn("Frontmatter for file: "+f.Name(), err)
			continue
		}

		// Create a page
		page := Page{
			Title:   frontMatter["Title"],
			Path:    frontMatter["Path"],
			PubTime: time.Now(), //FIXME
		}

		// Read rest of file
		cont, err := parsePost(scanner) // is cont copied?
		if err != nil {
			checkAndWarn("Post content for file: "+f.Name(), err)
			continue
		}

		// Execute all the necessary templates to get the final html for this post
		// put it back into post - Is this "idiomatic" go? Or should we return the result??
		// the others return.. FIXME
		err = parseAndExecuteTemplatesPages(&page, cont)
		if err != nil {
			checkAndWarn("Execute template for file: "+f.Name(), err)
			continue
		}

		// yeay, store the html w/ the path we want to reach it by
		pages[page.Path] = page
	}

	//log.Printf("%s", pages)
	return
}

// Parse and execute the html/template templates from file and store the resulting html in a Post
// Currently 3 files needed to make up a post:
//		layout.html 	- Main html sceleton, shared with other pages?
//		blogpost.html 	- Extra html sceleton just for blogposts
// 		post 			- This is the actual blogpost. Already read into memory.
//  					  Note this is a template too, not just data passed to Execut..()

func parseAndExecuteTemplatesPages(page *Page, cont *bytes.Buffer) error {

	// parse
	// execute!

	buf := bytes.NewBuffer(nil)

	// page as a template, already read in:
	tmpl := template.Must(template.New("").Parse(cont.String()))

	// layouts
	_, err := tmpl.ParseFiles("templates/layout.html")

	// execute them all, start with "layout" (defined in the tmpl)
	err = tmpl.ExecuteTemplate(buf, "layout", page)
	if err != nil {
		log.Fatalf("template execution: %s", err)
	}

	page.Content = template.HTML(buf.String()) // copied? is this necessary..?

	// debug stuff:
	// log.Println(buf)
	// ts := tmpl.Templates()
	// os.Stdout.WriteString("\n\n")
	// for _, t := range ts {
	// 	log.Println(t.Name())
	// }

	return nil

}
