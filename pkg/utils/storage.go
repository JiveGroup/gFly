package utils

import (
	"fmt"
	"github.com/gflydev/core/utils"
	"github.com/gflydev/modules/storage"
	storageDto "github.com/gflydev/modules/storage/dto"
	"github.com/gflydev/modules/storagecs3"
	"path/filepath"
	"strings"
)

// LegitimizeUploadedFile make object (Full URL) be validated data
func LegitimizeUploadedFile(objectUrl, dir string) string {
	if objectUrl == "" {
		return ""
	}

	legitimizeItem := storageDto.LegitimizeItem{
		File:          objectUrl,
		Name:          filepath.Base(objectUrl),
		Dir:           dir,
		LegitimizeURL: "",
	}

	// Get Request path of object
	object, _ := utils.RequestPath(objectUrl)

	// Prefix of Contabo Object Storage
	bucketPathCS3 := fmt.Sprintf("/%s:%s/",
		utils.Getenv("CS_BUCKET_CODE", ""),
		utils.Getenv("CS_BUCKET", ""),
	)
	// Prefix of Local Storage
	localPath := fmt.Sprintf("/%s/",
		utils.Getenv("STORAGE_DIR", ""),
	)

	legitimizeItems := make([]storageDto.LegitimizeItem, 0)
	if strings.HasPrefix(object, bucketPathCS3) {
		legitimizeItems = storagecs3.LegitimizeFiles([]storageDto.LegitimizeItem{legitimizeItem})
	} else if strings.HasPrefix(object, localPath) {
		legitimizeItems = storage.LegitimizeFiles([]storageDto.LegitimizeItem{legitimizeItem})
	}

	if len(legitimizeItems) > 0 {
		firstItem := legitimizeItems[0]

		objectUrl = firstItem.LegitimizeURL
	}

	return objectUrl
}
