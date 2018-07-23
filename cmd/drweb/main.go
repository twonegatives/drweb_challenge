package main

import (
	"fmt"

	filepkg "github.com/twonegatives/drweb_challenge/pkg/file"
	uploaders "github.com/twonegatives/drweb_challenge/pkg/localuploader"
	encoders "github.com/twonegatives/drweb_challenge/pkg/sha256encoder"
)

func main() {
	input := []byte("This is an example file")

	file := filepkg.File{
		Body:    input,
		Encoder: &encoders.SHA256Encoder{},
		Uploader: &uploaders.LocalUploader{
			BasePath: ".",
			FileMode: 0600,
		},
	}

	filepath, err := file.Save()

	if err != nil {
		panic(fmt.Sprintf("could not save the file: %s", err))
	}

	loadedBack, err := filepkg.LoadFile(filepath)

	if err != nil {
		panic(fmt.Sprintf("could not load the file: %s", err))
	}

	fmt.Println("saved and loaded back successfully")
	fmt.Println(string(loadedBack.Body[:]))

	err = filepkg.DeleteFile(filepath)

	if err != nil {
		panic(fmt.Sprintf("could not delete the file: %s", err))
	}

	fmt.Println("deleted aswell")
}
