package filesavehooks

import (
	"fmt"

	"github.com/twonegatives/drweb_challenge/pkg/drweb"
)

type PrintlnHook struct {
}

func (h *PrintlnHook) Before(file *drweb.File) error {
	fmt.Println("before file save")
	return nil
}

func (h *PrintlnHook) After(file *drweb.File, filename string, filepath string) error {
	fmt.Println("after file save")
	return nil
}
