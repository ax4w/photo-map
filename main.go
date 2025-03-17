package main

import (
	"fmt"
	"github.com/a-h/templ"
	"log"
	"net/http"
	"os"
	"photo-map/backend"
	"photo-map/frontend"
)

func init() {
	backend.InitDB()
	go backend.FileWorker()
}

func main() {
	var index = frontend.Index(os.Getenv("title"))
	http.HandleFunc("/api/images/", backend.Cors(backend.Images))
	http.HandleFunc("/api/regions/", backend.Cors(backend.Regions))
	http.HandleFunc("/images/", backend.Cors(backend.Image))
	http.HandleFunc("/thumbs/", backend.Cors(backend.Thumbnail))
	http.Handle("/", templ.Handler(index))

	fmt.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
