package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
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
	createFile := drweb.CreateFileHandler(&storage, &filenamegenerator)
	router.HandleFunc("/files", drweb.WithCallbacks(createFile, &startSaveCbk, &finishSaveCbk)).Methods("POST")
	router.HandleFunc("/files/{hashstring}", drweb.RetrieveFileHandler(&storage)).Methods("GET")
	router.HandleFunc("/files/{hashstring}", drweb.DeleteFileHandler(&storage)).Methods("DELETE")

	srv := &http.Server{
		Handler:      router,
		Addr:         cfg.GetString("LISTEN"),
		WriteTimeout: cfg.GetDuration("WRITE_TIMOUT") * time.Second,
		ReadTimeout:  cfg.GetDuration("READ_TIMEOUT") * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
