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
)

type retrieveSuccessCase struct {
	Filename      string
	Contents      []byte
	ContentType   string
	ContentLength int
}

type retrieveFailureCase struct {
	Filename     string
	Contents     []byte
	ContentType  string
	StorageError error
	ServerCode   int
	ServerError  string
}

func TestRetrieveSuccess(t *testing.T) {
	txtFile, err := ioutil.ReadFile("./../testdata/alice.txt")
	if err != nil {
		t.Fatal(err)
	}

	jpgFile, err := ioutil.ReadFile("./../testdata/gopher.jpg")
	if err != nil {
		t.Fatal(err)
	}

	var objects = map[string]retrieveSuccessCase{
		"text file": {
			Filename:      "text_saved_file",
			Contents:      txtFile,
			ContentType:   "text/plain; charset=utf-8",
			ContentLength: 4094,
		},
		"image file": {
			Filename:      "image_saved_file",
			Contents:      jpgFile,
			ContentType:   "image/jpeg",
			ContentLength: 6707,
		},
	}

	for testName, testObject := range objects {
		t.Run(testName, func(t *testing.T) {
			file := drweb.File{Body: ioutil.NopCloser(bytes.NewReader(testObject.Contents)), Size: int64(len(testObject.Contents))}

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			storage := mocks.NewMockStorage(mockCtrl)
			storage.EXPECT().Load(testObject.Filename).Return(&file, nil)

			req, err := http.NewRequest("GET", fmt.Sprintf("/files/%s", testObject.Filename), nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/files/{hashstring}", drweb.RetrieveFileHandler(storage))
			router.ServeHTTP(rr, req)

			assert.Equal(t, http.StatusOK, rr.Code)
			assert.Equal(t, string(testObject.Contents), rr.Body.String())
			assert.Equal(t, testObject.ContentType, rr.Header().Get("Content-Type"))
			assert.Equal(t, fmt.Sprintf("%d", testObject.ContentLength), rr.Header().Get("Content-Length"))
			assert.Equal(t, fmt.Sprintf("attachment; filename=%s", testObject.Filename), rr.Header().Get("Content-Disposition"))
		})
	}
}

func TestRetrieveFailure(t *testing.T) {
	var objects = map[string]retrieveFailureCase{
		"file does not exist": {
			Filename:     "unexistant_file",
			ContentType:  "application/json",
			ServerCode:   http.StatusNotFound,
			StorageError: errors.Wrap(os.ErrNotExist, "some description"),
			ServerError:  "",
		},
		"file read from storage failed": {
			Filename:     "unlucky_file",
			ContentType:  "application/json",
			ServerCode:   http.StatusInternalServerError,
			StorageError: errors.Wrap(errors.New("some error"), "some description"),
			ServerError:  "some error",
		},
	}

	for testName, testObject := range objects {
		t.Run(testName, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			storage := mocks.NewMockStorage(mockCtrl)
			storage.EXPECT().Load(testObject.Filename).Return(nil, testObject.StorageError)

			req, err := http.NewRequest("GET", fmt.Sprintf("/files/%s", testObject.Filename), nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/files/{hashstring}", drweb.RetrieveFileHandler(storage))
			router.ServeHTTP(rr, req)

			var response map[string]string
			json.Unmarshal(rr.Body.Bytes(), &response)

			assert.Equal(t, testObject.ServerCode, rr.Code)
			assert.Equal(t, testObject.ContentType, rr.Header().Get("Content-Type"))
			assert.Contains(t, response["error"], testObject.ServerError)
		})
	}
}
