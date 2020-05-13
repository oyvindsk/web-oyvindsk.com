package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"golang.org/x/net/html"
)

// Input:
// BlogpostDate
// BlogpostTitle
// BlogpostSubtitle
// BlogpostAuthor
// BlogpostBody
const templateTest = `<!DOCTYPE html>
<html lang="en">
	<title> </title>
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<link rel="stylesheet" href="https://unpkg.com/tachyons/css/tachyons.min.css">
  
	<body>
		<article>
  			<header class="bg-gold sans-serif">
    			<div class="mw9 center pa4 pt5-ns ph7-l">
      				<time class="f6 mb2 dib ttu tracked"><small>{{ .BlogpostDate }}</small></time>
      				<h3 class="f2 f1-m f-headline-l measure-narrow lh-title mv0">
        				<span class="bg-black-90 lh-copy white pa1 tracked-tight">
						{{ .BlogpostTitle }}
        				</span>
      				</h3>
      				<h4 class="f3 fw1 georgia i">{{ .BlogpostSubtitle }}</h4>
      				<h5 class="f6 ttu tracked black-80">By {{ .BlogpostAuthor }}</h5>
    			</div>
  			</header>

			<div class="pa4 ph7-l georgia mw9-l center">
  
				<!-- auto generert herfrra -->

				{{ .BlogpostBody }}

				<!-- auto generert hit -->

			</div>

		</article>

	</body>
</html>
`

type tachyons struct {
	skip bool
}

func (mt *tachyons) nextToken(t html.Token) {

	// Should we set skip?
	if t.Type == html.StartTagToken {
		switch t.Data {
		case "blockquote":
			// fmt.Println("nextToken: setting skip")
			mt.skip = true
		}
	}

	// Should we clear skip?
	if t.Type == html.EndTagToken {
		switch t.Data {
		case "blockquote":
			// fmt.Println("nextToken: clearing skip")
			mt.skip = false
		}
	}
}

func (mt tachyons) getClasses(tt html.TokenType, t html.Token, orgClasses string) (bool, string) {

	// Are we in skip mode? Not all tokens should have classes, they can override some we just higher in the tree
	if mt.skip {
		return false, ""
	}

	// fmt.Printf("%#v\n%q\n", t, t.Type.String())

	switch orgClasses {
	case "":
		return false, ""
	case "sect2":
		return true, fmt.Sprintf("%s %s", orgClasses, "f4 f2-ns lh-copy measure sans-serif")

	case "quoteblock":
		return true, fmt.Sprintf("%s %s", orgClasses, "f6 f5-ns lh-copy measure i bl bw1 b--gold mb4")

	case "paragraph lead":
		return true, fmt.Sprintf("%s %s", orgClasses, "f5 f3-ns lh-copy measure georgia")
	case "paragraph":
		return true, fmt.Sprintf("%s %s", orgClasses, "f5 f4-ns lh-copy measure mb4 georgia")
	default:
		return true, orgClasses
	}

}

func main() {
	fmt.Println(fooPath("test-1"))
}

func fooPath(dirpath string) error {

	err := os.Chdir(dirpath)
	if err != nil {
		return fmt.Errorf("fooPath: %s", err)
	}

	// read metadata file
	blogmetadata, err := loadMetadata("metadata.toml")
	if err != nil {
		return fmt.Errorf("fooPath: %s", err)
	}

	file, err := os.Open("content.html")
	if err != nil {
		return fmt.Errorf("fooPath: %s", err)
	}

	blogBody2, err := postprocess(file)
	if err != nil {
		return fmt.Errorf("fooPath: %s", err)
	}

	// fmt.Println(blogMetadata.title)

	b, err := ioutil.ReadAll(blogBody2)
	if err != nil {
		return fmt.Errorf("fooPath: %s", err)
	}

	t := template.New("blogpost")

	template.Must(t.Parse(templateTest))

	tInput := struct {
		BlogpostDate     string
		BlogpostTitle    string
		BlogpostSubtitle string
		BlogpostAuthor   string
		BlogpostBody     template.HTML // Unsafe / unencoded. Input must be safe, a it is here since it comes from ascidoc(tor)
	}{
		blogmetadata.Postmeta.Date.String(),
		blogmetadata.Postmeta.Title,
		blogmetadata.Postmeta.Subtitle,
		blogmetadata.Postmeta.Author,
		template.HTML(string(b)),
	}
	err = t.Execute(os.Stdout, tInput)
	if err != nil {
		return fmt.Errorf("fooPath: %s", err)
	}

	return nil
}

func postprocess(input io.Reader) (io.Reader, error) {

	myTachyons := tachyons{}

	z := html.NewTokenizer(input)

	var body strings.Builder
	var err error

	for {

		// Advance to next token
		tt := z.Next()
		if tt == html.ErrorToken {
			// This includes EOF, break out and deal with it later
			err = z.Err()
			break
		}

		t := z.Token()

		myTachyons.nextToken(t)

		if tt == html.StartTagToken {

			// Find and save the original classes, if any.
			var orgClasses string // class string
			oci := -1             // index to remove later
			for i := range t.Attr {
				if t.Attr[i].Key == "class" {
					orgClasses = t.Attr[i].Val
					oci = i
					break // assume just 1 class
				}
			}

			// Remove the original class attr. Assume order of is irrelevant
			if oci >= 0 {
				t.Attr[oci] = t.Attr[0]
				t.Attr = t.Attr[1:]
				// t.Attr[oci] = t.Attr[len(t.Attr)-1] // or just t.Attr[0] = t.Attr[oci] , t.Attr = t.Attr[1:] ??
				// t.Attr = t.Attr[:len(t.Attr)-1]
			}

			add, classes := myTachyons.getClasses(tt, t, orgClasses)

			if add {
				t.Attr = append(t.Attr, html.Attribute{Key: "class", Val: classes})
			}
		}

		body.WriteString(t.String())
	}

	// Any parse / state machine error from?
	if err != nil {
		if err != io.EOF {
			return nil, fmt.Errorf("postprocess: error when replacing: %s", err)
		}
		err = nil
	}

	return strings.NewReader(body.String()), nil
}

func findAttr(attrs []html.Attribute, key, val string) (bool, int) {
	for i := range attrs {
		if attrs[i].Key == key {
			if attrs[i].Val == val {
				return true, i // assume only 1 match
			}
		}
	}
	return false, 0
}
