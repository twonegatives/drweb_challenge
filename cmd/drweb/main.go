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

	file := drweb.File{
		Body:        input,
		Encoder:     &encoder,
		Storage:     &storage,
		HooksOnSave: &hooks,
	}

	_, err := file.Save()

	if err != nil {
		panic(fmt.Sprintf("could not save the file: %s", err))
	}

	loadedBack, err := storage.Load(string(encoder.Encode(input)))

	if err != nil {
		panic(fmt.Sprintf("could not load the file: %s", err))
	}

	fmt.Println("saved and loaded back successfully")
	fmt.Println(string(loadedBack))

	err = storage.Delete(string(encoder.Encode(input)))

	if err != nil {
		panic(fmt.Sprintf("could not delete the file: %s", err))
	}

	fmt.Println("deleted aswell")
}
