package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"

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
	path := path.Base(r.URL.Path)
	log.Printf("newServer: blogpost: looking for path: %q", path)

	content, ok := s.blogposts[path]
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
		Activepage string
		Body       template.HTML // Unsafe / unencoded. Input must be safe, a it is here since it comes from ascidoc(tor)
	}{
		content,
		"blog",
		template.HTML(body),
	}
	err = t.ExecuteTemplate(w, "blogpost", tInput)
	if err != nil {
		return fmt.Errorf("fooPath: %s", err)
	}

	return nil
}

func (s serverNew) servePage(w http.ResponseWriter, r *http.Request) error {

	log.Printf("newServer: servePage: %q", r.URL.Path)
	path := path.Base(r.URL.Path)
	log.Printf("newServer: servePage: looking for path: %q", path)

	content, ok := s.pages[path]
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
		Body       template.HTML // Unsafe / unencoded. Input must be safe, a it is here since it comes from ascidoc(tor)
	}{
		content,
		path,
		template.HTML(body),
	}
	err = t.ExecuteTemplate(w, "page", tInput)
	if err != nil {
		return fmt.Errorf("newServer: servePage: %w", err)
	}

	return nil
}

func (s serverNew) serve404(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "404 - Could Not find that =(", http.StatusNotFound)
}
