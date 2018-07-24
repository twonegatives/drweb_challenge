package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/twonegatives/drweb_challenge/pkg/config"
	"github.com/twonegatives/drweb_challenge/pkg/drweb"
	"github.com/twonegatives/drweb_challenge/pkg/encoders"
	"github.com/twonegatives/drweb_challenge/pkg/filesavehooks"
	"github.com/twonegatives/drweb_challenge/pkg/storages"
)

var storage = storages.FileSystemStorage{
	BasePath: ".",
	FileMode: 0700,
}

func main() {
	cfg := config.NewConfig()

	router := mux.NewRouter()
	router.HandleFunc("/files", CreateFile).Methods("POST")
	router.HandleFunc("/files/{hashstring}", RetrieveFile).Methods("GET")
	router.HandleFunc("/files/{hashstring}", DeleteFile).Methods("DELETE")

	srv := &http.Server{
		Handler:      router,
		Addr:         cfg.GetString("LISTEN"),
		WriteTimeout: cfg.GetDuration("WRITE_TIMOUT") * time.Second,
		ReadTimeout:  cfg.GetDuration("READ_TIMEOUT") * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func CreateFile(w http.ResponseWriter, r *http.Request) {
	file, err := drweb.NewFile(r.Body, &storage, &filesavehooks.PrintlnHook{}, &encoders.SHA256Encoder{})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	_, err = file.Save()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(file)
}

func RetrieveFile(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	file, err := storage.Load(vars["hashstring"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	io.Copy(w, file.Body)
}

func DeleteFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if err := storage.Delete(vars["hashstring"]); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
