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

	// Serve static files from the React build directory
	fs := http.FileServer(http.Dir(filepath.Join("frontend-react", "dist")))
	http.Handle("/", backend.Cors(func(w http.ResponseWriter, r *http.Request) {
		// For any path that doesn't match our API routes, serve the React app
		if r.URL.Path != "/" && !fileExists(filepath.Join("frontend-react", "dist", r.URL.Path)) {
			// If the file doesn't exist, serve index.html (for client-side routing)
			http.ServeFile(w, r, filepath.Join("frontend-react", "dist", "index.html"))
			return
		}
		fs.ServeHTTP(w, r)
	}))

	fmt.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Helper function to check if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
