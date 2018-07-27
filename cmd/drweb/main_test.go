package main_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	main "github.com/twonegatives/drweb_challenge/cmd/drweb"
	"github.com/twonegatives/drweb_challenge/pkg/mocks"
)

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
