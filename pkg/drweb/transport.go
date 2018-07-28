package drweb

import (
	"encoding/json"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func CreateFileHandler(storage Storage, filenamegenerator FileNameGenerator) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var formFile multipart.File
		var formFileHeader *multipart.FileHeader
		var file *File
		var filename string
		var err error

		if formFile, formFileHeader, err = r.FormFile("file"); err != nil {
			log.WithError(err).Error("failed to get a form file")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		file = &File{
			Body:          formFile,
			NameGenerator: filenamegenerator,
			Extension:     filepath.Ext(formFileHeader.Filename),
		}

		if filename, err = storage.Save(file); err != nil {
			log.WithError(err).Error("failed to save file")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// TODO check json is everywhere
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"Hashstring": filename})
	}
}

func RetrieveFileHandler(storage Storage) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		var file *File
		var err error

		vars := mux.Vars(req)
		filename := vars["hashstring"]

		if file, err = storage.Load(filename); err != nil {
			if os.IsNotExist(errors.Cause(err)) {
				http.Error(w, err.Error(), http.StatusNotFound)
			} else {
				log.WithError(err).Error("failed to load file from storage")
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		mimetype := "application/octet-stream"
		if file.Extension != "" {
			if inferred := mime.TypeByExtension(file.Extension); inferred != "" {
				mimetype = inferred
			}
		}

		w.Header().Set("Content-Type", mimetype)

		if _, err = io.Copy(w, file.Body); err != nil {
			log.WithError(err).Error("failed to stream file to client")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func DeleteFileHandler(storage Storage) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		if err := storage.Delete(vars["hashstring"]); err != nil {
			if os.IsNotExist(errors.Cause(err)) {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}

			log.WithError(err).Error("failed to delete file from storage")
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
