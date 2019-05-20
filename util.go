package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// checkError checks if the error exists, and if so, outputs it to stdout and panics
func checkError(err error) {
	if err != nil {
		fmt.Print(err)
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
	return "", errors.New("No main file found")
}

// lines will divide a file's contents up into lines and put them into a slice
func lines(file *os.File) []string {
	// set the capacity to 32 lines
	result := make([]string, 0, 32)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}
	return result
}
