package drweb

import (
	"io"
)

//go:generate mockgen -source=drweb.go -destination ../mocks/mock_drweb.go -package mocks

type Encoder interface {
	Encode(input []byte) []byte
}

type Callback interface {
	Invoke(args ...interface{})
}

type Storage interface {
	Save(f *File) (string, error)
	Load(filename string) (*File, error)
	Delete(filename string) error
	Filepath(filename string) (string, error)
}

type File struct {
	Body          io.Reader
	MimeType      string
	NameGenerator FileNameGenerator
}

type FileNameGenerator interface {
	Generate(input io.Reader, mimeType string) (string, error)
}

type FilePathGenerator interface {
	Generate(filename string) (string, error)
}
