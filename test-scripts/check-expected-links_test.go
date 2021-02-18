// A simple program to check that all expected links still work and return what we expect
// Avoid dead links! Redirect if we have to
//
// run with `go test ./check-expected-links_test.go` (not a propper package dir structure here yet)
package linktest

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

const baseURL = "https://oyvindsk.com" // "http://localhost:8080"

type linkInfo struct {
	url                    string
	needle                 string
	pdfSizeMin, pdfSizeMax int64 // in kilobytes, match res.contentLength
	// the pdf urls are automatically ..s

}

func TestLinks(t *testing.T) {

	// Pages and blog posts
	// We also look for the pdf that all pages and blogposts should have
	// don't scan their content, but look at the size returened for some sort of comfort

	var links = []linkInfo{
		{"/", "I’m an independent consultant, developer and architect", 50, 60},
		{"/hire-me", "Are you working on building something new?", 70, 80},
		{"/contact", "Contact me", 28, 35},                                   // added in 2020
		{"/writing", "A few places you can find Golang remote jobs", 20, 30}, // added in 2020
		{"/projects", "Intolife", 110, 130},                                  // added in 2020
		//	{"/about", "After studying computer science"}, // add in 2020?
		//	{"/now", "Traveling?"},                        // add in 2020?

		// Blog posts:
		{"/writing/how-to-use-google-cloud-storage-with-golang", "I have to rant a little about the Google Cloud", 90, 95},
		{"/writing/common-golang-mistakes-1", "Print out the numbers from 0 to 2", 100, 150},
		{"/writing/go-remote-jobs", "This is probably the site with the most Go specific listings", 49, 52},
		{"/writing/docker-build-from-source", "Being lazy, I", 80, 90},
		{"/writing/reasons-redis-is-awesome", "Even the basic datatypes become super useful", 43, 50},
	}

	// Custom http cklient to set a shorter timeout
	client := &http.Client{
		Timeout: time.Second * 5,
	}

	// Check normal html pages and blog posts
	for _, l := range links {

		//
		// Normal html page

		// GET
		url := baseURL + l.url
		res, err := client.Get(url)
		if err != nil {
			t.Errorf("%q: failed: err: %s", l.url, err)
			continue
		}

		defer res.Body.Close()

		// Checks
		if res.StatusCode != http.StatusOK {
			t.Errorf("%q: failed: Could not fetch, got status: %q", l.url, res.Status)
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("%q: failed: err: %s", l.url, err)
		}

		if !bytes.Contains(body, []byte(l.needle)) {
			t.Errorf("%q: failed: Could not find needle in body. Body len: %d", l.url, len(body))
		}

		//
		// PDF urls

		// GET
		pdfURL := baseURL + l.url + "/full.pdf"
		if l.url == "/" {
			// special case for home  FIXME TODO ?
			pdfURL = baseURL + "/full.pdf"
		}
		res, err = client.Get(pdfURL)
		if err != nil {
			t.Errorf("%q: (%q) PDF failed: err: %s", l.url, pdfURL, err)
			continue
		}

		defer res.Body.Close()

		// Checks
		if res.StatusCode != http.StatusOK {
			t.Errorf("%q: (%q) PDF failed: Could not fetch, got status: %q", l.url, pdfURL, res.Status)
		}

		if res.ContentLength < l.pdfSizeMin*1020 || res.ContentLength > l.pdfSizeMax*1020 {
			t.Errorf("%q: (%q) PDF failed: PDF size is %d bytes, expected to be %d><%d", l.url, pdfURL, res.ContentLength, l.pdfSizeMin*1020, l.pdfSizeMax*1020)
		}

		fmt.Printf("\t%s done\n", l.url)
	}

}
