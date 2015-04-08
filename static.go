package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"path"
	"time"
)

// read all files in static/
// store in mem as string => []byte map ?
// make a reader
// serve with ServeContent
// detect changes :)


// Normal implementation - FileServer 




type StaticFile struct {
	PubTime time.Time
	ContentRaw []byte
	Path       string
}

var files map[string]StaticFile

func readFile(p string, f os.FileInfo, err error) error {
    base := path.Base(p)
	log.Printf("Visited: %s, shortpath: %s\n", p, base)

    // stat the file to skip directories etc
    info, err := os.Stat(p)
    if err != nil {
        checkAndWarn("readFile stat", err)
        return err
    }

    // skip non-regular files
    if info.Mode().IsRegular() == false {
	    log.Println("    ignored non-regular file: ", p)
        return nil
    }

    // slurp the whole file
    fileCont, err := ioutil.ReadFile(p)
    if err != nil {
        checkAndWarn("readFile ReadFile", err)
        return err
    }

    // store it
    files[base] = StaticFile{
        PubTime: time.Now(), // FIXME
        Path: p,
        ContentRaw: fileCont,
    }

	return nil
}

func readStaticFiles(dir string) (map[string]StaticFile, error) {

    files = make(map[string]StaticFile)

    // walk all the files recursivly
	err := filepath.Walk("static", readFile)
    checkAndDie("file walk in dir: " + dir, err)

    return files, err
}
