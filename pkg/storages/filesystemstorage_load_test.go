package storages_test

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/twonegatives/drweb_challenge/pkg/testutils"
)

func TestLoadFailure(t *testing.T) {
	t.Run("broken filepath", func(t *testing.T) {
		filename := "filename"
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		storage := testutils.GenerateStorage(filename, "", errors.New("generation error"), mockCtrl)

		_, err := storage.Load(filename)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "failed to generate filepath")
	})

	t.Run("unexistant file", func(t *testing.T) {
		filename := "load_unexistant"
		path := path.Join("../../tmp", filename)

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		storage := testutils.GenerateStorage(filename, path, nil, mockCtrl)

		_, err := storage.Load(filename)
		assert.NotNil(t, err)
		assert.Equal(t, true, os.IsNotExist(errors.Cause(err)))
	})
}

type loadSuccessCase struct {
	Filename          string
	Filemode          os.FileMode
	Contents          []byte
	ExpectedExtension string
}

func TestLoadSuccess(t *testing.T) {
	objects := map[string]loadSuccessCase{
		"with extension": {
			Filename:          "load_me1.txt",
			Filemode:          0644,
			ExpectedExtension: ".txt",
			Contents:          []byte("some contents"),
		},
		"without extension": {
			Filename:          "load_me1",
			Filemode:          0644,
			ExpectedExtension: "",
			Contents:          []byte("contents without ext"),
		},
	}

	for testName, testObject := range objects {
		t.Run(testName, func(t *testing.T) {
			path := path.Join("../../tmp", testObject.Filename)
			if err := testutils.CreateFile(path, testObject.Contents, testObject.Filemode); err != nil {
				t.Fatal(err)
			}

			defer os.Remove(path)

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			storage := testutils.GenerateStorage(testObject.Filename, path, nil, mockCtrl)

			file, err := storage.Load(testObject.Filename)
			defer file.Body.Close()
			assert.Nil(t, err)
			assert.NotNil(t, file.Extension)
			assert.Equal(t, testObject.ExpectedExtension, file.Extension)

			contents, err := ioutil.ReadAll(file.Body)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, testObject.Contents, contents)
		})
	}
}
