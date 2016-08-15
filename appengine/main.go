// Buildtag to only build on AE?
// PLUS build appengine

package blog

import (
	"log"
	"net/http"
	"path"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

const (
	GCSPath = "http://storage.googleapis.com/stunning-symbol-139515.appspot.com/oyvindsk.com-static"
)

func init() {
	http.HandleFunc("/", servePageFromDS)
	http.HandleFunc("/init", handleFileLoads)
	http.HandleFunc("/writing/", servePostFromDS)
	// static files are served directly from Google Cloud Storage
}

func handleFileLoads(w http.ResponseWriter, r *http.Request) {
	// Read pages and blogpost from file to update the Datastore
	c := appengine.NewContext(r)
	err := loadPagesIntoDS(c, "pages")
	checkAndDie("Reading Pages", err)

	err = loadPostsIntoDS(c, "blogposts")
	checkAndDie("Reading BlogPosts", err)
}

func servePageFromDS(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	key := datastore.NewKey(c, "Page", path.Base(r.URL.Path), 0, nil)
	var page Page
	err := datastore.Get(c, key, &page)
	if err != nil {
		http.NotFound(w, r)
		log.Printf("! 404 for Page?: Path: %s, Title: %s, err: %s", path.Base(r.URL.Path), page.Title, err)
		return
	}
	log.Printf("Serving page: Path: %s, Title: %s", path.Base(r.URL.Path), page.Title)
	w.Write([]byte(page.Content))
}

func servePostFromDS(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	key := datastore.NewKey(c, "Post", path.Base(r.URL.Path), 0, nil)
	var post Post
	err := datastore.Get(c, key, &post)
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
