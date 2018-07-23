package file

import (
	"github.com/twonegatives/drweb_challenge/pkg/drweb"
)

type File struct {
	Body    []byte
	Encoder drweb.FileEncoder
	Storage drweb.Storage
}

func (f *File) encode() string {
	hashbytes := f.Encoder.Encode(f.Body)
	return string(hashbytes[:])
}

func (f *File) Save() (string, error) {
	filename := f.encode()
	return f.Storage.Save(f.Body, filename)
}
