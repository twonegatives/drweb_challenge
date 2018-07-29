package storages_test

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/twonegatives/drweb_challenge/pkg/testutils"
)

func TestDeleteFailure(t *testing.T) {
	t.Run("unexistant file", func(t *testing.T) {
		filename := "delete_me1"
		path := path.Join("../../tmp", filename)

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		storage := testutils.GenerateStorage(filename, path, nil, mockCtrl)

		err := storage.Delete(filename)
		assert.NotNil(t, err)
		assert.Equal(t, true, os.IsNotExist(err))
	})

	t.Run("broken filepath", func(t *testing.T) {
		filename := "delete_me1"

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		storage := testutils.GenerateStorage(filename, "", errors.New("generation error"), mockCtrl)

		err := storage.Delete(filename)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "failed to generate filepath")
	})
}

func TestDeleteSuccess(t *testing.T) {
	filename := "delete_me1"
	path := path.Join("../../tmp", filename)
	err := ioutil.WriteFile(path, []byte("contents"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	storage := testutils.GenerateStorage(filename, path, nil, mockCtrl)

	err = storage.Delete(filename)
	assert.Nil(t, err)

	_, err = os.Lstat(path)
	assert.NotNil(t, err)
	assert.Equal(t, true, os.IsNotExist(err))
}
