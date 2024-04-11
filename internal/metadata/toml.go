package metadata

import (
	"fmt"
	"os"
	"time"

	"github.com/pelletier/go-toml/v2"
)

type M struct {
	Author            string
	Title             string
	Date              time.Time
	Servepath         string
	Published         bool
	Thisformatversion int // metameta.. the format of the metdatada file and this type

	// Only in blogposts
	BlogSubtitle string
	BlogTags     []string
}

func Fromfile(filepath string) (M, error) {

	var res M

	file, err := os.Open(filepath)
	if err != nil {
		return res, fmt.Errorf("Fromfile: path %q failed: %w", filepath, err)
	}

	dec := toml.NewDecoder(file)

	err = dec.Decode(&res)
	if err != nil {
		return res, fmt.Errorf("Fromfile: path %q failed: %w", filepath, err)
	}

	// sanity check
	if len(res.Title) < 3 {
		return res, fmt.Errorf("Fromfile: title was suspiciously short for: %q", filepath)
	}

	return res, nil
}
