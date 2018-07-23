package main

import (
	"fmt"

	filepkg "github.com/twonegatives/drweb_challenge/pkg/file"
)

func main() {
	input := []byte("This is an example file")
	file := filepkg.File{Body: input}
	hash, err := file.Save()

	if err != nil {
		panic(fmt.Sprintf("could not save the file: %s", err))
	}

	loadedBack, err := filepkg.LoadFile(hash)

	if err != nil {
		panic(fmt.Sprintf("could not load the file: %s", err))
	}

	fmt.Println("saved and loaded back successfully")
	fmt.Println(string(loadedBack.Body[:]))

	err = filepkg.DeleteFile(hash)

	if err != nil {
		panic(fmt.Sprintf("could not delete the file: %s", err))
	}

	fmt.Println("deleted aswell")
}
