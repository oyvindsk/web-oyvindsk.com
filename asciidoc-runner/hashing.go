package main

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"io"
	"io/ioutil"
	"os"
)

// filesChanged looks in the dir given and computes the hash of two files:
//												metadata.toml 		and content.adoc
// and compare them to the hashes found in: 	metadata.toml.hash 	and content.adoc.hahs
// return true if they are not all equal OR there are any errors
// All of this is just an optimalization to avoid regenerating the same content all the time
// (also the pdf changes a little every time they are generated, anoying to have in the repo)
func filesChanged(dirpath string) bool {

	for _, fn := range []string{asciidocFilename, metadataFilename} {

		// Slurp the .hash file
		hashfile, err := ioutil.ReadFile(dirpath + "/" + fn + hashFilenameSuffix)
		if err != nil {
			if !os.IsNotExist(err) {
				fmt.Printf("filesChanged: error: %s\n", err)
			}
			return true
		}

		// Hash the content of the "real" file
		realfile, err := os.Open(dirpath + "/" + fn)
		if err != nil {
			fmt.Printf("filesChanged: error: %s\n", err)
			return true
		}
		defer realfile.Close() // ignore errors

		h := fnv.New128a()
		_, err = io.Copy(h, realfile)
		if err != nil {
			fmt.Printf("filesChanged: error: %s\n", err)
			return true
		}

		// Compare
		if bytes.Compare(hashfile, h.Sum(nil)) != 0 {
			return true
		}
	}
	return false
}

func updateHashes(dirpath string) {

	for _, fn := range []string{asciidocFilename, metadataFilename} {

		// Hash the content of the "real" file
		realfile, err := os.Open(dirpath + "/" + fn)
		if err != nil {
			fmt.Printf("updateHashes: error: %s\n", err)
			return
		}
		defer realfile.Close() // ignore errors

		h := fnv.New128a()
		_, err = io.Copy(h, realfile)
		if err != nil {
			fmt.Printf("updateHashes: error: %s\n", err)
			return
		}

		// Write a new hash file
		ioutil.WriteFile(dirpath+"/"+fn+hashFilenameSuffix, h.Sum(nil), 0600)
		if err != nil {
			fmt.Printf("updateHashes: error: %s\n", err)
			return
		}
	}
}
