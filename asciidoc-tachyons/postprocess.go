package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/net/html"
)

const header = `<!DOCTYPE html>
<html lang="en">
	<title> </title>
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<link rel="stylesheet" href="https://unpkg.com/tachyons/css/tachyons.min.css">
  
	<body>
		<article>
  			<header class="bg-gold sans-serif">
    			<div class="mw9 center pa4 pt5-ns ph7-l">
      				<time class="f6 mb2 dib ttu tracked"><small>...DATE...</small></time>
      				<h3 class="f2 f1-m f-headline-l measure-narrow lh-title mv0">
        				<span class="bg-black-90 lh-copy white pa1 tracked-tight">
          					...TITLE...
        				</span>
      				</h3>
      				<h4 class="f3 fw1 georgia i">...SUBTITLE...</h4>
      				<h5 class="f6 ttu tracked black-80">By ..AUTHOR..</h5>
    			</div>
  			</header>

			<div class="pa4 ph7-l georgia mw9-l center">
  
				<!-- auto generert herfrra -->
`

const footer = `
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

	// <div class="openblock">
	// <div class="content">
	// <div class="paragraph">
	// <p>|| Adam Morse || Too many tools and frameworks: subTTT || 2015 || /foo/bar || Subtitle: The definitive guide to the javascript tooling landscape in 2015 ||</p>
	// </div>
	// </div>
	// </div>

	file, err := os.Open("doctor.html")
	if err != nil {
		fmt.Println(err)
		return
	}

	z := html.NewTokenizer(file)

	var exclude bool
	var excludeAfter = []string{}
	var foundString bool

	/*
		States at the beginning of the big for loop:

		exclude			excludeAfter		foundString 		MEANING
		0 				0 					0					Looking for start div
		0 				0 					1					Done excluding!
		0 				1 					0					N/A : Turned on together, exclude turns off when !excludeAfter
		0 				1 					1					N/A : Turned on together, exclude turns off when !excludeAfter

		1 				0 					0					N/A : Turned on together, excludeAfter is not turned off when !foundString
		1 				0 					1					Emptied excludeAfter but haven't turned exclude off yet. BUG??
		1 				1 					0					Looking for foundString. Excluding and saving tags in excludeAfter
		1 				1 					1					Looking for end div, removing tags from excludeAfter
	*/
	var excludedCnt int

	for {
		tt := z.Next()

		if tt == html.ErrorToken {
			fmt.Println(z.Err().Error()) // FIXME
			break                        // return
		}

		t := z.Token()

		// Look for metadata if we haven't found it yet
		if !foundString && !exclude {
			if tt == html.StartTagToken && t.Data == "div" {
				for i := range t.Attr {
					if t.Attr[i].Key == "class" {
						if t.Attr[i].Val == "openblock" {
							exclude = true
							excludeAfter = append(excludeAfter, "div")
							fmt.Println("EXCLUDE ON")
							break // assume just 1 class
						}
					}
				}
			}
		} else if exclude {

			fmt.Printf("--   %q   %q\n", t.Data, t.Type.String())

			if foundString {

				fmt.Printf("%#v\n", excludeAfter)

				if len(excludeAfter) == 0 {
					fmt.Println("EXCLUDE OFF")
					exclude = false
				} else {
					if tt == html.EndTagToken && t.Data == excludeAfter[len(excludeAfter)-1] {
						excludeAfter = excludeAfter[:len(excludeAfter)-1]
					}
				}

			}

			if !foundString {

				if tt == html.StartTagToken {
					excludeAfter = append(excludeAfter, t.Data)
				}

				if t.Type.String() == "Text" && strings.HasPrefix(t.Data, "||") {
					fmt.Println("FOUND")
					foundString = true
				}
			}
		}

		if exclude {
			excludedCnt++
			if excludedCnt > 20 {
				fmt.Println("Error: Excluded too much from blogpost wehn looking for metadata, giving up!")
				return // FIXME
			}
		} else {
			fmt.Print(t.String())
		}
		// FIXME
	}

}

func postprocess() {

	myTachyons := tachyons{}

	file, err := os.Open("blogpost.html")
	if err != nil {
		fmt.Println(err)
		return
	}

	z := html.NewTokenizer(file)

	fmt.Print(header) // FIXME

	for {
		tt := z.Next()

		if tt == html.ErrorToken {
			fmt.Println(z.Err().Error()) // FIXME
			break                        // return
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

		fmt.Print(t.String()) // FIXME
	}
	fmt.Print(footer) // FIXME

}
