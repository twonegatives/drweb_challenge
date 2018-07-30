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
	Filename string
	Path     string
	Size     int64
}

func TestLoadSuccess(t *testing.T) {
	objects := map[string]loadSuccessCase{
		"text file": {
			Filename: "alice",
			Path:     "../testdata/alice.txt",
			Size:     4094,
		},
		"image file": {
			Filename: "gopher",
			Path:     "../testdata/gopher.jpg",
			Size:     6707,
		},
	}

	for testName, testObject := range objects {
		t.Run(testName, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			storage := testutils.GenerateStorage(testObject.Filename, testObject.Path, nil, mockCtrl)

			file, err := storage.Load(testObject.Filename)
			assert.Nil(t, err)
			assert.Equal(t, testObject.Size, file.Size)

			_, err = ioutil.ReadAll(file.Body)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
