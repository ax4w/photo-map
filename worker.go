package main

import (
	"crypto/sha256"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"
)

//go:embed generate.sh
var generateThumbsScript []byte

func fsWorker() {
	for {
		fsWorkerLogic()
		time.Sleep(10 * time.Minute)
	}
}

func insertNewFolder(name string) {
	var resp, err = http.Get(fmt.Sprintf("https://nominatim.openstreetmap.org/search?%s", name))
	if err != nil {
		println(err.Error())
		return
	}
	bodyInBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		println(err.Error())
		return
	}
	println(string(bodyInBytes))
}

func fsWorkerLogic() {
	var (
		needsRegeneration bool
	)
	if _, err := os.Stat(imagesBasePath); errors.Is(err, os.ErrNotExist) {
		os.MkdirAll(imagesBasePath, os.ModePerm)
		println("Creating Image Folder")
	}
	if _, err := os.Stat(thumbnailsBasePath); errors.Is(err, os.ErrNotExist) {
		os.MkdirAll(thumbnailsBasePath, os.ModePerm)
		println("Creating Thumbs Folder")
	}
	if _, err := os.Stat("generate.sh"); errors.Is(err, os.ErrNotExist) {
		os.WriteFile("generate.sh", generateThumbsScript, os.ModePerm)
		println("Creating generate.sh script")
		err := exec.Command("chmod", "+x", "generate.sh").Run()
		if err != nil {
			println(err.Error())
		}
	}
	entries, err := os.ReadDir(imagesBasePath)
	if err != nil {
		println(err.Error())
		return
	}
	for _, v := range entries {
		var (
			info, _ = v.Info()
			str     = fmt.Sprintf("%d%d", info.ModTime().UnixMilli(), info.Size())
			hash    = fmt.Sprintf("%x", sha256.Sum256([]byte(str)))
			region  = Region{Name: v.Name()}
			tx      = pgConn.First(&region)
		)
		if tx.RowsAffected == 0 {
			insertNewFolder(v.Name())
			continue
		}
		if hash != region.Hash {
			needsRegeneration = true
			tx = pgConn.Model(&Region{}).Where("name = ?", v.Name()).Update("hash", hash)
			if tx.Error != nil {
				println(err.Error())
			}
		}
	}
	if needsRegeneration {
		err := exec.Command("/bin/sh", "generate.sh").Run()
		if err != nil {
			println(err.Error())
		}
	}
}
