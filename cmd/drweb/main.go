package main

import (
	"fmt"

	"github.com/twonegatives/drweb_challenge/pkg/drweb"
	"github.com/twonegatives/drweb_challenge/pkg/encoders"
	"github.com/twonegatives/drweb_challenge/pkg/filesavehooks"
	"github.com/twonegatives/drweb_challenge/pkg/storages"
)

func main() {
	input := []byte("This is an example file")
	encoder := encoders.SHA256Encoder{}
	hooks := filesavehooks.PrintlnHook{}

	storage := storages.FileSystemStorage{
		BasePath: ".",
		FileMode: 0600,
	}

	encoded := string(encoder.Encode(input)[:])

	file := drweb.File{
		Body:        input,
		Storage:     &storage,
		HooksOnSave: &hooks,
		Encoder:     &encoder,
	}

	_, err := file.Save()

	if err != nil {
		panic(fmt.Sprintf("could not save the file: %s", err))
	}

	loadedBack, err := storage.Load(encoded)

	if err != nil {
		panic(fmt.Sprintf("could not load the file: %s", err))
	}

	fmt.Println("saved and loaded back successfully")
	fmt.Println(string(loadedBack.Body))

	err = storage.Delete(encoded)

	if err != nil {
		panic(fmt.Sprintf("could not delete the file: %s", err))
	}

	fmt.Println("deleted aswell")
}
