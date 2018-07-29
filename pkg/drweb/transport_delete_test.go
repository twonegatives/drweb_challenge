package drweb_test

import (
	"encoding/json"
	"fmt"
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

type deleteFailureCase struct {
	Filename     string
	StorageError error
	ContentType  string
	ServerError  string
	ServerCode   int
}

func TestDeleteFileHandlerFailure(t *testing.T) {
	var objects = map[string]deleteFailureCase{
		"internal error": {
			Filename:     "delete_me_test_main",
			StorageError: errors.Wrap(errors.New("some error"), "failure 500"),
			ContentType:  "application/json",
			ServerError:  "some error",
			ServerCode:   http.StatusInternalServerError,
		},
		"not found": {
			Filename:     "not_exist",
			StorageError: errors.Wrap(os.ErrNotExist, "failure 404"),
			ContentType:  "application/json",
			ServerError:  "",
			ServerCode:   http.StatusNotFound,
		},
	}

	for testName, testObject := range objects {
		t.Run(testName, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			storage := mocks.NewMockStorage(mockCtrl)
			storage.EXPECT().Delete(testObject.Filename).Return(testObject.StorageError)

			req, err := http.NewRequest("DELETE", fmt.Sprintf("/files/%s", testObject.Filename), nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/files/{hashstring}", drweb.DeleteFileHandler(storage))
			router.ServeHTTP(rr, req)

			var response map[string]string
			json.Unmarshal(rr.Body.Bytes(), &response)

			assert.Equal(t, testObject.ServerCode, rr.Code)
			assert.Equal(t, testObject.ContentType, rr.Header().Get("Content-Type"))
			assert.Contains(t, response["error"], testObject.ServerError)
		})
	}
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
