package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
)

const (
	gcsPath = "https://storage.googleapis.com/stunning-symbol-139515.appspot.com/oyvindsk.com-static"
)

type server struct {
	pages     map[string]pageContent
	blogposts map[string]blogpostContent
}

func main() {

	s := server{
		pages:     make(map[string]pageContent),
		blogposts: make(map[string]blogpostContent),
	}

	// Parse all our templates on disk
	// ATM, we only do this here, so a restart is required to pick up new content
	// Not an issue since we have to redeploy for new content anyway =/
	// Also, all the template stuff are in the functions, since they're so small and it's not worth the overhead
	pages, err := slurpAndParseAllPages(gcsPath, "pages")
	if err != nil {
		log.Fatalln(err)
	}
	s.pages = pages

	blogposts, err := slurpAndParseAllPosts(gcsPath, "blogposts")
	if err != nil {
		log.Fatalln(err)
	}
	s.blogposts = blogposts

	// HTTP handlers
	// static files are served directly from Google Cloud Storage
	http.HandleFunc("/", s.handlePage)
	http.HandleFunc("/writing/", s.handleBlogpost)

	// HTTP Listen
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Getting ready to listen on port: %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func (s server) handlePage(w http.ResponseWriter, r *http.Request) {

	// Look for the page in the templates we have parsed already
	path := path.Base(r.URL.Path)

	log.Printf("handlePage: looking for path: %q", path)

	pc, ok := s.pages[path]
	if !ok {
		log.Fatalln("handlePage: Could not find: ", path)
	}

	log.Printf("handlePage: found: %#v", pc)

	// execute them all, start with "layout" (defined in the tmpl)
	err := pc.template.ExecuteTemplate(w, "layout", pc)
	if err != nil {
		log.Fatalf("handlePage: template execution: %s", err)
	}
}

func (s server) handleBlogpost(w http.ResponseWriter, r *http.Request) {

	// Look for the blogpost in the templates we have parsed already
	path := path.Base(r.URL.Path)

	log.Printf("handleBlogpost: looking for path: %q", path)

	bpc, ok := s.blogposts[path]
	if !ok {
		log.Fatalln("handleBlogpost: Could not find: ", path)
	}

	log.Printf("handleBlogpost: found: %#v", bpc)

	// execute them all, start with "layout" (defined in the tmpl)
	err := bpc.template.ExecuteTemplate(w, "layout", bpc)
	if err != nil {
		log.Fatalf("template execution: %s", err)
	}
}
