package drweb

import (
	"fmt"
	"io"
	"time"

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
}

type File struct {
	Body     io.Reader
	Filename string `json:"filename"`
	Storage  Storage
}

func NewFile(body io.Reader, storage Storage, encoder Encoder) (*File, error) {
	file := File{
		Body:    body,
		Storage: storage,
	}

	leadingChars := make([]byte, 50)
	if _, err := file.Body.Read(leadingChars); err != nil {
		return &file, errors.Wrap(err, "failed to build a new file object")
	}

	file.Filename = fmt.Sprintf("%x-%d", encoder.Encode(leadingChars), time.Now().UnixNano())

	return &file, nil
}
