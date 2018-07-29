package storages

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/twonegatives/drweb_challenge/pkg/drweb"
)

type FileSystemStorage struct {
	FileMode          os.FileMode
	FilePathGenerator drweb.FilePathGenerator
}

func (s *FileSystemStorage) filepath(filename string) (string, error) {
	return s.FilePathGenerator.Generate(filename)
}

func (s *FileSystemStorage) Save(file *drweb.File) (string, error) {
	var filename string
	var path string
	var err error

	if file.NameGenerator == nil {
		return filename, errors.New("failed to save file without name generator")
	}

	// NOTE: we use temp file here as it solves filename uniqueness for us
	// although it requires us to chmod and rename it later
	tmpfile, err := ioutil.TempFile(os.TempDir(), "prefix")
	if err != nil {
		return filename, errors.Wrap(err, "failed to create file")
	}

	defer tmpfile.Close()

	if err = tmpfile.Chmod(s.FileMode); err != nil {
		return filename, errors.Wrap(err, "failed to set requested file mode")
	}

	filenameReader := io.TeeReader(file.Body, tmpfile)
	filename, err = file.NameGenerator.Generate(filenameReader, file.Extension)

	if err != nil {
		return filename, errors.Wrap(err, "failed to generate filename")
	}

	if path, err = s.filepath(filename); err != nil {
		return filename, errors.Wrap(err, "failed to generate filepath")
	}

	if err = os.MkdirAll(filepath.Dir(path), s.FileMode); err != nil {
		return filename, errors.Wrap(err, "failed to create nested folders")
	}

	err = os.Rename(tmpfile.Name(), path)

	return filename, errors.Wrap(err, "failed to write to file")
}

func (s *FileSystemStorage) Load(filename string) (*drweb.File, error) {
	var file *drweb.File
	var reader io.ReadCloser
	var path string
	var err error

	if path, err = s.filepath(filename); err != nil {
		return nil, errors.Wrap(err, "failed to generate filepath")
	}

	if reader, err = os.Open(path); err != nil {
		return nil, errors.Wrap(err, "failed to open file")
	}

	file = &drweb.File{
		Body:      reader,
		Extension: filepath.Ext(filename),
	}

	return file, nil
}

func (s *FileSystemStorage) Delete(filename string) error {
	var path string
	var err error

	if path, err = s.filepath(filename); err != nil {
		return errors.Wrap(err, "failed to generate filepath")
	}

	return os.Remove(path)
}
