package localuploader

import (
	"io/ioutil"
	"os"
	"path"
)

type LocalUploader struct {
	BasePath string
	FileMode os.FileMode
}

func (u *LocalUploader) Upload(contents []byte, filename string) (string, error) {
	filepath := path.Join(u.BasePath, filename)
	err := ioutil.WriteFile(filepath, contents, u.FileMode)
	return filepath, err
}
