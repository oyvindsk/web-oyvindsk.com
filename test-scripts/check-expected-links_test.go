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

const baseURL = "http://localhost:8080" // "https://oyvindsk.com"

func TestLinks(t *testing.T) {

	var links = []struct {
		url    string
		needle string
	}{
		{"/", "I’m an independent consultant, developer and architect"},
		{"/hire-me", "Need additional Golang resources"},
		{"/contact", "Contact me"},                    // v2
		{"/about", "After studying computer science"}, // v2
		{"/writing", "FIXME"},                         // v2
		{"/projects", "Intolife"},                     // v2
		{"/now", "Traveling?"},                        // v2

		{"/writing/how-to-use-google-cloud-storage-with-golang", "I have to rant a little about the Google Cloud"},
		{"/writing/common-golang-mistakes-1", "Print out the numbers from 0 to 2"},
		{"/writing/go-remote-jobs", "This is probably the site with the most Go specific listings"},
		{"/writing/docker-build-from-source", "Being lazy, I"},
		{"/writing/reasons-redis-is-awesome", "Even the basic datatypes become super useful"},
	}

	// Custom http cklient to set a shorter timeout
	client := &http.Client{
		Timeout: time.Second * 5,
	}

	for _, l := range links {

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

		fmt.Printf("\t%s done\n", l.url)
	}

}
