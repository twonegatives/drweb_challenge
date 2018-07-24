package storages

import (
	"io"
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
	path := s.filepath(file.Filename)
	if err := os.MkdirAll(filepath.Dir(path), s.FileMode); err != nil {
		return path, err
	}

	output, err := os.Create(path)
	defer output.Close()
	if err != nil {
		return path, err
	}

	_, err = io.Copy(output, file.Body)

	return path, err
}

func (s *FileSystemStorage) Load(filename string) (*drweb.File, error) {
	reader, err := os.Open(filename)

	if err != nil {
		return nil, errors.Wrap(err, "failed to load file")
	}

	return &drweb.File{Body: reader}, nil
}

func (s *FileSystemStorage) Delete(filename string) error {
	return os.Remove(s.filepath(filename))
}
