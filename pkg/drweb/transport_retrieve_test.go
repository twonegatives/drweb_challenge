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
	Filename    string
	Extension   string
	Contents    []byte
	ContentType string
}

type retrieveFailureCase struct {
	Filename     string
	Extension    string
	Contents     []byte
	ContentType  string
	StorageError error
	ServerCode   int
	ServerError  string
}

func TestRetrieveSuccess(t *testing.T) {
	var objects = map[string]retrieveSuccessCase{
		"unknown file type": {
			Filename:    "some_saved_file",
			Extension:   "",
			Contents:    []byte("Loaded file"),
			ContentType: "application/octet-stream",
		},
		"known file type": {
			Filename:    "another_saved_file",
			Extension:   ".jpg",
			Contents:    []byte("Loaded file"),
			ContentType: "image/jpeg",
		},
	}

	for testName, testObject := range objects {
		t.Run(testName, func(t *testing.T) {
			file := drweb.File{Body: ioutil.NopCloser(bytes.NewReader(testObject.Contents)), Extension: testObject.Extension}

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
		})
	}
}

func TestRetrieveFailure(t *testing.T) {
	var objects = map[string]retrieveFailureCase{
		"unknown file type": {
			Filename:     "unexistant_file",
			ContentType:  "application/json",
			ServerCode:   http.StatusNotFound,
			StorageError: errors.Wrap(os.ErrNotExist, "some description"),
			ServerError:  "",
		},
		"known file type": {
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
