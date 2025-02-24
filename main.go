package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

type latlong struct {
	Lat, Long float32
}

/*var allowedRegions = map[string]latlong{
	"mallorca":   {Lat: 39.6, Long: 3.00},
	"berlin":     {Lat: 52.52, Long: 13.40},
	"bari":       {Lat: 41.11, Long: 16.87},
	"dubrovnik":  {Lat: 42.64, Long: 18.09},
	"montenegro": {Lat: 42.70, Long: 19.37},
	"venice":     {Lat: 45.44, Long: 12.31},
	"zadar":      {Lat: 44.11, Long: 15.23},
}*/

var pgConn *sql.DB

func main() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err.Error())
	}
	pgConn = db
	http.HandleFunc("/api/images/", handleImageAPI)
	http.HandleFunc("/api/regions/", handleRegionsAPI)
	http.HandleFunc("/images/", handleImageServe)
	http.HandleFunc("/", handleWebsite)

	fmt.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleWebsite(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./index.html")
}

func handleRegionsAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	rows, err := pgConn.Query("SELECT * from region")
	if err != nil {
		println(err.Error())
		return
	}
	var regions = make(map[string]latlong)
	for rows.Next() {
		var (
			name string
			lat  float32
			long float32
		)
		err := rows.Scan(&name, &lat, &long)
		if err != nil {
			println(err.Error())
			return
		}
		regions[name] = latlong{Lat: lat, Long: long}
	}
	json.NewEncoder(w).Encode(regions)
}

func allowedRegion(s string) bool {
	rows, err := pgConn.Query("SELECT * from region where name=$1", s)
	if err != nil {
		println(err.Error())
		return false
	}
	return rows.Next()
}

func handleImageAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/images/"), "/")
	if len(pathParts) < 1 {
		http.Error(w, "Missing region", http.StatusBadRequest)
		return
	}

	region := strings.ToLower(pathParts[0])
	if !allowedRegion(region) || strings.Contains(region, "..") {
		http.Error(w, "Invalid region", http.StatusForbidden)
		return
	}

	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 30 {
		limit = 30
	}

	imagePath := filepath.Join("images", region)
	files, err := os.ReadDir(imagePath)
	if err != nil {
		http.Error(w, "Region not found", http.StatusNotFound)
		return
	}

	var images []string
	count := 0
	totalFiles := 0

	for _, file := range files {
		if isImage(file.Name()) {
			totalFiles++
		}
	}

	for _, file := range files {
		if !isImage(file.Name()) {
			continue
		}
		if count >= offset+limit {
			break
		}
		if count >= offset {
			images = append(images, file.Name())
		}
		count++
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"images":   images,
		"has_more": count < totalFiles,
	})
}

func handleImageServe(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/images/"), "/")
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

	http.ServeFile(w, r, filepath.Join("images", region, filename))
}

func isImage(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp":
		return true
	}
	return false
}
