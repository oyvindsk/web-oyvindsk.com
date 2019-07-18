package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// readFile reads the metadata (front matter) and template content from a file.
// If there is no metadata, it's not a valid file
// inspired by https://github.com/PuerkitoBio/trofaf/blob/master/tpldata.go
//
// Returns a map with metadata, the unparsed template and an error
func readFile(filep string) (map[string]string, bytes.Buffer, error) {

	// Only handel files with the .html extention
	if filepath.Ext(filep) != ".html" {
		return nil, bytes.Buffer{}, fmt.Errorf("readFile: not a .html file: %q", filep)
	}

	// Open the file
	file, err := os.Open(filep)
	defer file.Close()
	if err != nil {
		return nil, bytes.Buffer{}, err
	}

	// Page metadata. We expect to find this in the  the start of the file between ---
	md := make(map[string]string)

	scanner := bufio.NewScanner(file)
	inm := false // in metadata (front matter)

	for scanner.Scan() {
		l := strings.TrimSpace(scanner.Text())

		// The metadata delimited by 3 dashes
		if l == "---" {

			if inm {
				// This signals the end of the metda data. Success!
				break
			}

			// This is the start of the metadata
			inm = true

		} else if inm {

			sections := strings.SplitN(l, ":", 2)
			if len(sections) != 2 {
				return nil, bytes.Buffer{}, errors.New("Invalid metadata line")
			}

			// Ok, looks like a valid metadata line
			md[strings.TrimSpace(sections[0])] = strings.TrimSpace(sections[1])

		} else if l != "" {
			// No metadata, quit
			return nil, bytes.Buffer{}, errors.New("No metadata")
		}
	}

	if err := scanner.Err(); err != nil {
		// The scanner stopped because of an error
		return nil, bytes.Buffer{}, err
	}

	// Parse the rest of the file, the stuff after metadata
	// basically just read evrything left into a buffer
	var content bytes.Buffer
	for scanner.Scan() {
		content.WriteString(scanner.Text() + "\n")
	}

	if err := scanner.Err(); err != nil {
		return nil, bytes.Buffer{}, err
	}

	return md, content, nil
}
