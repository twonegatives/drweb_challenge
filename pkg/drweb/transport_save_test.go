package drweb_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/twonegatives/drweb_challenge/pkg/drweb"
	"github.com/twonegatives/drweb_challenge/pkg/mocks"
	"github.com/twonegatives/drweb_challenge/pkg/testutils"
)

func TestSaveFailure(t *testing.T) {
	t.Run("no file form", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		storage := mocks.NewMockStorage(mockCtrl)
		filenamegenerator := mocks.NewMockFileNameGenerator(mockCtrl)

		req, err := http.NewRequest("POST", "/files", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/files", drweb.CreateFileHandler(storage, filenamegenerator))
		router.ServeHTTP(rr, req)

		var response map[string]string
		json.Unmarshal(rr.Body.Bytes(), &response)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
		assert.NotNil(t, response["error"])
	})

	t.Run("storage failure", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		storage := mocks.NewMockStorage(mockCtrl)
		storage.EXPECT().Save(gomock.Any()).Return("", errors.New("storage is corrupted"))
		filenamegenerator := mocks.NewMockFileNameGenerator(mockCtrl)

		multipartBody, multipartBoundary, err := testutils.FileToMultipartForm("original_filename", []byte("Byte file contents"), "file")
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest("POST", "/files", multipartBody)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Content-Type", fmt.Sprintf("multipart/form-data; boundary=\"%s\"", multipartBoundary))

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/files", drweb.CreateFileHandler(storage, filenamegenerator))
		router.ServeHTTP(rr, req)

		var response map[string]string
		json.Unmarshal(rr.Body.Bytes(), &response)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
		assert.NotNil(t, response["error"])
	})
}

func TestSaveFileHandlerSuccess(t *testing.T) {
	filename := "filename_to_user"
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	storage := mocks.NewMockStorage(mockCtrl)
	storage.EXPECT().Save(gomock.Any()).Return(filename, nil)
	filenamegenerator := mocks.NewMockFileNameGenerator(mockCtrl)

	multipartBody, multipartBoundary, err := testutils.FileToMultipartForm("original_filename", []byte("Byte file contents"), "file")
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/files", multipartBody)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", fmt.Sprintf("multipart/form-data; boundary=\"%s\"", multipartBoundary))

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/files", drweb.CreateFileHandler(storage, filenamegenerator))
	router.ServeHTTP(rr, req)

	var response map[string]string
	json.Unmarshal(rr.Body.Bytes(), &response)

	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	assert.Equal(t, map[string]string{"hashstring": filename}, response)
}
