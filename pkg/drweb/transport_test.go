package drweb_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/twonegatives/drweb_challenge/pkg/drweb"
	"github.com/twonegatives/drweb_challenge/pkg/mocks"
	"github.com/twonegatives/drweb_challenge/pkg/testutils"
)

func TestSaveFileHandlerWithoutFileForm(t *testing.T) {
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
}

func TestSaveFileHandlerWithStorageFail(t *testing.T) {
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

func TestRetrieveFileHandlerUnknownTypeSuccess(t *testing.T) {
	filename := "some_saved_file"
	contents := []byte("Loaded file")
	file := drweb.File{Body: ioutil.NopCloser(bytes.NewReader(contents))}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	storage := mocks.NewMockStorage(mockCtrl)
	storage.EXPECT().Load(filename).Return(&file, nil)

	req, err := http.NewRequest("GET", "/files/some_saved_file", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/files/{hashstring}", drweb.RetrieveFileHandler(storage))
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, string(contents), rr.Body.String())
	assert.Equal(t, "application/octet-stream", rr.Header().Get("Content-Type"))
}

func TestRetrieveFileHandlerWithExtensionSuccess(t *testing.T) {
	filename := "some_saved_file.jpg"
	contents := []byte("Loaded file")
	file := drweb.File{Body: ioutil.NopCloser(bytes.NewReader(contents)), Extension: ".jpg"}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	storage := mocks.NewMockStorage(mockCtrl)
	storage.EXPECT().Load(filename).Return(&file, nil)

	req, err := http.NewRequest("GET", "/files/some_saved_file.jpg", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/files/{hashstring}", drweb.RetrieveFileHandler(storage))
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, string(contents), rr.Body.String())
	assert.Equal(t, "image/jpeg", rr.Header().Get("Content-Type"))
}

func TestRetrieveFileHandlerNotFound(t *testing.T) {
	filename := "unexistant_file"
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	storage := mocks.NewMockStorage(mockCtrl)
	storage.EXPECT().Load(filename).Return(nil, errors.Wrap(os.ErrNotExist, "some description"))

	req, err := http.NewRequest("GET", "/files/unexistant_file", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/files/{hashstring}", drweb.RetrieveFileHandler(storage))
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
}

func TestRetrieveFileHandlerLoadFailed(t *testing.T) {
	filename := "unlucky_file"
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	storage := mocks.NewMockStorage(mockCtrl)
	storage.EXPECT().Load(filename).Return(nil, errors.Wrap(errors.New("some error"), "some description"))

	req, err := http.NewRequest("GET", "/files/unlucky_file", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/files/{hashstring}", drweb.RetrieveFileHandler(storage))
	router.ServeHTTP(rr, req)

	var response map[string]string
	json.Unmarshal(rr.Body.Bytes(), &response)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	assert.Contains(t, response["error"], "some error")
}

func TestDeleteFileHandlerSuccess(t *testing.T) {
	filename := "delete_me_test_main"
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	storage := mocks.NewMockStorage(mockCtrl)
	storage.EXPECT().Delete(filename).Return(nil)

	req, err := http.NewRequest("DELETE", "/files/delete_me_test_main", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/files/{hashstring}", drweb.DeleteFileHandler(storage))
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	assert.Empty(t, rr.Body.String())
}

func TestDeleteFileHandlerFail(t *testing.T) {
	filename := "delete_me_test_main"
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	storage := mocks.NewMockStorage(mockCtrl)
	storage.EXPECT().Delete(filename).Return(errors.Wrap(errors.New("some error"), "some description"))

	req, err := http.NewRequest("DELETE", "/files/delete_me_test_main", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/files/{hashstring}", drweb.DeleteFileHandler(storage))
	router.ServeHTTP(rr, req)

	var response map[string]string
	json.Unmarshal(rr.Body.Bytes(), &response)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	assert.Contains(t, response["error"], "some error")
}

func TestDeleteFileHandlerNotFound(t *testing.T) {
	filename := "delete_me_test_main"
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	storage := mocks.NewMockStorage(mockCtrl)
	storage.EXPECT().Delete(filename).Return(errors.Wrap(os.ErrNotExist, "some description"))

	req, err := http.NewRequest("DELETE", "/files/delete_me_test_main", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/files/{hashstring}", drweb.DeleteFileHandler(storage))
	router.ServeHTTP(rr, req)

	var response map[string]string
	json.Unmarshal(rr.Body.Bytes(), &response)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	assert.NotNil(t, response["error"])
}

func TestWithCallbacks(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	before := mocks.NewMockCallback(mockCtrl)
	after := mocks.NewMockCallback(mockCtrl)

	rr := httptest.NewRecorder()
	handler := func(http.ResponseWriter, *http.Request) {}

	gomock.InOrder(
		before.EXPECT().Invoke(gomock.Any(), gomock.Any()),
		after.EXPECT().Invoke(gomock.Any(), gomock.Any()),
	)

	function := drweb.WithCallbacks(handler, before, after)
	function(rr, &http.Request{})
}
