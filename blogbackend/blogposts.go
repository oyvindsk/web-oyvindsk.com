// http://0value.com/build-a-blog-engine-in-Go

package main

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"time"

	"cloud.google.com/go/datastore"
)

// The Post structure defines all the data for a post, including metadata / front-matter and resulting HTML
// also includes some template data
type Post struct {
	Title             string
	PubTime           time.Time
	Content           []byte // the result of parsing and executing the templates, for easy serving.
	Path              string // the path - last part of the url - that is the address of this post. Might want to make this an [] ?
	StorageBucketPath string
}

// PostDS holds the data we store in DS for this post.
// Since already have "compiled" the template into html, we only realy need the path (key) and the content (html)
// Indexing is the default now (when?), so we have to tell ds not to index the fields that are too long (or just uneccessary)
// TODO no need for two types, ds can ignore fields.. also.. why? :D
type PostDS struct {
	Title   string
	Content []byte `datastore:",noindex"` // the result of parsing and executing the templates, for easy serving
}

// Read all blogposts and its front-matter into Datastore for later serving
func loadPostsIntoDS(c context.Context, dir string) error {

	// Read the blogpost from file
	log.Println("Reading blogposts from:", dir)
	files, err := ioutil.ReadDir(dir)
	checkAndDie("Listing blogpost", err)

	files = fileFilterHTML(files)

	// Loop the blogpost files
	for _, f := range files {

		log.Println("\t", f.Name())

		// Create a new blogpost from a file, loads front matter and content
		post, err := newPost(dir + "/" + f.Name())
		if err != nil {
			checkAndWarn("Creating Post for file: "+f.Name(), err)
			continue
		}

		// Execute all the necessary templates to get the final html for this post
		err = post.parseAndExecuteTemplates()
		if err != nil {
			checkAndWarn("Execute template for file: "+f.Name(), err)
			continue
		}

		// yeay, store the html w/ the path we want to reach it by
		data := &PostDS{Title: post.Title, Content: post.Content}
		key := datastore.NameKey("Post", post.Path, nil)
		if _, err := dsClient.Put(c, key, data); err != nil {
			log.Printf("Failed to import Post: %s, : %s", f.Name(), err)
		}
	}

	return nil
}

func newPost(filepath string) (*Post, error) {

	// read the file, frontmatter + contentpage
	fm, content, err := readFile(filepath)
	if err != nil {
		return nil, err
	}
	//log.Printf("read file:\n%+v\n\n%s\n", fm, content)

	post := Post{}
	var ok bool
	post.Title, ok = fm["Title"]
	if !ok {
		return nil, fmt.Errorf("No Title for page: %s", filepath)
	}
	post.Path, ok = fm["Path"]
	if !ok {
		return nil, fmt.Errorf("No Path for page: %s", filepath)
	}
	post.StorageBucketPath = GCSPath // global const, but still needs to be passed to the template (?)

	post.Content = content

	log.Printf("\tloaded blogpost with content size: %d", len(post.Content))

	return &post, nil
}

// Parse and execute the html/template templates from file and store the resulting html in a Post
// Currently 3 files needed to make up a post:
//		layout.html 	- Main html sceleton, shared with other pages?
//		blogpost.html 	- Extra html sceleton just for blogposts
// 		post 			- This is the actual blogpost. Already read into memory.
//  					  Note this is a template too, not just data passed to Execut..()

// Don't store the html in DS but in Memcache instead??

func (p *Post) parseAndExecuteTemplates() error {

	// parse
	// execute!

	buf := bytes.NewBuffer(nil)

	// page as a template, already read in:
	tmpl := template.Must(template.New("blogpost").Parse(string(p.Content)))

	// layouts
	_, err := tmpl.ParseFiles("templates/layout.html", "templates/blogpost.html")

	// execute them all, start with "layout" (defined in the tmpl)
	err = tmpl.ExecuteTemplate(buf, "layout", p)
	if err != nil {
		return fmt.Errorf("template execution: %s", err)
	}

	p.Content = buf.Bytes()
	return nil

	// debug stuff:
	// ts := tmpl.Templates()
	// os.Stdout.WriteString("\n\n")
	// for _, t := range ts {
	// 	log.Println(t.Name())
	// }

}

// Parse and execute the html/template templates from file and store the resulting html in a Post
// Currently 3 files needed to make up a post:
//		layout.html 	- Main html sceleton, shared with other pages?
//		blogpost.html 	- Extra html sceleton just for blogposts
// 		post 			- This is the actual blogpost. Already read into memory.
//  					  Note this is a template too, not just data passed to Execut..()

/*
func parseAndExecuteTemplatesss(post *Post, cont *bytes.Buffer) error {

	// parse
	// execute!

	buf := bytes.NewBuffer(nil)

	// blogpost as a template, already read in:
	tmpl := template.Must(template.New("blogpost").Parse(cont.String()))

	// layouts
	_, err := tmpl.ParseFiles("templates/layout.html", "templates/blogpost.html")

	// execute them all, start with "layout" (defined in the tmpl)
	err = tmpl.ExecuteTemplate(buf, "layout", post)
	if err != nil {
		log.Fatalf("template execution: %s", err)
	}

	post.Content = template.HTML(buf.String()) // copied? is this necessary..?

	// debug stuff:
	// log.Println(buf)
	// ts := tmpl.Templates()
	// os.Stdout.WriteString("\n\n")
	// for _, t := range ts {
	// 	log.Println(t.Name())
	// }

	return nil

}
*/
