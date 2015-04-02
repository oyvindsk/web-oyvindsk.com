// http://0value.com/build-a-blog-engine-in-Go

package main

import (
	"bufio"
	"bytes"
	"errors"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// The Post structure defines all the data for a post, including metadata / front-matter and resulting HTML
type Post struct {
	Title   string
	PubTime time.Time
	Content template.HTML // the result of parsing and executing the templates, for easy serving
	Path    string        // the path - last part of the url - that is the address of this post. Might want to make this an [] ?

	// Unused ATM:
	//Slug        string
	//Description string
	//Lang        string
	//ModTime time.Time
}

// The TemplateData structure contains all the relevant information passed to the
// template to generate the static HTML.
type TemplateData struct {
	Post *Post
	//Recent []*Post
	//Prev   *ShortPost
	//Next   *ShortPost
}

// Read all posts and its front-matter into memory - Few posts, shouldn't be a problem
func readPosts(dir string) (posts map[string]Post, err error) {

	posts = make(map[string]Post)

	// Read the blogposts from file
	log.Println("Reading posts from:", dir)
	files, err := ioutil.ReadDir(dir)
	checkAndDie("Listing posts", err)

	// Filter out the files we actually want to read  - From http://0value.com/build-a-blog-engine-in-Go
	for i := 0; i < len(files); {
		if files[i].IsDir() || filepath.Ext(files[i].Name()) != ".html" {
			files[i], files = files[len(files)-1], files[:len(files)-1] // replace the one we don't want with the last one, and shorten our slice by 1
		} else {
			i++
		}
	}

	// Loop the post files
	for _, f := range files {

		log.Println("\t", f.Name())

		// open the file
		file, err := os.Open(dir + "/" + f.Name()) // fixme?
		defer file.Close()
		if err != nil {
			checkAndWarn("Open file", err)
			continue
		}

		// Use a scanner to read the post file line-by-line
		scanner := bufio.NewScanner(file)

		// parse the front matter
		frontMatter, err := parseFrontmatter(scanner)
		if err != nil {
			checkAndWarn("Frontmatter for file: "+f.Name(), err)
			continue
		}

		// Create a post
		post := Post{
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
		err = parseAndExecuteTemplates(&post, cont)
		if err != nil {
			checkAndWarn("Execute template for file: "+f.Name(), err)
			continue
		}

		// yeay, store the html w/ the path we want to reach it by
		posts[post.Path] = post
	}

	//log.Printf("%s", posts)
	return
}

// Read the front matter from the post. If there is no front matter, this is
// not a valid post.
// from https://github.com/PuerkitoBio/trofaf/blob/master/tpldata.go

func parseFrontmatter(s *bufio.Scanner) (m map[string]string, err error) {
	m = make(map[string]string)
	infm := false
	for s.Scan() {
		l := strings.Trim(s.Text(), " ")
		if l == "---" { // The front matter is delimited by 3 dashes
			if infm {
				// This signals the end of the front matter
				return m, nil
			} else {
				// This is the start of the front matter
				infm = true
			}
		} else if infm {
			sections := strings.SplitN(l, ":", 2)
			if len(sections) != 2 {
				// Invalid front matter line
				return nil, errors.New("Invalid front matter line")
			}
			m[sections[0]] = strings.Trim(sections[1], " ")
		} else if l != "" {
			// No front matter, quit
			return nil, errors.New("No front matter")
		}
	}
	if err := s.Err(); err != nil {
		// The scanner stopped because of an error
		return nil, err
	}
	return nil, errors.New("Empty post file")
}

// Parse the rest of the blogpost, the stuff after front-matter
// basically just ready evrything left into a buffer

func parsePost(s *bufio.Scanner) (buf *bytes.Buffer, err error) {

	buf = bytes.NewBuffer(nil)
	for s.Scan() {
		buf.WriteString(s.Text() + "\n")
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return
}

// Parse and execute the html/template templates from file and store the resulting html in a Post
// Currently 3 files needed to make up a post:
//		layout.html 	- Main html sceleton, shared with other pages?
//		blogpost.html 	- Extra html sceleton just for blogposts
// 		post 			- This is the actual blogpost. Already read into memory.
//  					  Note this is a template too, not just data passed to Execut..()

func parseAndExecuteTemplates(post *Post, cont *bytes.Buffer) error {

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
