package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/twonegatives/drweb_challenge/pkg/callbacks"
	"github.com/twonegatives/drweb_challenge/pkg/config"
	"github.com/twonegatives/drweb_challenge/pkg/drweb"
	"github.com/twonegatives/drweb_challenge/pkg/encoders"
	"github.com/twonegatives/drweb_challenge/pkg/storages"
)

var storage = storages.FileSystemStorage{
	BasePath: ".",
	FileMode: 0700,
}

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
	file, err := drweb.NewFile(r.Body, &storage, &encoders.SHA256Encoder{})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	_, err = storage.Save(file)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(file)
}

func retrieveFile(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	filename := vars["hashstring"]
	file, err := storage.Load(filename)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if _, err := io.Copy(w, file.Body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func deleteFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if err := storage.Delete(vars["hashstring"]); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
