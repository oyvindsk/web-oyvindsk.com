package blog

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// The Page structure defines all the data for a page, including metadata / front-matter and resulting HTML
// also includes some template data
type Page struct {
	Title             string
	PubTime           time.Time
	Content           []byte // the result of parsing and executing the templates, for easy serving
	Path              string // the path - last part of the url - that is the address of this page. Might want to make this an [] ?
	StorageBucketPath string
}

// PageDS holds the data we store in DS for this page.
// Since already have "compiled" the template into html, we only realy need the path (key) and the content (html)
type PageDS struct {
	Title   string
	Content []byte // the result of parsing and executing the templates, for easy serving
}

// Read all pages and its front-matter into Datastore for later serving
func loadPagesIntoDS(c context.Context, dir string) error {

	// Read the page from file
	log.Println("Reading pages from:", dir)
	files, err := ioutil.ReadDir(dir)
	checkAndDie("Listing pages", err)

	files = fileFilterHTML(files)

	// Loop the page files
	for _, f := range files {

		log.Println("\t", f.Name())

		// Create a new page from a file, loads front matter and content
		page, err := newPage(dir + "/" + f.Name())
		if err != nil {
			checkAndWarn("Creating Page for file: "+f.Name(), err)
			continue
		}

		// Execute all the necessary templates to get the final html for this post
		err = page.parseAndExecuteTemplates()
		if err != nil {
			checkAndWarn("Execute template for file: "+f.Name(), err)
			continue
		}

		// yeay, store the html w/ the path we want to reach it by
		data := &PageDS{Title: page.Title, Content: page.Content}
		key := datastore.NewKey(c, "Page", page.Path, 0, nil)
		if _, err := datastore.Put(c, key, data); err != nil {
			log.Printf("Failed to store Page: %s, : %s", f.Name(), err)
		}
	}

	return nil
}

func newPage(filepath string) (*Page, error) {

	// read the file, frontmatter + contentpage
	fm, content, err := readFile(filepath)
	if err != nil {
		return nil, err
	}
	//log.Printf("read file:\n%+v\n\n%s\n", fm, content)

	page := Page{}
	var ok bool
	page.Title, ok = fm["Title"]
	if !ok {
		return nil, fmt.Errorf("No Title for page: %s", filepath)
	}
	page.Path, ok = fm["Path"]
	if !ok {
		return nil, fmt.Errorf("No Path for page: %s", filepath)
	}
	// todo - copy inn data
	page.Content = content

	log.Printf("\tloaded page with content size: %d", len(page.Content))

	return &page, nil
}

// Parse and execute the html/template templates from file and store the resulting html in a Post
// Currently 3 files needed to make up a post:
//		layout.html 	- Main html sceleton, shared with other pages?
//		blogpost.html 	- Extra html sceleton just for blogposts
// 		post 			- This is the actual blogpost. Already read into memory.
//  					  Note this is a template too, not just data passed to Execut..()

// Don't store the html in DS but in Memcache instead??

func (p *Page) parseAndExecuteTemplates() error {

	// parse
	// execute!

	buf := bytes.NewBuffer(nil)

	// page as a template, already read in:
	tmpl := template.Must(template.New("").Parse(string(p.Content)))

	// layouts
	_, err := tmpl.ParseFiles("templates/layout.html")

	// execute them all, start with "layout" (defined in the tmpl)
	err = tmpl.ExecuteTemplate(buf, "layout", p)
	if err != nil {
		return fmt.Errorf("template execution: %s", err)
	}

	p.Content = buf.Bytes()
	return nil

	// debug stuff:
	// log.Println(buf)
	// ts := tmpl.Templates()
	// os.Stdout.WriteString("\n\n")
	// for _, t := range ts {
	// 	log.Println(t.Name())
	// }

}
