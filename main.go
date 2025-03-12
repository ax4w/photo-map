package main

import (
	"fmt"
	"log"
	"net/http"
	"pics/backend"
)

func init() {
	backend.InitDB()
	go backend.FileWorker()
}

func main() {
	http.HandleFunc("/api/images/", backend.Images)
	http.HandleFunc("/api/regions/", backend.Regions)
	http.HandleFunc("/images/", backend.Image)
	http.HandleFunc("/thumbs/", backend.Thumbnail)
	http.HandleFunc("/", backend.Website)

	fmt.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
