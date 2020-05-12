package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"regexp"
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

	file, err := os.Open("blog-asciidoctor.html")
	if err != nil {
		fmt.Println(err)
		return
	}

	blogBody1, blogMetadata, err := cutMetadata(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	blogBody2, err := postprocess(blogBody1)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(blogMetadata.title)
	b, err := ioutil.ReadAll(blogBody2)
	if err != nil {
		fmt.Println(err)
		return
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
		blogMetadata.date,
		blogMetadata.title,
		blogMetadata.subtitle,
		blogMetadata.author,
		template.HTML(string(b)),
	}
	err = t.Execute(os.Stdout, tInput)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func cutMetadata(input io.Reader) (io.Reader, blogMetadata, error) {

	z := html.NewTokenizer(input)

	/*
		States at the beginning of the big for loop:

		Current State								Input that switches state 			Next state
		----------------------------------------------------------------------------------------------------------
		(lfop) Looking for opening token			div token with openblock class 		lfml
		(lfml) Looking for magic lines				line that starts with "||"			iml
		(iml)  In magic								line that does not start with ..	lfet
		(lfet) Looking for tokens we should exclude last tokens in exclude list			done
		(done) Done looking and exluding			n/a									n/a
	*/

	// state machine variables
	state := "lfot"                 // current state
	var prevToken html.Token        // we sometimes have to look back to the previous token
	var endTokensToExlcude []string // usen when excluding tokens around the metadata lines. Typically 3 tokens before and 3 after.

	// results of the state machine loop (not the func)
	var magicLines []string
	var body strings.Builder
	var err error

MACHINE:
	for {

		// Advance to next token
		tt := z.Next()
		if tt == html.ErrorToken {
			// This includes EOF, break out and deal with it later
			err = z.Err()
			break MACHINE
		}

		thisToken := z.Token() // The token we are currenlty looking at, as opposed to prevToken

		// Switch on the 5 known states. See above.
		// this could of course be something other than a string, otoh ..
		// we do not really enforce that all transitions are valid, but that would require a bug in the code (?)
		switch state {

		case "lfot":

			// Look for opening div of metadata, with class openblock
			var found bool
			if tt == html.StartTagToken && thisToken.Data == "div" {
				if found, _ = findAttr(thisToken.Attr, "class", "openblock"); found {
					state = "lfml" // fmt.Println("\n\t==>\t Looking for magic lines")

					// add this div to the list of tokens we want to exlude after the magic lines (in lfet)
					endTokensToExlcude = append(endTokensToExlcude, thisToken.Data)
				}
			}

			// Include token unless it was the opening div we are looking for
			if !found {
				body.WriteString(thisToken.String())
			}

		case "lfml":

			if thisToken.Type.String() == "Text" && strings.HasPrefix(thisToken.Data, "||") {
				state = "iml" // fmt.Println("\n\t==>\t In magic")
				break
			}

			// Add tokens we see before the firts line of magic
			// to the list of tokens we want to exlude after the magic lines (in lfet)
			if thisToken.Type == html.StartTagToken {
				endTokensToExlcude = append(endTokensToExlcude, thisToken.Data)
			}

		case "iml":

			// Save the magic lines(s) for later
			// syntax from ascidoc(tor) puts it on 1 line with a \n, so ..
			magicLines = append(magicLines, strings.Split(prevToken.String(), "\n")...)

			if thisToken.Type.String() != "Text" || !strings.HasPrefix(thisToken.Data, "||") {
				state = "lfet" // fmt.Printf("\n\t==>\t Looking for tags we should exclude\n")
			}

		case "lfet":

			if prevToken.Type == html.EndTagToken && prevToken.Data == endTokensToExlcude[len(endTokensToExlcude)-1] {
				endTokensToExlcude = endTokensToExlcude[:len(endTokensToExlcude)-1]
			}

			if len(endTokensToExlcude) == 0 {
				state = "done" //	fmt.Println("\n\t==>\t DONE!")
			}

		case "done":
			body.WriteString(thisToken.String())

		default:
			err = fmt.Errorf("unknown state seen: %q", state)
			break MACHINE
		}

		prevToken = thisToken
	}

	// Any parse / state machine error from?
	if err != nil {
		if err != io.EOF {
			return nil, blogMetadata{}, fmt.Errorf("cutMetadata: error when running state machine: %s", err)
		}
		err = nil
	}

	// Convert the magic lines we found into blogMetadata
	metadata, err := blogMetadataFromMagicLines(magicLines)
	if err != nil {
		return nil, blogMetadata{}, fmt.Errorf("cutMetadata: %s", err)
	}

	return strings.NewReader(body.String()), metadata, nil
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

type blogMetadata struct {
	author    string
	title     string
	subtitle  string
	date      string
	servePath string
	tags      []string
}

// Input looks like this:
// "|| Adam Morse || Too many tools and frameworks: subTTT || 2015 || /foo/bar || Subtitle: The definitive guide to the javascript tooling landscape in 2015"
// "|| foo bar go golang javascript"
func blogMetadataFromMagicLines(magicLines []string) (blogMetadata, error) {
	if !(len(magicLines) == 1 || len(magicLines) == 2) {
		return blogMetadata{}, fmt.Errorf("blogMetadataFromMagicLines: Expect 1 or 2 magix lines, got: %d", len(magicLines))
	}

	// First line, || separated, everything but the tags
	l1 := regexp.MustCompile(`\s?\|\|\s?`).Split(magicLines[0], 100)
	l1 = l1[1:] // first is always bogus since we start out line with ||

	m := blogMetadata{
		author:    l1[0],
		title:     l1[1],
		subtitle:  l1[4],
		date:      l1[2],
		servePath: l1[3],
	}

	// add tags if any
	if len(magicLines) > 1 && len(magicLines[1]) > 4 {
		if !strings.HasPrefix(magicLines[1], "|| ") {
			return blogMetadata{}, fmt.Errorf("blogMetadataFromMagicLines: Tag line invalid, must start with '|| '")
		}

		//l2 := regexp.MustCompile(`\|?\|?\s`).Split(magicLines[1], 100)
		m.tags = strings.Fields(magicLines[1][3:]) // split on space after '|| '
	}

	return m, nil
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
