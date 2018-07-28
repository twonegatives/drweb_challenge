package main

import (
	"encoding/json"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/twonegatives/drweb_challenge/pkg/callbacks"
	"github.com/twonegatives/drweb_challenge/pkg/config"
	"github.com/twonegatives/drweb_challenge/pkg/drweb"
	"github.com/twonegatives/drweb_challenge/pkg/namegenerators"
	"github.com/twonegatives/drweb_challenge/pkg/pathgenerators"
	"github.com/twonegatives/drweb_challenge/pkg/storages"
)

func main() {
	cfg := config.GetConfig()

	log.SetFormatter(&log.JSONFormatter{})

	pathgen := pathgenerators.NestedGenerator{
		Levels:       cfg.GetInt("PATH_NESTED_LEVELS"),
		FolderLength: cfg.GetInt("PATH_NESTED_FOLDERS_LENGTH"),
		BasePath:     cfg.GetString("PATH_BASE"),
	}

	storage := storages.FileSystemStorage{
		FileMode:          os.FileMode(cfg.GetInt("STORAGE_FILE_MODE")),
		FilePathGenerator: &pathgen,
	}

	filenamegenerator := namegenerators.SHA256{}

	router := mux.NewRouter()
	startSaveCbk := callbacks.LogCallback{Content: "Started to save a file"}
	finishSaveCbk := callbacks.LogCallback{Content: "Finished file saving process"}
	createFile := CreateFileHandler(&storage, &filenamegenerator)
	router.HandleFunc("/files", WithCallbacks(createFile, &startSaveCbk, &finishSaveCbk)).Methods("POST")
	router.HandleFunc("/files/{hashstring}", RetrieveFileHandler(&storage)).Methods("GET")
	router.HandleFunc("/files/{hashstring}", DeleteFileHandler(&storage)).Methods("DELETE")

	srv := &http.Server{
		Handler:      router,
		Addr:         cfg.GetString("LISTEN"),
		WriteTimeout: cfg.GetDuration("WRITE_TIMOUT") * time.Second,
		ReadTimeout:  cfg.GetDuration("READ_TIMEOUT") * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func WithCallbacks(handler func(http.ResponseWriter, *http.Request), before drweb.Callback, after drweb.Callback) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		before.Invoke(w, r)
		defer after.Invoke(w, r)
		handler(w, r)
	}
}

func CreateFileHandler(storage drweb.Storage, filenamegenerator drweb.FileNameGenerator) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var formFile multipart.File
		var formFileHeader *multipart.FileHeader
		var file *drweb.File
		var filename string
		var err error

		if formFile, formFileHeader, err = r.FormFile("file"); err != nil {
			log.WithError(err).Error("failed to get a form file")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		file = &drweb.File{
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

func RetrieveFileHandler(storage drweb.Storage) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		var file *drweb.File
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

func DeleteFileHandler(storage drweb.Storage) func(http.ResponseWriter, *http.Request) {
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
