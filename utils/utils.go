package utils

import (
	"os"
	"path/filepath"
)

func GetDirs(dirPth string) []string {
	var files []string

	root := dirPth
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return files
}
