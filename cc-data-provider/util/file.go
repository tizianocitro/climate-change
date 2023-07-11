package util

import (
	"io/fs"
	"strings"

	"github.com/tizianocitro/climate-change/cc-data-provider/data"
)

func GetEmbeddedFilePathByName(fileName string) (string, error) {
	filePaths, err := fs.Glob(data.Data, "*.csv")
	if err != nil {
		return "", err
	}
	for _, filePath := range filePaths {
		if strings.Contains(filePath, fileName) {
			return filePath, nil
		}
	}
	return "", err
}
