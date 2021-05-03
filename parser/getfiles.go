package parser

import (
	"os"
	"path/filepath"
	"strings"
)

//GetFiles returns all files in a directory that end in .tusk
func GetFiles(dir string) (files []string, err error) {

	var tfiles []string

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})

	for _, v := range tfiles {
		if strings.HasPrefix(v, ".tusk") {
			files = append(files, v)
		}
	}

	return
}
