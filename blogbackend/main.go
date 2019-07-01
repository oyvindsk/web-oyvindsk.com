// Buildtag to only build on AE?
// PLUS build appengine

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"cloud.google.com/go/datastore"
)

const (
	GCSPath = "https://storage.googleapis.com/stunning-symbol-139515.appspot.com/oyvindsk.com-static"
)

var (
	dsClient *datastore.Client
	ctx      context.Context // TODO is there one in the http request we can use instead? but then.. hmm..
)

func main() {

	// Init an empty context
	// TODO is there one in the http request we can use instead? but then.. hmm..
	ctx = context.Background()

	// Create a datastore client. In a typical application, you would create
	// a single client which is reused for every datastore operation.
	var err error
	dsClient, err = datastore.NewClient(ctx, "stunning-symbol-139515")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", servePageFromDS)
	http.HandleFunc("/init", handleFileLoads)
	http.HandleFunc("/writing/", servePostFromDS)

	// static files are served directly from Google Cloud Storage

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" 
	}

	log.Printf("Getting ready to listen on port: %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func handleFileLoads(w http.ResponseWriter, r *http.Request) {
	// Read pages and blogpost from file to update the Datastore
	err := loadPagesIntoDS(ctx, "pages")
	checkAndDie("Reading Pages", err)

	err = loadPostsIntoDS(ctx, "blogposts")
	checkAndDie("Reading BlogPosts", err)
}

func servePageFromDS(w http.ResponseWriter, r *http.Request) {
	key := datastore.NameKey("Page", path.Base(r.URL.Path), nil)
	var page Page
	err := dsClient.Get(ctx, key, &page)
	if err != nil {
		http.NotFound(w, r)
		log.Printf("! 404 for Page?: Path: %s, Title: %s, err: %s", path.Base(r.URL.Path), page.Title, err)
		return
	}
	log.Printf("Serving page: Path: %s, Title: %s", path.Base(r.URL.Path), page.Title)
	w.Write([]byte(page.Content))
}

func servePostFromDS(w http.ResponseWriter, r *http.Request) {
	key := datastore.NameKey("Post", path.Base(r.URL.Path), nil)
	var post Post
	err := dsClient.Get(ctx, key, &post)
	if err != nil {
		http.NotFound(w, r)
		log.Printf("! 404 for Post?: Path: %s, Title: %s, err: %s", path.Base(r.URL.Path), post.Title, err)
		return
	}
	log.Printf("Serving post: Path: %s, Title: %s", path.Base(r.URL.Path), post.Title)
	w.Write([]byte(post.Content))
}

/*
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

	//log.Println("Listening...")
	//log.Fatal(http.ListenAndServe(":3001", nil))

*/
