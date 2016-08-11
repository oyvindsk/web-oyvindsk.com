package blog

import (
	"net/http"

	"google.golang.org/appengine"
)

const (
	storageBucketPath = "http://storage.googleapis.com/stunning-symbol-139515.appspot.com/oyvindsk.com-static"
)

// Gobal variables
var (
//Posts       map[string]Post
// Pages map[string]Page
//StaticFiles map[string]StaticFile
)

func init() {
	http.HandleFunc("/init", handleFileLoads)
	http.HandleFunc("/", servePageFromDS)
	//log.Fatal("I'm done!")
}

func handleFileLoads(w http.ResponseWriter, r *http.Request) {
	// Read pages and blogpost from file to update the Datastore
	c := appengine.NewContext(r)
	err := loadPagesIntoDS(c, "pages")
	checkAndDie("Reading Pages", err)
	err = loadPagesIntoDS(c, "blogposts")
	checkAndDie("Reading BlogPosts", err)
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
