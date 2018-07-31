package drweb

import (
	"io"
)

//go:generate mockgen -source=drweb.go -destination ../mocks/mock_drweb.go -package mocks

type Callback interface {
	Invoke(args ...interface{})
}

type Storage interface {
	Save(f *FileCreateRequest) (string, error)
	Load(filename string) (*File, error)
	Delete(filename string) error
}

type FileCreateRequest struct {
	Body          io.ReadCloser
	NameGenerator FileNameGenerator
}

func (f *FileCreateRequest) Close() error {
	return f.Body.Close()
}

type File struct {
	Body io.ReadCloser
	Size int64
}

func (f *File) Close() error {
	return f.Body.Close()
}

type FileNameGenerator interface {
	Generate(input io.Reader) (string, error)
}

type FilePathGenerator interface {
	Generate(filename string) (string, error)
}
