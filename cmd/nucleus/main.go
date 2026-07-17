package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"nucleus/internal/api"
	"nucleus/internal/db"
	"nucleus/internal/scheduler"
)

//go:embed dist/*
var staticFiles embed.FS

type spaHandler struct {
	staticPath string
	indexPath  string
	fs         http.FileSystem
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	// Check if file exists
	f, err := h.fs.Open(path)
	if os.IsNotExist(err) {
		// fallback to index.html
		r.URL.Path = "/"
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		// If file exists, check if it's a directory
		stat, err := f.Stat()
		if err == nil && stat.IsDir() {
			r.URL.Path = "/"
		}
		f.Close()
	}

	http.FileServer(h.fs).ServeHTTP(w, r)
}

func main() {
	db.InitDB("app.db")
	scheduler.InitScheduler()

	mux := http.NewServeMux()

	api.RegisterRoutes(mux)

	// Setup embedded SPA
	distFS, err := fs.Sub(staticFiles, "dist")
	if err != nil {
		// If running locally without build, this might fail, fallback to a simple message
		log.Println("dist/ not found in embedded filesystem. Ensure frontend is built before compiling. Serving API only.")
	} else {
		spa := spaHandler{staticPath: "/", indexPath: "index.html", fs: http.FS(distFS)}
		mux.Handle("/", spa)
	}

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
