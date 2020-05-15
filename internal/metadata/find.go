package metadata

import (
	"fmt"
	"os"
	"path/filepath"
)

func Find(startpath, metadataFilename string) ([]string, error) {

	res := []string{}

	err := filepath.Walk(startpath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walk: %w for path %q", err, path)
		}
		if info.Mode().IsRegular() && info.Name() == metadataFilename {

			// we preffer absolute paths, easier to work with
			abs, err := filepath.Abs(path)
			if err != nil {
				return fmt.Errorf("walk: abs: %w for path %q", err, path)
			}

			// ignore filename
			p := filepath.Dir(abs)
			if len(p) < 3 {
				return fmt.Errorf("walk: len of path %q is suspiciously short", p)
			}

			res = append(res, p)
		}
		return nil
	})

	if err != nil {
		return res, fmt.Errorf("findMetadataFiles: %w", err)
	}

	return res, nil
}
