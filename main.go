package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// start
// compile templates
// read posts
// server
//      reload     ------ /\

// Gobal variables
var (
	Posts       map[string]Post
	Pages       map[string]Page
	StaticFiles map[string]StaticFile
)

func main() {

	// Initialize directories
	pwd, err := os.Getwd()
	checkAndDie("Getting working dir", err)
	StaticDir := filepath.Join(pwd, "static")
	PostsDir := filepath.Join(pwd, "blogposts")
	PagesDir := filepath.Join(pwd, "pages")

	// Read and parse all blogposts, compile and execute templtes, store result in mem.
	Posts, err = readPosts(PostsDir)
	checkAndDie("Reading Posts", err)

	// do the same with the pages (non-blogposts)
	Pages, err = readPages(PagesDir)
	checkAndDie("Reading Pages", err)

	// and static files (maybe not a good idea..)
	StaticFiles, err = readStaticFiles(StaticDir)
	checkAndDie("Reading Static files", err)

	// Serve static files
	fs := http.FileServer(http.Dir(StaticDir))
	http.Handle("/static2/", http.StripPrefix("/static2/", fs))
	http.HandleFunc("/static/", serveStaticFilestFromMem)

	// Serve blogposts
	http.HandleFunc("/writing/", serveBlogPostFromMem)

	// Serve other pages, index etc.
	http.HandleFunc("/", servePagesFromMem)

	log.Println("Listening...")
	log.Fatal(http.ListenAndServe(":3001", nil))

}
