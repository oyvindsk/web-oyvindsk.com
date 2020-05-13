package main

import (
	"fmt"
	"os"
	"time"

	"github.com/pelletier/go-toml"
)

type Postmeta struct {
	Author   string
	Title    string
	Subtitle string
	Date     time.Time
}

type Enginemeta struct {
	Tags      []string
	Servepath string
}

type Metameta struct {
	Thisformatversion int
}

type BlogMetadata struct {
	Postmeta   Postmeta
	Enginemeta Enginemeta
	Metameta   Metameta
}

func loadMetadata(filepath string) (BlogMetadata, error) {

	var res BlogMetadata

	file, err := os.Open(filepath)
	if err != nil {
		return res, fmt.Errorf("loadMetadata: %s", err)
	}

	dec := toml.NewDecoder(file)

	err = dec.Decode(&res)
	if err != nil {
		return res, fmt.Errorf("loadMetadata: %s", err)
	}

	return res, nil
}
