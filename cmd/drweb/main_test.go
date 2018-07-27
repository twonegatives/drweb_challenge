package main_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	main "github.com/twonegatives/drweb_challenge/cmd/drweb"
	"github.com/twonegatives/drweb_challenge/pkg/drweb"
	"github.com/twonegatives/drweb_challenge/pkg/mocks"
)

func TestRetrieveFileHandlerSuccess(t *testing.T) {
	filename := "some_saved_file"
	contents := []byte("Loaded file")
	file := drweb.File{Body: ioutil.NopCloser(bytes.NewReader(contents))}

	mockCtrl := gomock.NewController(t)
	storage := mocks.NewMockStorage(mockCtrl)
	storage.EXPECT().Load(filename).Return(&file, nil)

	req, err := http.NewRequest("GET", "/files/some_saved_file", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/files/{hashstring}", main.RetrieveFileHandler(storage))
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, string(contents), rr.Body.String())
}

func TestRetrieveFileHandlerNotFound(t *testing.T) {
	filename := "unexistant_file"
	mockCtrl := gomock.NewController(t)
	storage := mocks.NewMockStorage(mockCtrl)
	storage.EXPECT().Load(filename).Return(nil, errors.Wrap(os.ErrNotExist, "some description"))

	req, err := http.NewRequest("GET", "/files/unexistant_file", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/files/{hashstring}", main.RetrieveFileHandler(storage))
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestDeleteFileHandlerSuccess(t *testing.T) {
	filename := "delete_me_test_main"
	mockCtrl := gomock.NewController(t)
	storage := mocks.NewMockStorage(mockCtrl)
	storage.EXPECT().Delete(filename).Return(nil)

	req, err := http.NewRequest("DELETE", "/files/delete_me_test_main", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/files/{hashstring}", main.DeleteFileHandler(storage))
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Empty(t, rr.Body.String())
}

func TestDeleteFileHandlerNotFound(t *testing.T) {
	filename := "delete_me_test_main"
	mockCtrl := gomock.NewController(t)
	storage := mocks.NewMockStorage(mockCtrl)
	storage.EXPECT().Delete(filename).Return(errors.Wrap(os.ErrNotExist, "some description"))

	req, err := http.NewRequest("DELETE", "/files/delete_me_test_main", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/files/{hashstring}", main.DeleteFileHandler(storage))
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}
