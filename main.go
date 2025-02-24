package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var allowedRegions = map[string]bool{
	"mallorca": true,
	"canaries": true,
}

func main() {
	http.HandleFunc("/api/images/", handleImageAPI)
	http.HandleFunc("/images/", handleImageServe)
	http.HandleFunc("/", handleWebsite)

	fmt.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleWebsite(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./index.html")
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
	if !allowedRegions[region] || strings.Contains(region, "..") {
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
	if !allowedRegions[region] ||
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
