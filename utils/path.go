package utils

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

var basePath string

// BasePath return executable file directory
func BasePath() string {
	if basePath != `` {
		return basePath
	}
	filename, e := exec.LookPath(os.Args[0])
	if e != nil {
		log.Fatalln(e)
	}

	filename, e = filepath.Abs(filename)
	if e != nil {
		log.Fatalln(e)
	}
	basePath = filepath.Dir(filename)
	return basePath
}

// Abs Use bashPath as the working directory to return the absolute path
func Abs(bashPath, path string) string {
	if filepath.IsAbs(path) {
		path = filepath.Clean(path)
	} else {
		path = filepath.Clean(filepath.Join(bashPath, path))
	}
	return path
}
