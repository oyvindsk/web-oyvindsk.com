// A simple program to loop all domains and www. subdomains and fetch the html.
// The result should always include the same http snippet
package domaintest

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestDomains(t *testing.T) {

	defNeedle := "I’ve taken contract work since 2014 and worked"

	var domains = []struct {
		domain  string
		needle  string
		skipWWW bool
	}{
		{domain: "oyvindsk.com", needle: defNeedle},
		{domain: "odots.org", needle: defNeedle},

		{domain: "devbyte.no", needle: defNeedle},
		{domain: "devbyte-consulting.com", needle: defNeedle},

		{domain: "skaarsolutions.no", needle: defNeedle},
		{domain: "skaar-solutions.no", needle: defNeedle},
		{domain: "skaarsolutions.com", needle: defNeedle},
		{domain: "skaar-solutions.com", needle: defNeedle},

		{domain: "intolife.skaarsolutions.com", needle: "Food Waste Projects", skipWWW: true},
	}

	// Custom http cklient to set a shorter timeout
	client := &http.Client{
		Timeout: time.Second * 5,
	}

	for _, d := range domains {

		// Test all combinations: www. and http(s)
		for _, u := range []string{
			"http://",
			"https://",
			"http://www.",
			"https://www.",
		} {

			if d.skipWWW && (u == "http://www." || u == "https://www.") {
				continue
			}

			// GET
			url := u + d.domain
			res, err := client.Get(url)
			if err != nil {
				t.Errorf("%q: %q failed: err: %s", d.domain, url, err)
				continue
			}

			defer res.Body.Close()

			// Checks
			if res.StatusCode != http.StatusOK {
				t.Errorf("%q: %q failed: Could not fetch, got status: %q", d.domain, url, res.Status)
			}

			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Errorf("%q: %q failed: err: %s", d.domain, url, err)
			}

			if !bytes.Contains(body, []byte(d.needle)) {
				t.Errorf("%q: %q failed: Could not find needle in body. Body len: %d", d.domain, url, len(body))
			}

		}
		fmt.Printf("\t%s done\n", d.domain)
	}

}

func checkOneDomain(t *testing.T) {

}
