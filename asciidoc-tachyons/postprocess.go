package main

import (
	"fmt"
	"os"

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
      				<time class="f6 mb2 dib ttu tracked"><small>27 July, 2015</small></time>
      				<h3 class="f2 f1-m f-headline-l measure-narrow lh-title mv0">
        				<span class="bg-black-90 lh-copy white pa1 tracked-tight">
          					Too many tools and frameworks
        				</span>
      				</h3>
      				<h4 class="f3 fw1 georgia i">The definitive guide to the javascript tooling landscape in 2015.</h4>
      				<h5 class="f6 ttu tracked black-80">By Adam Morse</h5>
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

	myTachyons := tachyons{}

	file, err := os.Open("blogpost.html")
	if err != nil {
		fmt.Println(err)
		return
	}

	z := html.NewTokenizer(file)

	fmt.Print(header)

	for {
		tt := z.Next()

		if tt == html.ErrorToken {
			fmt.Println(z.Err().Error())
			break // return
		}

		t := z.Token()

		myTachyons.nextToken(t)

		if tt == html.StartTagToken {

			// 				Purge all classes for Attr,

			// Find and cut class attr, if any. Assume order of attrs is irrelevant
			var orgClasses string
			for i := range t.Attr {
				if t.Attr[i].Key == "class" {
					orgClasses = t.Attr[i].Val
					t.Attr[i] = t.Attr[0]
					t.Attr = t.Attr[1:]
				}
			}

			add, classes := myTachyons.getClasses(tt, t, orgClasses)

			if add {
				t.Attr = append(t.Attr, html.Attribute{Key: "class", Val: classes})
			}
		}

		// for i := range t.Attr {
		// 	//fmt.Printf("%d\t%+v\t%#v\n%q\n", i, a, a, tt.String())
		// 	// tml.Attribute{Namespace:"", Key:"class", Val:"sect2"}
		// 	// html.Render(os.Stdout, )

		// 	if t.Attr[i].Key == "class" {
		// 		// t.Attr[i].Val = "FOO"
		// 		t.Attr[i] = html.Attribute{}
		// 	}
		// }

		fmt.Print(t.String())

		// switch {
		// case tt == html.ErrorToken:
		// 	// End of the document, we're done
		// 	return // return z.Err()
		// case tt == html.StartTagToken:
		// 	t := z.Token()

		// 	// fmt.Printf("%q\n", t.String())

		// 	isAnchor := t.Data == "div"
		// 	if isAnchor {
		// 		// fmt.Println("We found a link!")

		// 		for range t.Attr {
		// 			// fmt.Printf("%d\t%+v\t%#v\n%q\n", i, a, a, tt.String())
		// 			// html.Render(os.Stdout, )

		// 			t.Attr = []html.Attribute{}
		// 		}
		// 	}

		// 	//fmt.Printf("%q\n", z.Raw())

		// }
		// fmt.Printf("%s", z.Raw())
	}
	fmt.Print(footer)

}
