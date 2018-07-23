package file

import (
	"crypto/sha256"
	"io/ioutil"
	"os"
)

type File struct {
	Body []byte
}

func (f *File) toHash() string {
	hashbytes := sha256.Sum256(f.Body)
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
