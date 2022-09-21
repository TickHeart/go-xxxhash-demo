package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
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
	for _, file := range files {
		fmt.Println(root, "==========", file)
	}
	return files
}

func GetYamlFiles(dirPth string) []string {
	var files []string

	root := dirPth

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			matchString, _ := regexp.MatchString("\\.yaml$", info.Name())
			if matchString {
				files = append(files, path)
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		fmt.Println(root, "==========", file)
	}
	return files
}
