package main

import (
	"errors"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"
)

// checkError checks if the error exists, and if so, outputs it to stdout and panics
func checkError(err error) {
	if err != nil {
		// fmt.Print(err)
		panic(err)
	}
}

// findMainFile takes in a directory path and finds a file of any extension with the name 'main.*'
// and returns the path to that file.
// If multiple files exist that are called 'main', then it only returns the first one alphabetically by the extension
func findMainFile(dirname string) (string, error) {
	files, err := ioutil.ReadDir(dirname)
	checkError(err)
	for _, info := range files {
		name := info.Name()
		base := filepath.Base(name)
		if strings.HasPrefix(base, "main.") {
			return path.Join(dirname, base), nil
		}
	}
	return "", errors.New("no main file found")
}
