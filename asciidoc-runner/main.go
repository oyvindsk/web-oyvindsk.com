package main

import (
	"fmt"
	"os"

	"github.com/oyvindsk/web-oyvindsk.com/internal/metadata"
)

// We expect all files to have special names, too much work to pass these along :[]
const (
	metadataFilename   = "metadata.toml" // input
	asciidocFilename   = "content.adoc"  // input
	htmlFilename       = "content.html"  // output
	docBookFilename    = "full.xml"      // output
	hashFilenameSuffix = ".hash"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Expects 1 argument: a root path")
		os.Exit(1)
	}

	metadataPaths, err := metadata.Find(os.Args[1], metadataFilename)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, fp := range metadataPaths {
		if !filesChanged(fp) {
			fmt.Printf("Skipping %q, hash match!\n", fp)
			continue
		}

		err = generateAll(fp)
		if err != nil {
			fmt.Println(err)
			return
		}

		updateHashes(fp) // ignorer errors
	}

	fmt.Println("done!")
}
