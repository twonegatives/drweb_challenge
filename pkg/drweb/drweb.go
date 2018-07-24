package drweb

import (
	"fmt"
	"io"
	"io/ioutil"

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

func NewFile(body io.Reader, storage Storage, hooks FileSaveHooks, encoder Encoder) *File {
	file := File{
		Body:        body,
		Storage:     storage,
		HooksOnSave: hooks,
	}

	data, _ := ioutil.ReadAll(file.Body)
	file.Filename = fmt.Sprintf("%x", encoder.Encode(data))

	return &file
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
