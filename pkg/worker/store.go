package worker

import (
	"net/http"
	"os"
	"path"
)

type dirStore struct {
	dirPath string
	fs http.FileSystem
}

func NewDirStore(dirPath string) (*dirStore, error) {
	return &dirStore{
		dirPath: dirPath,
	}, nil
}

func (s *dirStore) Get(sessionID string) (*Job, error) {
	s.fs.Open()
	dirpath := path.Join(s.dirPath, sessionID)
	dir, err := os.Open()
}


