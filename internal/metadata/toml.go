package metadata

import (
	"fmt"
	"os"
	"time"

	"github.com/pelletier/go-toml"
)

type BlogPost struct {
	Author   string
	Title    string
	Subtitle string
	Date     time.Time
}

type BlogEngine struct {
	Servepath string
	Published bool
	Tags      []string
}

type Metameta struct {
	Thisformatversion int
}

type Blogpost struct {
	Postmeta   BlogPost
	Enginemeta BlogEngine
	Metameta   Metameta
}

func LoadMetadata(filepath string) (Blogpost, error) {

	var res Blogpost

	file, err := os.Open(filepath)
	if err != nil {
		return res, fmt.Errorf("loadMetadata: path %q failed: %w", filepath, err)
	}

	dec := toml.NewDecoder(file)

	err = dec.Decode(&res)
	if err != nil {
		return res, fmt.Errorf("loadMetadata: path %q failed: %w", filepath, err)
	}

	return res, nil
}
