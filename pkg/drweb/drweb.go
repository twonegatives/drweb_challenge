package drweb

import (
	"io"
	"net/textproto"

	"github.com/pkg/errors"
)

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
	Body     io.Reader
	Filename string `json:"filename"`
	Storage  Storage
	MimeType textproto.MIMEHeader
}

type FileNameGenerator interface {
	Generate(file *File) (string, error)
}

type FilePathGenerator interface {
	Generate(filename string) (string, error)
}

func NewFile(body io.Reader, mimetype textproto.MIMEHeader, storage Storage, nameGenerator FileNameGenerator) (*File, error) {
	file := File{
		Body:     body,
		Storage:  storage,
		MimeType: mimetype,
	}

	filename, err := nameGenerator.Generate(&file)
	if err != nil {
		return &file, errors.Wrap(err, "failed to generate filename")
	}

	file.Filename = filename
	return &file, nil
}
