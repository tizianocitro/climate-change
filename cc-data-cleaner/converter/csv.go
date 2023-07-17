package converter

import (
	"encoding/csv"
	"os"

	"github.com/tizianocitro/climate-change/cc-data-cleaner/model"
)

type CSVConverter struct{}

func NewCSVConverter() *CSVConverter {
	return &CSVConverter{}
}

func (cc *CSVConverter) ToCSV(data model.WritableData, outputPath string) (bool, error) {
	file, err := os.Create(outputPath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write(data.Header)
	if err != nil {
		return false, err
	}

	for _, row := range data.Rows {
		err := writer.Write(row)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}
