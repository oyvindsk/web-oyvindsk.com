package tachyons

import (
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
)

// FIXME TODO new package?

const (
	magicParagraphExcludeText = "You can email me at"

	contactForm = `
	<!-- Contact form-->
	<!-- Also set the url in the js in template/layout.html -->
	<form method="post" action="/contact-form" class="contact-form" id="contactForm" novalidate="">
		<div>
			<input style="width:500px;" name="name" placeholder="Your name" type="text">
		</div>
		<div>
			<input style="width:500px;" name="from" placeholder="Your email / phone number" type="email">
		</div>
	
		<div>
			<input style="width:500px;" name="subject" placeholder="Subject" type="text">
		</div>
	
		<div>
			<textarea style="width:500px;" name="text" placeholder="Hi. I would like to have a friendly chat about.." rows="7"></textarea>
		</div>
	
		<div class="get-it">
			<button name="submitButton" type="submit">Send</button>
		</div>
	</form>
	
	<div id="messages">msg</div>
	`

	urlToRelativize = "https://oyvindsk.com/"
)

// FOO Look for a contact-me div block. Cut in and include html for a contact me form instead
// assumes there's only 1 concat form =/
func FOO(input io.Reader) (io.Reader, error) {

	// <div class="paragraph">
	//
	// <p>You can email me at
	//		<a href="mailto:hello@oyvindsk.com">hello@oyvindsk.com</a>
	//  	or use
	// 		<a href="https://oyvindsk.com/hire-me#contact" class="bare">https://oyvindsk.com/hire-me#contact</a>
	// </p>
	//
	// </div>

	z := html.NewTokenizer(input)

	/*
		States at the beginning of the big for loop:

		Current State								Input that switches state 			Next state
		----------------------------------------------------------------------------------------------------------
		(lft)  Looking for magic text               <p> containing text                 lfep
		(lfep) Looking for end of p					</p>								done
		(done) Done looking and exluding			n/a									n/a
	*/

	// state machine variables
	state := "lft" // current state

	// results of the state machine loop (not the func)
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

		thisToken := z.Token() // The token we are currenlty looking at

		// log.Printf("MACHINE: token: %q  %q\tstate: %s\n", thisToken.Type.String(), thisToken.Data, state)

		// Switch on the 3 known states. See above.
		// this could of course be something other than a string, otoh ..
		// we do not really enforce that all transitions are valid, but that would require a bug in the code (?)
		switch state {

		case "lft":

			if thisToken.Type.String() == "Text" && strings.Contains(thisToken.Data, magicParagraphExcludeText) {
				fmt.Println("\n\t==>\t Looking for end of p")
				state = "lfep"
			} else {
				body.WriteString(thisToken.String())
			}

		case "lfep":
			if tt == html.EndTagToken && thisToken.Data == "p" {
				fmt.Println("\n\t==>\t Done")
				state = "done"

				body.WriteString("DONE HERE:" + contactForm)
				body.WriteString(thisToken.String())
			}

		case "done":
			body.WriteString(thisToken.String())

		default:
			err = fmt.Errorf("unknown state seen: %q", state)
			break MACHINE

		}
	}

	// Any parse / state machine error from?
	if err != nil {
		if err != io.EOF {
			return nil, fmt.Errorf("FOO: error when running state machine: %s", err)
		}
	}

	return strings.NewReader(body.String()), nil
}

// FOO2 Look for other html to replace or "fix"
func FOO2(input io.Reader) (io.Reader, error) {

	z := html.NewTokenizer(input)

	// results of the state machine loop (not the func)
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

		thisToken := z.Token() // The token we are currenlty looking at

		// Replace urls linking to ourself witha relative url
		// we use the full in pages beacuse the relative ones get effd in the PDFs for some reason =/
		if thisToken.Type == html.StartTagToken && thisToken.Data == "a" {
			if ok, i := findAttr(thisToken.Attr, "href"); ok {
				if strings.Contains(thisToken.Attr[i].Val, urlToRelativize) {
					// log.Printf("MACHINE: token: %q  %q  %d  %s\n", thisToken.Type.String(), thisToken.Data, i, thisToken.Attr[i].Val)
					thisToken.Attr[i].Val = strings.Replace(thisToken.Attr[i].Val, urlToRelativize, "/", 1)
				}
			}
		}

		body.WriteString(thisToken.String())

	}

	// Any parse / state machine error from?
	if err != nil {
		if err != io.EOF {
			return nil, fmt.Errorf("FOO2: error when running state machine: %s", err)
		}
	}

	return strings.NewReader(body.String()), nil
}
