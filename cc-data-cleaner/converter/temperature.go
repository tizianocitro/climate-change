package converter

import (
	"encoding/csv"
	"log"
	"os"
	"sort"
	"strconv"

	"github.com/tizianocitro/climate-change/cc-data-cleaner/model"
)

const (
	firstYear int = 1961
	lastYear  int = 2022
)

type TemperatureConverter struct{}

func NewTemperatureConverter() *TemperatureConverter {
	return &TemperatureConverter{}
}

func (tc *TemperatureConverter) ConvertAnnualSurfaceTemperature(inputPath string) (model.WritableData, error) {
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

	temperatureHeader := []string{"Country"}
	for i := firstYear; i <= lastYear; i++ {
		temperatureHeader = append(temperatureHeader, strconv.Itoa(i))
	}
	temperatureRows := [][]string{}
	for index, row := range records {
		if index == 0 {
			continue
		}
		country := row[1]
		yearRow := []string{country}
		firstYearIndex := getYearIndex(strconv.Itoa(firstYear))
		lastYearIndex := getYearIndex(strconv.Itoa(lastYear))
		for i := firstYearIndex; i <= lastYearIndex; i++ {
			rowYear := row[i]
			if rowYear == "" {
				rowYear = "0"
			}
			yearRow = append(yearRow, rowYear)
		}
		temperatureRows = append(temperatureRows, yearRow)
	}

	return model.WritableData{
		Header: temperatureHeader,
		Rows:   temperatureRows,
	}, nil
}

func (tc *TemperatureConverter) ConvertSurfaceTemperatureCO2(inputPath string) (model.WritableData, error) {
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

	yearAverages := make(map[string]model.Average)
	for index, row := range records {
		if index == 0 {
			continue
		}
		value := row[1]
		valueAsNumber, err := strconv.ParseFloat(value, 64)
		if err != nil {
			continue
		}
		date := row[0]
		year := date[:4]
		if year < "1958" {
			continue
		}
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

	temperatureRows := [][]string{}
	for key, value := range yearAverages {
		temperatureRows = append(temperatureRows, []string{key, strconv.FormatFloat(value.Sum/float64(value.Divider), 'f', -1, 64)})
	}

	sort.Slice(temperatureRows, func(i, j int) bool {
		return temperatureRows[i][0] < temperatureRows[j][0]
	})

	return model.WritableData{
		Header: []string{"Year", "SurfaceTemperatureChange"},
		Rows:   temperatureRows,
	}, nil
}

func getYearIndex(year string) int {
	yearAsNumber, err := strconv.Atoi(year)
	if err != nil {
		log.Println("Error converting year from string to int", err)
		return -1
	}
	return (yearAsNumber - 2022) + 71
}
