package testutils

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func CreateFile(path string, contents []byte, fileMode os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), fileMode); err != nil {
		return err
	}

	return ioutil.WriteFile(path, contents, fileMode)
}
