package backend

import "path/filepath"

var (
	basePath           = "map-data"
	imagesBasePath     = filepath.Join(basePath, "images")
	thumbnailsBasePath = filepath.Join(basePath, "thumbs")
)
