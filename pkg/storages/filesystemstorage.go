package storages

import (
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/twonegatives/drweb_challenge/pkg/drweb"
)

type FileSystemStorage struct {
	FileMode          os.FileMode
	FilePathGenerator drweb.FilePathGenerator
}

func (s *FileSystemStorage) Filepath(filename string) (string, error) {
	return s.FilePathGenerator.Generate(filename)
}

func (s *FileSystemStorage) Save(file *drweb.File) (string, error) {
	var path string
	var err error

	if path, err = s.Filepath(file.Filename); err != nil {
		return path, errors.Wrap(err, "failed to generate filepath")
	}

	if err = os.MkdirAll(filepath.Dir(path), s.FileMode); err != nil {
		return path, errors.Wrap(err, "failed to create nested folders")
	}

	output, err := os.Create(path)
	defer output.Close()
	if err != nil {
		return path, errors.Wrap(err, "failed to create file")
	}

	_, err = io.Copy(output, file.Body)

	return path, errors.Wrap(err, "failed to write to file")
}

func (s *FileSystemStorage) Load(filename string) (*drweb.File, error) {
	reader, err := os.Open(filename)

	if err != nil {
		return nil, errors.Wrap(err, "failed to load file")
	}

	return &drweb.File{Body: reader}, nil
}

func (s *FileSystemStorage) Delete(filename string) error {
	var path string
	var err error

	if path, err = s.Filepath(filename); err != nil {
		return errors.Wrap(err, "failed to generate filepath")
	}

	return os.Remove(path)
}
