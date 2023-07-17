package converter

import (
	"encoding/csv"
	"os"
	"strconv"

	"github.com/tizianocitro/climate-change/cc-data-cleaner/model"
)

const (
	seaFirstYear int = 1992
	seaLastYear  int = 2022
)

type SeaConverter struct{}

func NewSeaConverter() *SeaConverter {
	return &SeaConverter{}
}

func (sc *SeaConverter) ConvertSeaLevels(inputPath string) (model.WritableData, error) {
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

	seaHeader := []string{"Country"}
	for i := seaFirstYear; i <= seaLastYear; i++ {
		seaHeader = append(seaHeader, strconv.Itoa(i))
	}
	yearAverages := make(map[string]model.Average)
	for index, row := range records {
		if index == 0 {
			continue
		}
		value := row[12]
		valueAsNumber, err := strconv.ParseFloat(value, 64)
		if err != nil {
			continue
		}
		date := row[11]
		year := date[7:]
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

	seaRows := [][]string{}
	seaRow := []string{"World"}
	for i := seaFirstYear; i <= seaLastYear; i++ {
		yearAverage := yearAverages[strconv.Itoa(i)]
		seaRow = append(seaRow, strconv.FormatFloat(yearAverage.Sum/float64(yearAverage.Divider), 'f', -1, 64))
	}
	seaRows = append(seaRows, seaRow)

	return model.WritableData{
		Header: seaHeader,
		Rows:   seaRows,
	}, nil
}
