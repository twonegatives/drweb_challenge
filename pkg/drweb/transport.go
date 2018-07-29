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

func writeJSONError(writer *http.ResponseWriter, err error) {
	json.NewEncoder(*writer).Encode(map[string]string{"error": err.Error()})
}

func CreateFileHandler(storage Storage, filenamegenerator FileNameGenerator) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var formFile multipart.File
		var formFileHeader *multipart.FileHeader
		var file *File
		var filename string
		var err error

		w.Header().Set("Content-Type", "application/json")

		if formFile, formFileHeader, err = r.FormFile("file"); err != nil {
			log.WithError(err).Error("failed to get a form file")
			w.WriteHeader(http.StatusBadRequest)
			writeJSONError(&w, err)
			return
		}
		defer formFile.Close()

		file = &File{
			Body:          formFile,
			NameGenerator: filenamegenerator,
			Extension:     filepath.Ext(formFileHeader.Filename),
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
		var file *File
		var err error

		vars := mux.Vars(req)
		filename := vars["hashstring"]

		w.Header().Set("Content-Type", "application/json")

		if file, err = storage.Load(filename); err != nil {
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

		mimetype := "application/octet-stream"
		if file.Extension != "" {
			if inferred := mime.TypeByExtension(file.Extension); inferred != "" {
				mimetype = inferred
			}
		}

		if _, err = io.Copy(w, file.Body); err != nil {
			log.WithError(err).Error("failed to stream file to client")
			w.WriteHeader(http.StatusInternalServerError)
			writeJSONError(&w, err)
			return
		}

		w.Header().Set("Content-Type", mimetype)
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
