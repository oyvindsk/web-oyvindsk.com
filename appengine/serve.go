package blog

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
)

//    if err != nil {
//        http.Error(w, err.Error(), http.StatusInternalServerError)
//    }

/*
func servePagesFromMem(w http.ResponseWriter, r *http.Request) {
	page, exists := Pages[path.Base(r.URL.Path)]
	if !exists {
		http.NotFound(w, r)
		log.Printf("! 404: Path: %s, Title: %s", path.Base(r.URL.Path), page.Title)
		return
	}
	log.Printf("Serving page: Path: %s, Title: %s", path.Base(r.URL.Path), page.Title)
	w.Write([]byte(page.Content))
}

func serveBlogPostFromMem(w http.ResponseWriter, r *http.Request) {
	post, exists := Posts[path.Base(r.URL.Path)]
	if !exists {
		http.NotFound(w, r)
		log.Printf("! 404: Path: %s, Title: %s", path.Base(r.URL.Path), post.Title)
		return
	}
	log.Printf("Serving blogpost: Path: %s, Title: %s", path.Base(r.URL.Path), post.Title)
	w.Write([]byte(post.Content))
}

func serveStaticFilestFromMem(w http.ResponseWriter, r *http.Request) {
	staticFile, exists := StaticFiles[path.Base(r.URL.Path)]
	if !exists {
		http.NotFound(w, r)
		log.Printf("! 404: Path: %s", path.Base(r.URL.Path))
		return
	}

	reader := bytes.NewReader(staticFile.ContentRaw)

	log.Printf("Serving static file: Path: %s, StaticFile.Path %s", path.Base(r.URL.Path), staticFile.Path)

	http.ServeContent(w, r, path.Base(r.URL.Path), staticFile.PubTime, reader) // FIXME
}
*/
func servePages(w http.ResponseWriter, r *http.Request) {

	lp := path.Join("templates", "layout.html")
	fp := path.Join("templates", r.URL.Path)

	log.Println("Serving: '" + r.URL.Path + "' '" + fp + "'")

	// Return 404 if the template doesn't exist
	info, err := os.Stat(fp)

	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
	}

	// Return 404 if the req is for a dir
	if r.URL.Path == "/" {
		log.Println("index")
		fp = path.Join("templates", "index.html")
	} else if info.IsDir() {
		http.NotFound(w, r)
		return
	}

	// Parse the template files
	tmpl, err := template.ParseFiles(lp, fp)

	if checkErrHttp(err, w) {
		return
	}

	// Execute the templates
	err = tmpl.ExecuteTemplate(w, "layout", nil)

	if checkErrHttp(err, w) {
		return
	}
}

func serveBlogPost(w http.ResponseWriter, r *http.Request) {

	lp := path.Join("templates", "layout.html")
	fp := path.Join("templates", "blogpost.html")
	pp := path.Join("posts", path.Base(r.URL.Path)+".html") // This is the same as r.URL.Path + .html .. using path.Base is a security measure.. is it neccessary??

	log.Println("Serving blogpost: '" + pp + "'")

	// Return 404 if the blogpost doesn't exist
	_, err := os.Stat(pp)

	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
	}

	// Parse the template files
	// when parsing thm all at once, use template.ParseGlob() and/or tmpl.ParseGlob()
	// or maybe for .. { tmpl.ParseFile() }
	tmpl, err := template.ParseFiles(lp, fp, pp)

	if checkErrHttp(err, w) {
		return
	}

	// Execute the templates
	err = tmpl.ExecuteTemplate(w, "layout", "A few (more than one and less than ten) reasons Redis is awesome!") // FIXME !!

	if checkErrHttp(err, w) {
		return
	}
}
