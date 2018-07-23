package file

import (
	"io/ioutil"
	"os"

	"github.com/twonegatives/drweb_challenge/pkg/drweb"
)

type File struct {
	Body    []byte
	Encoder drweb.FileEncoder
}

func (f *File) toHash() string {
	hashbytes := f.Encoder.Encode(f.Body)
	return string(hashbytes[:])
}

func (f *File) Save() (string, error) {
	hashstring := f.toHash()
	err := ioutil.WriteFile(hashstring, f.Body, 0600)
	return hashstring, err
}

func LoadFile(hashstring string) (*File, error) {
	content, err := ioutil.ReadFile(hashstring)
	return &File{Body: content}, err
}

func DeleteFile(hashstring string) error {
	return os.Remove(hashstring)
}
