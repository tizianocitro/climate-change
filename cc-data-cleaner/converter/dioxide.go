package converter

import (
	"encoding/csv"
	"os"
	"strconv"

	"github.com/tizianocitro/climate-change/cc-data-cleaner/model"
)

const (
	dioxideFirstYear int = 1958
	dioxideLastYear  int = 2022
)

type DioxideConverter struct{}

func NewDioxideConverter() *DioxideConverter {
	return &DioxideConverter{}
}

func (dc *DioxideConverter) ConvertDioxideConcentrations(inputPath string) (model.WritableData, error) {
	file, err := os.Open(inputPath)
	if err != nil {
		return model.WritableData{}, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return model.WritableData{}, err
	}

	dioxideHeader := []string{"Country"}
	for i := dioxideFirstYear; i <= dioxideLastYear; i++ {
		dioxideHeader = append(dioxideHeader, strconv.Itoa(i))
	}
	yearAverages := make(map[string]model.Average)
	for index, row := range records {
		unit := row[5]
		if index == 0 || unit == "Percent" {
			continue
		}
		value := row[11]
		valueAsNumber, err := strconv.ParseFloat(value, 64)
		if err != nil {
			continue
		}
		date := row[10]
		year := date[:4]
		if yearAverage, ok := yearAverages[year]; ok {
			yearAverages[year] = model.Average{
				Sum:     yearAverage.Sum + valueAsNumber,
				Divider: yearAverage.Divider + 1,
			}
			continue
		}
		yearAverages[year] = model.Average{
			Sum:     valueAsNumber,
			Divider: 1,
		}
	}

	dioxideRows := [][]string{}
	dioxideRow := []string{"World"}
	for i := dioxideFirstYear; i <= dioxideLastYear; i++ {
		yearAverage := yearAverages[strconv.Itoa(i)]
		dioxideRow = append(dioxideRow, strconv.FormatFloat(yearAverage.Sum/float64(yearAverage.Divider), 'f', -1, 64))
	}
	dioxideRows = append(dioxideRows, dioxideRow)

	return model.WritableData{
		Header: dioxideHeader,
		Rows:   dioxideRows,
	}, nil
}
