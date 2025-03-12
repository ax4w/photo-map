package backend

import (
	"crypto/sha256"
	_ "embed"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

//go:embed generate.sh
var generateThumbsScript []byte

func FileWorker() {
	for {
		logic()
		time.Sleep(2 * time.Minute)
	}
}

func createFolderIfNotExist(path string) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		println("Creating Image Folder")
		err := os.Mkdir(imagesBasePath, 0755)
		if err != nil {
			println(err.Error())
		}
	}
}

func regenerateAllThumbnails(scriptPath string) {
	println("regenerating thumbnails")
	out, err := exec.Command("/bin/sh", scriptPath).CombinedOutput()
	if err != nil {
		println("error running generate", err.Error())
	}
	println(string(out))
	println("finished rengerating thumbnails")
}

func logic() {
	var (
		needsRegeneration bool
		scriptPath        = filepath.Join(basePath, "generate.sh")
	)

	createFolderIfNotExist(imagesBasePath)
	createFolderIfNotExist(thumbnailsBasePath)

	if _, err := os.Stat(scriptPath); errors.Is(err, os.ErrNotExist) {
		os.WriteFile(scriptPath, generateThumbsScript, 0755)
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
			info, _    = v.Info()
			str        = fmt.Sprintf("%d%d", info.ModTime().UnixMilli(), info.Size())
			hash       = fmt.Sprintf("%x", sha256.Sum256([]byte(str)))
			region, ok = getRegion(v.Name())
		)
		if n := strings.ReplaceAll(v.Name(), ".", ""); len(n) == 0 {
			continue
		}

		if !ok {
			println("found nothing for", v.Name())
			insertNewRegion(v.Name())
		}
		if hash != region.Hash {
			needsRegeneration = true
			updateHash(v.Name(), hash)
		}
	}
	if needsRegeneration {
		regenerateAllThumbnails(scriptPath)
	}
}
