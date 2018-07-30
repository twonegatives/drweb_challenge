package drweb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func writeJSONError(writer *http.ResponseWriter, err error) {
	json.NewEncoder(*writer).Encode(map[string]string{"error": err.Error()})
}

func CreateFileHandler(storage Storage, filenamegenerator FileNameGenerator) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var formFile multipart.File
		var file *FileCreateRequest
		var filename string
		var err error

		w.Header().Set("Content-Type", "application/json")

		if formFile, _, err = r.FormFile("file"); err != nil {
			log.WithError(err).Error("failed to get a form file")
			w.WriteHeader(http.StatusBadRequest)
			writeJSONError(&w, err)
			return
		}
		defer formFile.Close()

		file = &FileCreateRequest{
			Body:          formFile,
			NameGenerator: filenamegenerator,
		}

		if filename, err = storage.Save(file); err != nil {
			log.WithError(err).Error("failed to save file")
			w.WriteHeader(http.StatusInternalServerError)
			writeJSONError(&w, err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"hashstring": filename})
	}
}

func RetrieveFileHandler(storage Storage) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		var err error
		var file *File
		var leadingCnt int

		vars := mux.Vars(req)
		filename := vars["hashstring"]

		if file, err = storage.Load(filename); err != nil {
			w.Header().Set("Content-Type", "application/json")
			if os.IsNotExist(errors.Cause(err)) {
				w.WriteHeader(http.StatusNotFound)
			} else {
				log.WithError(err).Error("failed to load file from storage")
				w.WriteHeader(http.StatusInternalServerError)
			}
			writeJSONError(&w, err)
			return
		}

		defer file.Close()

		max := func(x, y int64) int64 {
			if x < y {
				return x
			}
			return y
		}

		buffer := make([]byte, max(file.Size, 512))
		if leadingCnt, err = file.Body.Read(buffer); err != nil {
			w.Header().Set("Content-Type", "application/json")
			writeJSONError(&w, err)
		}

		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
		w.Header().Set("Content-Type", http.DetectContentType(buffer))
		w.Header().Set("Content-Length", fmt.Sprintf("%d", file.Size))
		_, err = io.Copy(w, io.MultiReader(bytes.NewReader(buffer[0:leadingCnt]), file.Body))

		if err != nil {
			// NOTE: streaming does not leave us much to do in case of failure
			// but to close the connection and assume client will check
			// hashsum or content-length by himself. in any case we can log this
			log.WithError(err).Error("file streaming over http failed")
		}
	}
}

func DeleteFileHandler(storage Storage) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		w.Header().Set("Content-Type", "application/json")
		if err := storage.Delete(vars["hashstring"]); err != nil {
			if os.IsNotExist(errors.Cause(err)) {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			log.WithError(err).Error("failed to delete file from storage")
			w.WriteHeader(http.StatusInternalServerError)
			writeJSONError(&w, err)
			return
		}
	}
}

func WithCallbacks(handler func(http.ResponseWriter, *http.Request), before Callback, after Callback) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		before.Invoke(w, r)
		defer after.Invoke(w, r)
		handler(w, r)
	}
}
