package filesystemstorage

import (
	"io/ioutil"
	"os"
	"path"
)

type FileSystemStorage struct {
	BasePath string
	FileMode os.FileMode
}

func (s *FileSystemStorage) filepath(filename string) string {
	return path.Join(s.BasePath, filename)
}

func (s *FileSystemStorage) Save(contents []byte, filename string) (string, error) {
	path := s.filepath(filename)
	err := ioutil.WriteFile(path, contents, s.FileMode)
	return path, err
}

func (s *FileSystemStorage) Load(filename string) ([]byte, error) {
	return ioutil.ReadFile(s.filepath(filename))
}

func (s *FileSystemStorage) Delete(filename string) error {
	return os.Remove(s.filepath(filename))
}
