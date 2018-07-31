package testutils

import (
	"bytes"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"
)

func CreateFile(path string, contents []byte, fileMode os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), fileMode); err != nil {
		return err
	}

	return ioutil.WriteFile(path, contents, fileMode)
}

func FileToFormData(filename string, contents []byte, paramName string) (*bytes.Buffer, string, error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filename)
	if err != nil {
		return nil, "", err
	}

	part.Write(contents)
	err = writer.Close()
	if err != nil {
		return nil, "", err
	}

	return body, writer.Boundary(), nil
}
