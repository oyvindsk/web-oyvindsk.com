// Package html does various html postprocessing
// It is therefore a part of the styling for the website, along with the templates
// It is not used with the pdf output (duh)
package html

import (
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/net/html"
)

const (
	// Images - we replace the filepaths of static files (works with pdfs) with a relative url
	staticRelRoot    = "../../../static_files/"
	staticWebRelRoot = "/"
)

// PostprocessFile opens a file on disk and runs all the normal postprocessing
// returns an error or the final html as a string
func PostprocessFile(filepath string) (string, error) {

	file, err := os.Open(filepath)
	if err != nil {
		return "", fmt.Errorf("PostprocessFile: %w", err)
	}

	body1, err := InsertContactForm(file)
	if err != nil {
		return "", fmt.Errorf("PostprocessFile FOO: %w", err)
	}

	body2, err := ReplaceOtherHTML(body1)
	if err != nil {
		return "", fmt.Errorf("PostprocessFile FOO2: %w", err)
	}

	bodyr, err := InsertTachyonsClasses(body2)
	if err != nil {
		return "", fmt.Errorf("PostprocessFile: %w", err)
	}

	body, err := ioutil.ReadAll(bodyr)
	if err != nil {
		return "", fmt.Errorf("PostprocessFile: %w", err)
	}

	return string(body), nil
}

// package internal helper funcs

func findAttr(attrs []html.Attribute, key string) (bool, int) {
	for i := range attrs {
		if attrs[i].Key == key {
			return true, i // assume only 1 match
		}
	}
	return false, 0
}

func findAttrVal(attrs []html.Attribute, key, val string) (bool, int) {
	for i := range attrs {
		if attrs[i].Key == key {
			if attrs[i].Val == val {
				return true, i // assume only 1 match
			}
		}
	}
	return false, 0
}
