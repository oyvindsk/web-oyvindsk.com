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
	Posts map[string]Post
	Pages map[string]Page
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

	// Serve static files
	fs := http.FileServer(http.Dir(StaticDir))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Serve blogposts
	http.HandleFunc("/writing/", serveBlogPostFromMem)

	// Serve other pages, index etc.
	http.HandleFunc("/", servePagesFromMem)

	log.Println("Listening...")
	log.Fatal(http.ListenAndServe(":84", nil))

}

// Blog functions
//func getBlogPost(title string) (*Post, error) {
//    file := path.Join("/writing", title + ".html")
//	body, err := ioutil.ReadFile(filename)
//	if err != nil {
//		return nil, err
//	}
//	return &Post{Title: title, Body: body}, nil
//
//}

// Error handling functions

func checkAndDie(m string, e error) {
	if e != nil {
		log.Fatal("!! ", m, " : ", e)
	}
}

func checkAndWarn(m string, e error) {
	if e != nil {
		log.Print("! ", m, " : ", e)
	}
}

func checkErr(m string, e error) {
	if e != nil {
		log.Fatal("!! ", m, " : ", e)
	}
}

func checkErrHttp(err error, w http.ResponseWriter) bool {
	if err != nil {
		log.Println("!! ", err)
		http.Error(w, http.StatusText(500), 500)
		return true
	}
	return false
}
