package main

import (
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/twonegatives/drweb_challenge/pkg/callbacks"
	"github.com/twonegatives/drweb_challenge/pkg/config"
	"github.com/twonegatives/drweb_challenge/pkg/drweb"
	"github.com/twonegatives/drweb_challenge/pkg/files"
	"github.com/twonegatives/drweb_challenge/pkg/namegenerators"
	"github.com/twonegatives/drweb_challenge/pkg/pathgenerators"
	"github.com/twonegatives/drweb_challenge/pkg/storages"
)

func main() {
	cfg := config.GetConfig()

	// NOTE: we use double folder nesting here in order to overcome
	// issue with too much files lying in a single folder.
	// in case we're ok with a single nesting Levels param may be changed
	var pathgen = pathgenerators.NestedGenerator{
		Levels:       2,
		FolderLength: 2,
		BasePath:     ".",
	}

	var storage = storages.FileSystemStorage{
		FileMode:          0700,
		FilePathGenerator: &pathgen,
	}

	var filenamegenerator = namegenerators.SHA256{}

	router := mux.NewRouter()
	startSaveCbk := callbacks.LogCallback{Content: "Started to save a file"}
	finishSaveCbk := callbacks.LogCallback{Content: "Finished file saving process"}
	createFile := createFileHandler(&storage, &filenamegenerator)
	router.HandleFunc("/files", withCallbacks(createFile, &startSaveCbk, &finishSaveCbk)).Methods("POST")
	router.HandleFunc("/files/{hashstring}", retrieveFileHandler(&storage)).Methods("GET")
	router.HandleFunc("/files/{hashstring}", deleteFileHandler(&storage)).Methods("DELETE")

	srv := &http.Server{
		Handler:      router,
		Addr:         cfg.GetString("LISTEN"),
		WriteTimeout: cfg.GetDuration("WRITE_TIMOUT") * time.Second,
		ReadTimeout:  cfg.GetDuration("READ_TIMEOUT") * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func withCallbacks(handler func(http.ResponseWriter, *http.Request), before drweb.Callback, after drweb.Callback) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		before.Invoke(w, r)
		defer after.Invoke(w, r)
		handler(w, r)
	}
}

func createFileHandler(storage drweb.Storage, filenamegenerator drweb.FileNameGenerator) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var formFile multipart.File
		var formFileHeader *multipart.FileHeader
		var file *drweb.File
		var filename string
		var err error

		if formFile, formFileHeader, err = r.FormFile("file"); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if file, err = files.NewFile(formFile, formFileHeader.Header.Get("Content-Type"), filenamegenerator); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if filename, err = storage.Save(file); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"Hashstring": filename})
	}
}

func retrieveFileHandler(storage drweb.Storage) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		var file *drweb.File
		var err error

		vars := mux.Vars(req)
		filename := vars["hashstring"]

		if file, err = storage.Load(filename); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if _, err = io.Copy(w, file.Body); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func deleteFileHandler(storage drweb.Storage) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		if err := storage.Delete(vars["hashstring"]); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
