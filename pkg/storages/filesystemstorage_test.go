package storages_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/twonegatives/drweb_challenge/pkg/drweb"
	"github.com/twonegatives/drweb_challenge/pkg/mocks"
	"github.com/twonegatives/drweb_challenge/pkg/storages"
)

func TestDeleteUnexistantFile(t *testing.T) {
	filename := "delete_me1"
	path := path.Join("../../tmp", filename)

	mockCtrl := gomock.NewController(t)
	pathgen := mocks.NewMockFilePathGenerator(mockCtrl)
	pathgen.EXPECT().Generate(filename).Return(path, nil)
	storage := storages.FileSystemStorage{
		FilePathGenerator: pathgen,
	}

	err := storage.Delete(filename)
	assert.NotNil(t, err)
	assert.Equal(t, true, os.IsNotExist(err))
}

func TestDeleteBrokenFilepath(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	pathgen := mocks.NewMockFilePathGenerator(mockCtrl)
	pathgen.EXPECT().Generate("filename").Return("", errors.New("generation error"))
	storage := storages.FileSystemStorage{
		FilePathGenerator: pathgen,
	}

	err := storage.Delete("filename")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "failed to generate filepath")
}

func TestDeleteSuccess(t *testing.T) {
	filename := "delete_me1"
	path := path.Join("../../tmp", filename)
	err := ioutil.WriteFile(path, []byte("contents"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	mockCtrl := gomock.NewController(t)
	pathgen := mocks.NewMockFilePathGenerator(mockCtrl)
	pathgen.EXPECT().Generate(filename).Return(path, nil)
	storage := storages.FileSystemStorage{
		FilePathGenerator: pathgen,
	}

	err = storage.Delete(filename)
	assert.Nil(t, err)

	_, err = os.Lstat(path)
	assert.NotNil(t, err)
	assert.Equal(t, true, os.IsNotExist(err))
}

func TestLoadBrokenFilepath(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	pathgen := mocks.NewMockFilePathGenerator(mockCtrl)
	pathgen.EXPECT().Generate("filename").Return("", errors.New("generation error"))
	storage := storages.FileSystemStorage{
		FilePathGenerator: pathgen,
	}

	_, err := storage.Load("filename")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "failed to generate filepath")
}

func TestLoadSuccess(t *testing.T) {
	filename := "load_me1"
	path := path.Join("../../tmp", filename)
	err := ioutil.WriteFile(path, []byte("contents"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(path)

	mockCtrl := gomock.NewController(t)
	pathgen := mocks.NewMockFilePathGenerator(mockCtrl)
	pathgen.EXPECT().Generate(filename).Return(path, nil)
	storage := storages.FileSystemStorage{
		FilePathGenerator: pathgen,
	}

	file, err := storage.Load(filename)
	assert.Nil(t, err)

	defer file.Body.Close()
	contents, err := ioutil.ReadAll(file.Body)
	assert.Nil(t, err)
	assert.Equal(t, contents, []byte("contents"))
}

func TestLoadUnexistantFile(t *testing.T) {
	filename := "load_unexistant"
	path := path.Join("../../tmp", filename)

	mockCtrl := gomock.NewController(t)
	pathgen := mocks.NewMockFilePathGenerator(mockCtrl)
	pathgen.EXPECT().Generate(filename).Return(path, nil)
	storage := storages.FileSystemStorage{
		FilePathGenerator: pathgen,
	}

	_, err := storage.Load(filename)
	assert.NotNil(t, err)
	assert.Equal(t, true, os.IsNotExist(errors.Cause(err)))
}

func TestSaveBlankNameGenerator(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	pathgen := mocks.NewMockFilePathGenerator(mockCtrl)
	pathgen.EXPECT().Generate(gomock.Any()).Return("some/path", nil)

	storage := storages.FileSystemStorage{
		FilePathGenerator: pathgen,
	}

	_, err := storage.Save(&drweb.File{})
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "failed to save file without name generator")
}

func TestSaveBrokenFilename(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	namegen := mocks.NewMockFileNameGenerator(mockCtrl)
	namegen.EXPECT().Generate(gomock.Any(), gomock.Any()).Return("", errors.New("name generation error"))

	pathgen := mocks.NewMockFilePathGenerator(mockCtrl)
	pathgen.EXPECT().Generate(gomock.Any()).Return("some/path", nil)

	storage := storages.FileSystemStorage{
		FilePathGenerator: pathgen,
	}

	_, err := storage.Save(&drweb.File{NameGenerator: namegen})
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "failed to generate filename")
}

func TestSaveBrokenFilepath(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	namegen := mocks.NewMockFileNameGenerator(mockCtrl)
	namegen.EXPECT().Generate(gomock.Any(), gomock.Any()).Return("encrypted", nil)

	pathgen := mocks.NewMockFilePathGenerator(mockCtrl)
	pathgen.EXPECT().Generate("encrypted").Return("", errors.New("path generation error"))

	storage := storages.FileSystemStorage{
		FilePathGenerator: pathgen,
	}

	_, err := storage.Save(&drweb.File{NameGenerator: namegen})
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "failed to generate filepath")
}

func TestSaveSuccess(t *testing.T) {
	filename := "encrypted1"
	path := path.Join("../../tmp", filename)
	contents := []byte("File contents")
	// NOTE: we use this instead of mock version as long as we use TeeReader in storage#Save
	// this forces us to consume reader in order to fill file with expected contents
	namegen := &staticFileNameGenerator{Name: filename}

	mockCtrl := gomock.NewController(t)
	pathgen := mocks.NewMockFilePathGenerator(mockCtrl)
	pathgen.EXPECT().Generate(filename).Return(path, nil)

	storage := storages.FileSystemStorage{
		FileMode:          0700,
		FilePathGenerator: pathgen,
	}

	file := drweb.File{
		Body:          ioutil.NopCloser(bytes.NewReader(contents)),
		MimeType:      "image/png",
		NameGenerator: namegen,
	}

	savedFileName, err := storage.Save(&file)
	defer os.Remove(path)

	if err != nil {
		t.Fatal(err)
	}

	bytes, err := ioutil.ReadFile(path)
	assert.Nil(t, err)
	assert.Equal(t, savedFileName, filename)
	assert.Equal(t, contents, bytes)
}

type staticFileNameGenerator struct {
	Name string
}

func (g *staticFileNameGenerator) Generate(input io.Reader, mimeType string) (string, error) {
	ioutil.ReadAll(input)
	return g.Name, nil
}
