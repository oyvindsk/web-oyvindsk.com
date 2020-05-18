package tachyons

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func PostprocessFile(filepath string) (string, error) {

	file, err := os.Open(filepath)
	if err != nil {
		return "", fmt.Errorf("PostprocessFile: %w", err)
	}

	bodyr, err := Postprocess(file)
	if err != nil {
		return "", fmt.Errorf("PostprocessFile: %w", err)
	}

	body, err := ioutil.ReadAll(bodyr)
	if err != nil {
		return "", fmt.Errorf("PostprocessFile: %w", err)
	}

	return string(body), nil
}

func Postprocess(input io.Reader) (io.Reader, error) {

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
