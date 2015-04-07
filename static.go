package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

//  Was:
//
//      ab -n 5000 -c 5 localhost:84/static/css.css
//
//      Time taken for tests:   0.888 seconds
//      Complete requests:      5000
//      Failed requests:        0
//      Requests per second:    5632.20 [#/sec] (mean)
//      Time per request:       0.888 [ms] (mean)
//      Time per request:       0.178 [ms] (mean, across all concurrent requests)
//      Transfer rate:          6578.24 [Kbytes/sec] received
//
//      Connection Times (ms)
//                    min  mean[+/-sd] median   max
//      Connect:        0    0   0.0      0       1
//      Processing:     0    1   0.3      1       4
//      Waiting:        0    1   0.3      1       4
//      Total:          0    1   0.3      1       4
//
//      Percentage of the requests served within a certain time (ms)
//        50%      1
//        66%      1
//        75%      1
//        80%      1
//        90%      1
//        95%      1
//        98%      2
//        99%      2
//       100%      4 (longest request)



// read all files in static/
// store in mem as string => []byte map ?
// make a reader
// serve with ServeContent
// detect changes :)

type StatucFile struct {
	ContentRaw bytes.Buffer
	Path       string
}

func visit(path string, f os.FileInfo, err error) error {
	fmt.Printf("Visited: %s\n", path)

	// open the file
	//  file, err := os.Open(dir + "/" + f.Name()) // fixme?
	//  defer file.Close()
	//  if err != nil {
	//  	checkAndWarn("Open file", err)
	//  	continue
	//  }

	//  // Use a scanner to read the page file line-by-line
	//  scanner := bufio.NewScanner(file)

	//  // parse the front matter
	//  frontMatter, err := parseFrontmatter(scanner)
	//  if err != nil {
	//  	checkAndWarn("Frontmatter for file: "+f.Name(), err)
	//  	continue
	//  }

	//  // Create a page
	//  page := Page{
	//  	Title:   frontMatter["Title"],
	//  	Path:    frontMatter["Path"],
	//  	PubTime: time.Now(), //FIXME
	//  }

	return nil

}

type InMemoryFile struct {
	at   int64
	Name string
	data []byte
	//fs   InMemoryFS
}

type ReadSeeker interface {
	Reader
	Seeker
}
type Reader interface {
	Read(p []byte) (n int, err error)
}
type Seeker interface {
	Seek(offset int64, whence int) (int64, error)
}

func (f *InMemoryFile) Seek(offset int64, whence int) (int64, error) {
    log.Println("Seek")
    return 0, nil
}

func (f *InMemoryFile) Read(p []byte) (n int, err error) {
    log.Println("Read")
    p = []byte("Hello2")
    return 6, nil
}

func main() {
	flag.Parse()
	root := flag.Arg(0)
	err := filepath.Walk(root, visit)
	fmt.Printf("filepath.Walk() returned %v\n", err)

        b := []byte("Hello Wordl!")
        reader := bytes.NewReader(b)

	http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		//file, err := os.Open("static/css.css")
		//if err != nil {
		//	log.Fatal("open() : ", err)
		//}

		_ = &InMemoryFile{
			at:   0,
			Name: "test",
			data: []byte("Hello Wordl!"),
		}


		http.ServeContent(w, r, "css.css", time.Now(), reader)
	})

	log.Println("Listening...")
	log.Fatal(http.ListenAndServe(":84", nil))
}
