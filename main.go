package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"photo-map/backend"
)

func init() {
	backend.InitDB()
	go backend.FileWorker()
}

func main() {
	// API endpoints
	http.HandleFunc("/api/images/", backend.Cors(backend.Images))
	http.HandleFunc("/api/regions/", backend.Cors(backend.Regions))
	http.HandleFunc("/images/", backend.Cors(backend.Image))
	http.HandleFunc("/thumbs/", backend.Cors(backend.Thumbnail))

	fs := http.FileServer(http.Dir(filepath.Join("frontend-react", "dist")))
	http.Handle("/", backend.Cors(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" && !fileExists(filepath.Join("frontend-react", "dist", r.URL.Path)) {
			http.ServeFile(w, r, filepath.Join("frontend-react", "dist", "index.html"))
			return
		}
		fs.ServeHTTP(w, r)
	}))

	fmt.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
