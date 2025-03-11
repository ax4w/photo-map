package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	basePath           = "map-data"
	imagesBasePath     = filepath.Join(basePath, "images")
	thumbnailsBasePath = filepath.Join(basePath, "thumbs")
	pgConn             *gorm.DB
)

type latlong struct {
	Lat, Long float64
}

func serveFile(w http.ResponseWriter, r *http.Request, folder string, pathParts []string) {
	if len(pathParts) < 2 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	region, filename := strings.ToLower(pathParts[0]), pathParts[1]
	if !allowedRegion(region) ||
		strings.Contains(region, "..") ||
		strings.Contains(filename, "..") {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	http.ServeFile(w, r, filepath.Join(folder, region, filename))
}

func isImage(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp":
		return true
	}
	return false
}

func allowedRegion(s string) bool {
	var (
		region Region
		tx     = pgConn.Where("name = ?", s).First(&region)
	)
	if tx.Error != nil {
		println(tx.Error.Error())
		return false
	}
	return region.Hash != ""
}

func init() {
	var dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("host"), os.Getenv("user"), os.Getenv("password"), os.Getenv("dbname"), os.Getenv("port"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}
	db.AutoMigrate(&Region{})
	pgConn = db
}

func main() {
	go fsWorker()
	http.HandleFunc("/api/images/", handleImageAPI)
	http.HandleFunc("/api/regions/", handleRegionsAPI)
	http.HandleFunc("/images/", handleImageServe)
	http.HandleFunc("/thumbs/", handleThumbServe)
	http.HandleFunc("/", handleWebsite)

	fmt.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleWebsite(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./index.html")
}
