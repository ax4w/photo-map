package backend

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

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

func Regions(w http.ResponseWriter, r *http.Request) {
	var (
		regionsM = make(map[string]latlong)
		regions  []Region
	)
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if tx := pgConn.Find(&regions); tx.Error != nil {
		println("error", tx.Error.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for _, v := range regions {
		regionsM[v.Name] = latlong{Lat: v.Lat, Long: v.Long}
	}
	json.NewEncoder(w).Encode(regionsM)
}

func Images(w http.ResponseWriter, r *http.Request) {
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

	imagePath := filepath.Join(imagesBasePath, region)
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

func Thumbnail(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/thumbs/"), "/")
	serveFile(w, r, thumbnailsBasePath, pathParts)
}

func Image(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/images/"), "/")
	serveFile(w, r, imagesBasePath, pathParts)
}

func Website(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./frontend/index.html")
}
