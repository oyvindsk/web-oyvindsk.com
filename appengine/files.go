package blog

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type fileBasedContent interface {
	loadFrontMatter(s *bufio.Scanner) (err error)
	loadContent(s *bufio.Scanner) error
	parseAndExecuteTemplates() error
}

// Read the front matter and content from a file.
// If there is no front matter, it's not a valid file
// based on  https://github.com/PuerkitoBio/trofaf/blob/master/tpldata.go

func readFile(filepath string) (map[string]string, []byte, error) {

	// Open the file
	file, err := os.Open(filepath)
	defer file.Close()
	if err != nil {
		return nil, nil, err
	}

	scanner := bufio.NewScanner(file)

	// Parse the front amtter
	fm := map[string]string{}
	infm := false
	for scanner.Scan() {
		l := strings.Trim(scanner.Text(), " ")
		if l == "---" { // The front matter is delimited by 3 dashes
			if infm {
				// This signals the end of the front matter
				// Success!
				break
			}
			// This is the start of the front matter
			infm = true

		} else if infm {
			sections := strings.SplitN(l, ":", 2)
			if len(sections) != 2 {
				// Invalid front matter line
				return nil, nil, errors.New("Invalid front matter line")
			}
			// Ok, looks like a valid front matter line
			fm[strings.Trim(sections[0], " ")] = strings.Trim(sections[1], " ")
		} else if l != "" {
			// No front matter, quit
			return nil, nil, errors.New("No front matter")
		}
	}
	if err := scanner.Err(); err != nil {
		// The scanner stopped because of an error
		return nil, nil, err
	}

	//return nil, nil, errors.New("Empty post file")

	// Parse the rest of the blogpost, the stuff after front-matter
	// basically just read evrything left into a buffer

	content := bytes.NewBuffer(nil)
	for scanner.Scan() {
		content.WriteString(scanner.Text() + "\n")
	}
	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	return fm, content.Bytes(), nil
}

// Filter out the files we actually want to read  - From http://0value.com/build-a-blog-engine-in-Go
func fileFilterHTML(files []os.FileInfo) []os.FileInfo {
	for i := 0; i < len(files); {
		if files[i].IsDir() || filepath.Ext(files[i].Name()) != ".html" {
			files[i], files = files[len(files)-1], files[:len(files)-1] // replace the one we don't want with the last one, and shorten our slice by 1
		} else {
			i++
		}
	}
	return files
}
