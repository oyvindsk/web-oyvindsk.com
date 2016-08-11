package blog

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"google.golang.org/appengine/datastore"

	"golang.org/x/net/context"
)

type fileBasedContent interface {
	loadFrontMatter(s *bufio.Scanner) (err error)
	loadContent(s *bufio.Scanner) error
	parseAndExecuteTemplates() error
}

type common struct {
	fileBasedContent

	Title   string
	PubTime time.Time
	Content []byte // the result of parsing and executing the templates, for easy serving
	Path    string // the path - last part of the url - that is the address of this page. Might want to make this an [] ?
}

// The Page structure defines all the data for a page, including metadata / front-matter and resulting HTML
type Page struct {
	common
}

// Read all pages and its front-matter into Datastore
func loadPagesIntoDS(c context.Context, dir string) error {

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

		// Create a new page from a file, loads front matter and content
		page, err := newPage(dir + "/" + f.Name())

		// Execute all the necessary templates to get the final html for this post
		err = page.parseAndExecuteTemplates()
		if err != nil {
			checkAndWarn("Execute template for file: "+f.Name(), err)
			continue
		}

		// yeay, store the html w/ the path we want to reach it by
		//pages[page.Path] = page

		key := datastore.NewKey(c, "Page", page.Path, 0, nil)
		if _, err := datastore.Put(c, key, page); err != nil {
			log.Printf("Failed to import Page: %s, : %s", f.Name(), err)
		}
	}

	//log.Printf("%s", pages)
	return nil
}

func newPage(filepath string) (*common, error) {

	// open the file
	file, err := os.Open(filepath)
	defer file.Close()
	if err != nil {
		return nil, err
	}

	page := Page{} // Use a scanner to read the page file line-by-line
	scanner := bufio.NewScanner(file)

	page.loadFrontMatter(scanner)
	if err != nil {
		// checkAndWarn("Frontmatter for file: "+filepath, err)
		return nil, err
	}

	page.loadContent(scanner)
	if err != nil {
		//checkAndWarn("Post content for file: "+f.Name(), err)
		return nil, err
	}

	// Read rest of file
	log.Printf("\tloaded page with content size: %d", len(page.Content))
	return &page, nil
}

// Read the front matter from the post. If there is no front matter, this is
// not a valid post.
// from https://github.com/PuerkitoBio/trofaf/blob/master/tpldata.go

func (p *Page) loadFrontMatter(s *bufio.Scanner) (err error) {
	infm := false
	for s.Scan() {
		l := strings.Trim(s.Text(), " ")
		if l == "---" { // The front matter is delimited by 3 dashes
			if infm {
				// This signals the end of the front matter
				// Success!
				return nil
			}
			// This is the start of the front matter
			infm = true

		} else if infm {
			sections := strings.SplitN(l, ":", 2)
			if len(sections) != 2 {
				// Invalid front matter line
				return errors.New("Invalid front matter line")
			}
			switch sections[0] {
			case "Title":
				p.Title = strings.Trim(sections[1], " ")
			case "Path":
				p.Path = strings.Trim(sections[1], " ")
			default:
				log.Println("frontamtter: ignored unknown:", sections[0])
			}

		} else if l != "" {
			// No front matter, quit
			return errors.New("No front matter")
		}
	}
	if err := s.Err(); err != nil {
		// The scanner stopped because of an error
		return err
	}
	return errors.New("Empty post file")
}

// Parse the rest of the blogpost, the stuff after front-matter
// basically just read evrything left into a buffer

func (p *Page) loadContent(s *bufio.Scanner) error {

	buf := bytes.NewBuffer(nil)
	for s.Scan() {
		buf.WriteString(s.Text() + "\n")
	}
	if err := s.Err(); err != nil {
		return err
	}

	p.Content = buf.Bytes()
	return nil
}

// Parse and execute the html/template templates from file and store the resulting html in a Post
// Currently 3 files needed to make up a post:
//		layout.html 	- Main html sceleton, shared with other pages?
//		blogpost.html 	- Extra html sceleton just for blogposts
// 		post 			- This is the actual blogpost. Already read into memory.
//  					  Note this is a template too, not just data passed to Execut..()

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
