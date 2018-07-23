package storages

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/twonegatives/drweb_challenge/pkg/drweb"
)

type FileSystemStorage struct {
	BasePath string
	FileMode os.FileMode
}

func (s *FileSystemStorage) filepath(filename string) string {
	return path.Join(s.BasePath, filename[0:2], filename[2:4], filename)
}

func (s *FileSystemStorage) Save(file *drweb.File) (string, error) {
	path := s.filepath(file.GetFilename())
	if err := os.MkdirAll(filepath.Dir(path), s.FileMode); err != nil {
		return path, err
	}
	err := ioutil.WriteFile(path, file.Body, s.FileMode)
	return path, err
}

func (s *FileSystemStorage) Load(filename string) (*drweb.File, error) {
	contents, err := ioutil.ReadFile(s.filepath(filename))

	if err != nil {
		return nil, errors.Wrap(err, "failed to load file")
	}

	return &drweb.File{Body: contents}, nil
}

func (s *FileSystemStorage) Delete(filename string) error {
	return os.Remove(s.filepath(filename))
}
