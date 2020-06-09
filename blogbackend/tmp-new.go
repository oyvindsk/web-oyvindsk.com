package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/oyvindsk/web-oyvindsk.com/blogbackend/internal/tachyons"
	"github.com/oyvindsk/web-oyvindsk.com/internal/metadata"
)

type serverNew struct {
	cfg       serverNewCfg
	pages     map[string]serverNewContent
	blogposts map[string]serverNewContent // servepath (from the metadata file) => metadata + dirpath

	// pages     map[string]pageContent
	// blogposts map[string]blogpostContent
	// email     struct {
	// 	toAddress, smtpServer, smtpUser, smtpPass string
	// 	smtpPort                                  int
	// }
}

type serverNewCfg struct {
	pathBlogposts string
	pathPages     string
	pathTemplates string
}

type serverNewContent struct {
	metadata.M
	dirpath string
}

func (s *serverNew) loadMetadata() error {

	// Load blogposts
	bppaths, err := metadata.Find(s.cfg.pathBlogposts, "metadata.toml")
	if err != nil {
		return err
	}

	s.blogposts = make(map[string]serverNewContent)

	for _, p := range bppaths {
		fmt.Println("b path:", p)
		m, err := metadata.Fromfile(p + "/metadata.toml")
		if err != nil {
			return fmt.Errorf("loadMetadata: %w", err)
		}
		fmt.Printf("\t%s [%s]\n", m.Title, m.Servepath)
		if _, f := s.blogposts[m.Servepath]; f {
			return fmt.Errorf("loadMetadata: can't load file: %q, Servepath %q aleready loaded", p, m.Servepath)
		}
		s.blogposts[m.Servepath] = serverNewContent{m, p}
	}

	// Load pages
	ppaths, err := metadata.Find(s.cfg.pathPages, "metadata.toml")
	if err != nil {
		return err
	}

	s.pages = make(map[string]serverNewContent)

	for _, p := range ppaths {
		fmt.Println("p path:", p)
		m, err := metadata.Fromfile(p + "/metadata.toml")
		if err != nil {
			return fmt.Errorf("loadMetadata: %w", err)
		}
		fmt.Printf("\t%s [%s]\n", m.Title, m.Servepath)
		if _, f := s.blogposts[m.Servepath]; f {
			return fmt.Errorf("loadMetadata: can't load file: %q, Servepath %q aleready loaded", p, m.Servepath)
		}
		s.pages[m.Servepath] = serverNewContent{m, p}
	}

	// fmt.Printf("%#v\n\n%#v\n", ppaths, ppaths)

	return nil

}

func (s serverNew) serveBlogpost(w http.ResponseWriter, r *http.Request) error {

	log.Printf("newServer: blogpost: %q", r.URL.Path)
	p := path.Base(r.URL.Path)
	log.Printf("newServer: blogpost: looking for path: %q", p)

	content, ok := s.blogposts[p]
	if !ok {
		// it could be an old post, return an error and let caller deal
		return fmt.Errorf("serveBlogpost: Not found")
	}
	log.Printf("newServer: blogpost: %s\n", content.Title)

	body, err := tachyons.PostprocessFile(content.dirpath + "/content.html")
	if err != nil {
		return fmt.Errorf("serveBlogpost: %w", err)
	}

	t, err := template.ParseGlob(s.cfg.pathTemplates + "/*.html")
	if err != nil {
		return fmt.Errorf("serveBlogpost: %w", err)
	}

	tInput := struct {
		serverNewContent
		Activepage string // used to highlite the link in the header,
		PDFurl     string
		Body       template.HTML // Unsafe / unencoded. Input must be safe, a it is here since it comes from ascidoc(tor)
	}{
		content,
		"writing",
		fmt.Sprintf("/writing/%s/full.pdf", content.Servepath),
		template.HTML(body),
	}
	err = t.ExecuteTemplate(w, "blogpost", tInput)
	if err != nil {
		return fmt.Errorf("serveBlogpost: %w", err)
	}

	return nil
}

func (s serverNew) serveBlogpostPDF(w http.ResponseWriter, r *http.Request) {

	// turn the req path into a relative disk path, e.g.:
	// 	/writing/how-to-use-google-cloud-storage-with-golang/full.pdf
	// 		==>
	// 	new-content/blogposts/how-to-use-google-cloud-storage-with-golang/full.pdf
	//
	// split so we can ignore the /writing/ part of the url
	ps := strings.Split(r.URL.Path, "/") // strigs split since this is not a filepath yet
	if len(ps) != 4 {
		log.Printf("newServer: serveBlogpostPDF: failed: path looks wrong: %q => %#v", r.URL.Path, ps)
		s.serve500(w, r)
		return
	}

	// full path on disk
	fp := filepath.Join(s.cfg.pathBlogposts, ps[2], ps[3])
	log.Printf("newServer: serveBlogpostPDF: looking for filepath: %q", fp)
	http.ServeFile(w, r, fp)
}

func (s serverNew) servePage(w http.ResponseWriter, r *http.Request) error {

	log.Printf("newServer: servePage: %q", r.URL.Path)
	p := path.Base(r.URL.Path)
	log.Printf("newServer: servePage: looking for path: %q", p)

	content, ok := s.pages[p]
	if !ok {
		// it could be an old page, return an error and let caller deal
		return fmt.Errorf("newServer: servePage: Not found")
	}
	log.Printf("newServer: servePage: page: %s\n", content.Title)

	body, err := tachyons.PostprocessFile(content.dirpath + "/content.html")
	if err != nil {
		return fmt.Errorf("newServer: servePage: %w", err)
	}

	t, err := template.ParseGlob(s.cfg.pathTemplates + "/*.html")
	if err != nil {
		return fmt.Errorf("newServer: servePage: %w", err)
	}

	tInput := struct {
		serverNewContent
		Activepage string
		PDFurl     string
		Body       template.HTML // Unsafe / unencoded. Input must be safe, a it is here since it comes from ascidoc(tor)
	}{
		serverNewContent: content,
		Activepage:       p,
		Body:             template.HTML(body),
	}

	// Figure out the PDF link
	// with a specialcase for the homepage, it's too different from the others
	if p == "/" {
		tInput.PDFurl = "/home/full.pdf"
	} else {
		tInput.PDFurl = fmt.Sprintf("/%s/full.pdf", p)
	}

	// Execute the blog template with all the data in tInput
	err = t.ExecuteTemplate(w, "page", tInput)
	if err != nil {
		return fmt.Errorf("newServer: servePage: %w", err)
	}

	return nil
}

func (s serverNew) servePagePDF(w http.ResponseWriter, r *http.Request) {
	filepath := path.Join(s.cfg.pathPages, r.URL.Path)
	log.Printf("newServer: servePDF: looking for filepath: %q", filepath)
	http.ServeFile(w, r, filepath)
}

func (s serverNew) serve404(w http.ResponseWriter, r *http.Request) {
	log.Printf("newServer: serving 404 for: %q", r.URL.Path)
	http.Error(w, "404 - Could Not find that =(", http.StatusNotFound)
}

func (s serverNew) serve500(w http.ResponseWriter, r *http.Request) {
	log.Printf("newServer: serving 500 for: %q", r.URL.Path)
	http.Error(w, "500 - I Failed :'(  Please try again a little later", http.StatusInternalServerError)
}
