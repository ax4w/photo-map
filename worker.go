package main

import (
	"crypto/sha256"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

//go:embed generate.sh
var generateThumbsScript []byte

func fsWorker() {
	for {
		fsWorkerLogic()
		time.Sleep(1 * time.Minute)
	}
}

func insertNewRegion(name string) {
	var (
		resp, err = http.Get(fmt.Sprintf("https://nominatim.openstreetmap.org/search?q=%s&format=json&limit=1", name))
		jsonData  []map[string]any
	)
	if err != nil {
		println("error in get", err.Error())
		return
	}
	bodyInBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		println("error reading response", err.Error())
		return
	}
	if err := json.Unmarshal(bodyInBytes, &jsonData); err != nil {
		println(err.Error())
		return
	}
	if len(jsonData) == 0 {
		return
	}
	var (
		lat, _  = strconv.ParseFloat(jsonData[0]["lat"].(string), 64)
		long, _ = strconv.ParseFloat(jsonData[0]["lon"].(string), 64)
	)
	pgConn.Create(&Region{
		Name: name,
		Lat:  lat,
		Long: long,
	})
}

func fsWorkerLogic() {
	var (
		needsRegeneration bool
		scriptPath        = filepath.Join(basePath, "generate.sh")
	)
	if _, err := os.Stat(imagesBasePath); errors.Is(err, os.ErrNotExist) {
		os.MkdirAll(imagesBasePath, os.ModePerm)
		println("Creating Image Folder")
	}
	if _, err := os.Stat(thumbnailsBasePath); errors.Is(err, os.ErrNotExist) {
		os.MkdirAll(thumbnailsBasePath, os.ModePerm)
		println("Creating Thumbs Folder")
	}
	if _, err := os.Stat(scriptPath); errors.Is(err, os.ErrNotExist) {
		os.WriteFile(scriptPath, generateThumbsScript, os.ModePerm)
		println("Creating generate.sh script")
		err := exec.Command("chmod", "+x", scriptPath).Run()
		if err != nil {
			println("error chmod +x file", err.Error())
		}
	}
	entries, err := os.ReadDir(imagesBasePath)
	if err != nil {
		println("error reading dir", err.Error())
		return
	}
	for _, v := range entries {
		var (
			info, _ = v.Info()
			str     = fmt.Sprintf("%d%d", info.ModTime().UnixMilli(), info.Size())
			hash    = fmt.Sprintf("%x", sha256.Sum256([]byte(str)))
			region  = Region{}
			tx      = pgConn.Where("name = ?", v.Name()).First(&region)
		)
		if region.Hash == "" {
			println("found nothing for", v.Name())
			insertNewRegion(v.Name())
		}
		if hash != region.Hash {
			needsRegeneration = true
			tx = pgConn.Model(&Region{}).Where("name = ?", v.Name()).Update("hash", hash)
			if tx.Error != nil {
				println(tx.Error.Error())
			}
		}
	}
	if needsRegeneration {
		println("regenerating thumbnails")
		err := exec.Command("/bin/sh", scriptPath).Run()
		if err != nil {
			println("error running generate", err.Error())
		}
	}
}
