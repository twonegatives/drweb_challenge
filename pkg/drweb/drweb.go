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

type Storage interface {
	Save(f *File) (string, error)
	Load(filename string) (*File, error)
	Delete(filename string) error
}

type FileSaveHooks interface {
	Before(file *File, args ...interface{}) error
	After(file *File, args ...interface{}) error
}

type File struct {
	Body        io.Reader
	Storage     Storage
	HooksOnSave FileSaveHooks
	Filename    string `json:"filename"`
}

func NewFile(body io.Reader, storage Storage, hooks FileSaveHooks, encoder Encoder) (*File, error) {
	file := File{
		Body:        body,
		Storage:     storage,
		HooksOnSave: hooks,
	}

	leadingChars := make([]byte, 50)
	if _, err := file.Body.Read(leadingChars); err != nil {
		return &file, errors.Wrap(err, "failed to build a new file object")
	}

	file.Filename = fmt.Sprintf("%x-%d", encoder.Encode(leadingChars), time.Now().UnixNano())

	return &file, nil
}

func (f *File) Save() (string, error) {
	if err := f.HooksOnSave.Before(f); err != nil {
		return "", errors.Wrap(err, "beforeSave hook failed")
	}

	filepath, err := f.Storage.Save(f)

	if err != nil {
		return "", errors.Wrap(err, "failed to save file")
	}

	if err := f.HooksOnSave.After(f); err != nil {
		return "", errors.Wrap(err, "afterSave hook failed")
	}

	return filepath, nil
}
