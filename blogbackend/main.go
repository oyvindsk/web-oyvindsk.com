package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/oyvindsk/post2mail"
)

const (
	gcsPath = "https://storage.googleapis.com/stunning-symbol-139515.appspot.com/oyvindsk.com-static"
)

type server struct {
	pages     map[string]pageContent
	blogposts map[string]blogpostContent
	email     struct {
		toAddress, smtpServer, smtpUser, smtpPass string
		smtpPort                                  int
	}
}

func main() {

	// A server type to hold global vars and make them accessible to the handlers
	// by convencion not passed further than the handlers
	// and only written to here in main() (otherwise a lock or guardian goroutine is needed)
	s := server{
		pages:     make(map[string]pageContent),
		blogposts: make(map[string]blogpostContent),
	}

	// Handle expected enironment variables

	// SMTP and email parameters for the contact-me backend. All required!
	s.email.toAddress = os.Getenv("EMAIL_TO")
	s.email.smtpServer = os.Getenv("SMTP_SERVER")
	s.email.smtpUser = os.Getenv("SMTP_USER")
	s.email.smtpPass = os.Getenv("SMTP_PASS")
	s.email.smtpPort = 587 // default for secure smtp
	if s.email.toAddress == "" || s.email.smtpServer == "" || s.email.smtpUser == "" || s.email.smtpPass == "" {
		log.Fatalln("At least one requried environment variable is missing. Giving up. Expects: EMAIL_TO SMTP_SERVER SMTP_USER SMTP_PASS")
	}

	// http port to listen on, from Cloud Run
	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8080"
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
	http.HandleFunc("/contact-form", s.handleContactform)

	// HTTP Listen
	log.Printf("Getting ready to listen on port: %s", httpPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", httpPort), nil))
}

// handlePage looks for the page in the templates we have parsed already
func (s server) handlePage(w http.ResponseWriter, r *http.Request) {
	path := path.Base(r.URL.Path)

	log.Printf("handlePage: looking for path: %q", path)

	pc, ok := s.pages[path]
	if !ok {
		log.Println("handlePage: Could not find: ", path)
		http.Error(w, "Page not found =(", http.StatusNotFound)
		return
	}

	log.Printf("handlePage: found: %#v", pc)

	// execute them all, start with "layout" (defined in the tmpl)
	err := pc.template.ExecuteTemplate(w, "layout", pc)
	if err != nil {
		log.Fatalf("handlePage: template execution: %s", err)
	}
}

// handleBlogpost looks for the blogpost in the templates we have parsed already
func (s server) handleBlogpost(w http.ResponseWriter, r *http.Request) {

	path := path.Base(r.URL.Path)

	log.Printf("handleBlogpost: looking for path: %q", path)

	bpc, ok := s.blogposts[path]
	if !ok {
		log.Println("handleBlogpost: Could not find: ", path)
		http.Error(w, "Blogpost not found =(", http.StatusNotFound)
		return
	}

	log.Printf("handleBlogpost: found: %#v", bpc)

	// execute them all, start with "layout" (defined in the tmpl)
	err := bpc.template.ExecuteTemplate(w, "layout", bpc)
	if err != nil {
		log.Fatalf("template execution: %s", err)
	}
}

// Status returned to the HTTP client for non-pages
// modifying this will affect the json returned to the client
// if so, remember the hardcoded error json string as well
type statusMSG struct {
	Success bool
	Status  string
}

// handleContactform handles http POSTs coming from the contact-me form and sends the info in an email using post2mail
func (s server) handleContactform(w http.ResponseWriter, req *http.Request) {
	log.Printf("handleContactform: Reffer: %s , from IP: %s", req.Referer(), req.RemoteAddr)

	var email post2mail.EmailData

	email.To = s.email.toAddress // Receiving email address, sat at server start

	email.FromName = req.FormValue("name")
	email.FromEmail = req.FormValue("from")
	email.Subject = "Someone filled out your form: " + req.FormValue("subject")
	email.Text = req.FormValue("text")

	// Do stupid spam filtering :P
	spam, reason := post2mail.IsSpam(email)
	if spam {
		log.Printf("handleContactform: skipping spammy post: reason: %q, Refferer: %q, IP: %q, UA: %q", reason, req.Referer(), req.RemoteAddr, req.UserAgent())
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, `{ "Status" : "Not acceptable" , "Success" : "false" }`)
		return
	}

	// Send the info on email
	err := post2mail.FormatAndSendEmail(
		email,
		post2mail.SMTPInfo{Server: s.email.smtpServer, Username: s.email.smtpUser, Password: s.email.smtpPass, Port: s.email.smtpPort},
	)

	// Return some infor to the client
	var status statusMSG
	if err == nil {
		status.Success = true
		status.Status = "OK"
	} else {
		status.Success = false
		status.Status = fmt.Sprintf("Error: %s", err)
	}

	// Return a json status
	w.Header().Set("Content-Type", "application/json")
	j, err := json.Marshal(status)
	// json encoding errors
	if err != nil {
		log.Printf("handleContactform: %d: failed to handle request: json encode failed: %s. Status: %v", http.StatusInternalServerError, err, status)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{ "Status" : "error: json encode failed: %s" , "Success" : "false" }`, err)
		return
	}

	// other problems
	if !status.Success {
		log.Printf("handleContactform: %d: failed to handle request: Status: %v", http.StatusInternalServerError, status)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", j)
		return
	}

	_, err = w.Write(j)
	if err != nil {
		log.Printf("handleContactform: failed write to client: error: %s. Status: %v", err, status)
	}

	// Success!
	log.Printf("handleContactform: Email sent:\t%+v", status)

}
