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
	"github.com/twonegatives/drweb_challenge/pkg/namegenerators"
	"github.com/twonegatives/drweb_challenge/pkg/pathgenerators"
	"github.com/twonegatives/drweb_challenge/pkg/storages"
)

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

// NOTE: we use some leading file bytes to generate hash and append unix time to it
// so that it is not required to read the whole file to generate its filename.
// in case we really want to build filename based on the whole file contents
// there is another generator for exactly this purpose: namegenerators.SHA256
var filenamegenerator = namegenerators.LeadingSHA256WithUnixTime{LeadingSize: 50}

func main() {
	cfg := config.NewConfig()

	router := mux.NewRouter()
	startSaveCbk := callbacks.LogCallback{Content: "Started to save a file"}
	finishSaveCbk := callbacks.LogCallback{Content: "Finished file saving process"}
	router.HandleFunc("/files", withCallbacks(createFile, &startSaveCbk, &finishSaveCbk)).Methods("POST")
	router.HandleFunc("/files/{hashstring}", retrieveFile).Methods("GET")
	router.HandleFunc("/files/{hashstring}", deleteFile).Methods("DELETE")

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
		handler(w, r)
		after.Invoke(w, r)
	}
}

func createFile(w http.ResponseWriter, r *http.Request) {
	var formFile multipart.File
	var formFileHeader *multipart.FileHeader
	var file *drweb.File
	var err error

	if formFile, formFileHeader, err = r.FormFile("file"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if file, err = drweb.NewFile(formFile, formFileHeader.Header, &storage, &filenamegenerator); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if _, err = storage.Save(file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(file)
}

func retrieveFile(w http.ResponseWriter, req *http.Request) {
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

func deleteFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if err := storage.Delete(vars["hashstring"]); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
