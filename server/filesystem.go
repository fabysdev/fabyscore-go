package server

import (
	"net/http"
	"strings"
)

// FileSystem is a custom http.FileSystem implementation to avoid serving directories
type FileSystem struct {
	fileSystem http.FileSystem
}

// Open returns the file or an error
func (fs FileSystem) Open(path string) (http.File, error) {
	f, err := fs.fileSystem.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if s.IsDir() {
		// serve the directory if it contains an index.html
		index := strings.TrimSuffix(path, "/") + "/index.html"

		if _, err := fs.fileSystem.Open(index); err != nil {
			return nil, err
		}
	}

	return f, nil
}
