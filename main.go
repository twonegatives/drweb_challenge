package main

import (
	"crypto/sha256"
	"fmt"
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

func loadFile(hashstring string) (*File, error) {
	content, err := ioutil.ReadFile(hashstring)
	return &File{Body: content}, err
}

func deleteFile(hashstring string) error {
	return os.Remove(hashstring)
}

func main() {
	input := []byte("This is an example file")
	file := File{Body: input}
	hash, err := file.Save()

	if err != nil {
		panic(fmt.Sprintf("could not save the file: %s", err))
	}

	loadedBack, err := loadFile(hash)

	if err != nil {
		panic(fmt.Sprintf("could not load the file: %s", err))
	}

	fmt.Println("saved and loaded back successfully")
	fmt.Println(string(loadedBack.Body[:]))

	err = deleteFile(hash)

	if err != nil {
		panic(fmt.Sprintf("could not delete the file: %s", err))
	}

	fmt.Println("deleted aswell")
}
