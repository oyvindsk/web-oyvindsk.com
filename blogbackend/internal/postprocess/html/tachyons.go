package html

import (
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
)

// These functions do html postprocessing to replace classes from the Asciidoc(tor) html output with Tachyons classes.
// http://tachyons.io/
// https://code.luasoftware.com/tutorials/web-development/tachyon-css-cheatsheet/

// InsertTachyonsClasses looks at the html from AsciiDoc(tor) and the classes specified and inserts tachyons css classes to get the layout we want
func InsertTachyonsClasses(input io.Reader) (io.Reader, error) {

	myTachyons := tachyons{} // used to remember some state

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

		thisToken := z.Token()

		myTachyons.nextToken(thisToken)

		if tt == html.StartTagToken {

			//
			// Classes

			// Find and save the original classes, if any.
			// Remove the original class attr
			var orgClasses string // class string
			if ok, oci := findAttr(thisToken.Attr, "class"); ok {
				orgClasses = thisToken.Attr[oci].Val

				// Remove the original class attr. Assume order of is irrelevant
				thisToken.Attr[oci] = thisToken.Attr[0]
				thisToken.Attr = thisToken.Attr[1:]

				// t.Attr[oci] = t.Attr[len(t.Attr)-1] // or just t.Attr[0] = t.Attr[oci] , t.Attr = t.Attr[1:] ??
				// t.Attr = t.Attr[:len(t.Attr)-1]
			}

			add, classes := myTachyons.tachyonsClasses(tt, thisToken, orgClasses)

			if add {
				thisToken.Attr = append(thisToken.Attr, html.Attribute{Key: "class", Val: classes})
			}
		}

		body.WriteString(thisToken.String())
	}

	// Any parse / state machine error from?
	if err != nil {
		if err != io.EOF {
			return nil, fmt.Errorf("InsertTachyonsClasses: error when replacing: %s", err)
		}
		err = nil
	}

	return strings.NewReader(body.String()), nil
}

// keep state while advancing ourtokens
type tachyons struct {
	skip bool // should this token have tachyons classes added?
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

func (mt tachyons) tachyonsClasses(tt html.TokenType, t html.Token, orgClasses string) (bool, string) {

	// Are we in skip mode? Not all tokens should have classes, they can override some we just higher in the tree
	if mt.skip {
		return false, ""
	}

	// fmt.Printf("%#v\n%q\n", t, t.Type.String())

	// Headers:
	// Title (h1 from asciidoc, in templates): f2 f1-m f-headline-l
	commonH := "measure lh-title" // for all headers
	switch t.Data {
	case "h2": // h2 == (sect1 from asciidoc)
		return true, fmt.Sprintf("%s %s", commonH, "f3 f2-m f1-l")
	case "h3": // h3 ===
		return true, fmt.Sprintf("%s %s", commonH, "f4 f3-m f2-l mv0")
	case "h4": // h4 ====
		return true, fmt.Sprintf("%s %s", commonH, "f5 f4-m f3-l mv0")
	case "h5": // h5 =====
		return true, fmt.Sprintf("%s %s", commonH, "f6 f5-m f4-l mv0")
	case "h6": // h6 ======
		return true, fmt.Sprintf("%s %s", commonH, "f7 f6-m f5-l mv0")
	}

	// Paragraphs and quote blocks
	switch orgClasses {

	case "quoteblock":
		return true, fmt.Sprintf("%s %s", orgClasses, "f6 f5-ns lh-copy measure i bl bw1 b--gold mb4")

	case "paragraph lead":
		return true, fmt.Sprintf("%s %s", orgClasses, "f5 f3-ns lh-copy measure georgia")

	case "paragraph":
		return true, fmt.Sprintf("%s %s", orgClasses, "f5 f4-ns lh-copy measure mb4 georgia")

	case "ulist": // Unordered list, ul
		return true, fmt.Sprintf("%s %s", orgClasses, "f5 f4-ns lh-copy measure mb4 georgia")

	}

	// Code blocks - Use a little more complicated matching than the others
	if t.Data == "code" && strings.HasPrefix(orgClasses, "language-") {
		return true, fmt.Sprintf("%s %s", orgClasses, "bg-washed-green f6 f5-ns code")
	}
	if t.Data == "pre" && orgClasses == "highlight" {
		return true, fmt.Sprintf("%s %s", orgClasses, "lh-solid") // outher code block element
	}

	// Default, if haven't matched anything
	return true, orgClasses // fmt.Sprintf("%s %s", orgClasses, "f5 f4-ns lh-copy measure georgia")
}
