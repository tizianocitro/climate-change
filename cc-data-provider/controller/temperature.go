package controller

import (
	"encoding/csv"
	"errors"
	"log"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/tizianocitro/climate-change/cc-data-provider/data"
	"github.com/tizianocitro/climate-change/cc-data-provider/model"
	"github.com/tizianocitro/climate-change/cc-data-provider/util"
)

type TemperatureController struct{}

func NewTemperatureController() *TemperatureController {
	return &TemperatureController{}
}

func (tc *TemperatureController) GetTemperatures(c *fiber.Ctx) error {
	organizationId := c.Params("organizationId")
	tableData := model.PaginatedTableData{
		Columns: temperaturesPaginatedTableData.Columns,
		Rows:    []model.PaginatedTableRow{},
	}
	for _, temperature := range temperaturesMap[organizationId] {
		tableData.Rows = append(tableData.Rows, model.PaginatedTableRow{
			ID:          temperature.ID,
			Name:        temperature.Name,
			Description: temperature.Description,
		})
	}
	return c.JSON(tableData)
}

func (tc *TemperatureController) GetTemperature(c *fiber.Ctx) error {
	return c.JSON(getTemperatureByID(c))
}

func (tc *TemperatureController) GetTemperatureDescription(c *fiber.Ctx) error {
	temperatureId := c.Params("temperatureId")
	return c.JSON(fiber.Map{"text": temperaturesDescriptionMap[temperatureId]})
}

func (tc *TemperatureController) GetTemperatureMap(c *fiber.Ctx) error {
	temperature := getTemperatureByID(c)
	if temperature == (model.Temperature{}) {
		return c.JSON(model.MapData{})
	}
	year := c.Query("year")
	if !isYearInTemperaturesMapRange(year) {
		return c.JSON(model.MapData{})
	}
	if temperature.Name == "World" {
		mapData, err := getWorldTemperatureMapData(year)
		if err != nil {
			return err
		}
		return c.JSON(fiber.Map{
			"items":      mapData.Items,
			"points":     mapData.Points,
			"range":      mapData.Range,
			"colorRange": mapData.ColorRange,
		})
	}
	return c.JSON(model.MapData{})
}

func (tc *TemperatureController) GetTemperatureChart(c *fiber.Ctx) error {
	temperature := getTemperatureByID(c)
	if temperature == (model.Temperature{}) {
		return c.JSON(model.SimpleLineChartData{})
	}
	if temperature.Name == "World" {
		chartData, err := getWorldTemperatureChartData()
		if err != nil {
			return err
		}
		return c.JSON(chartData)
	}
	return c.JSON(model.SimpleLineChartData{})
}

func getTemperatureByID(c *fiber.Ctx) model.Temperature {
	organizationId := c.Params("organizationId")
	temperatureId := c.Params("temperatureId")
	for _, temperature := range temperaturesMap[organizationId] {
		if temperature.ID == temperatureId {
			return temperature
		}
	}
	return model.Temperature{}
}

func getWorldTemperatureMapData(year string) (model.MapData, error) {
	filePath, err := util.GetEmbeddedFilePathByName("annual_surface_temperature_change.csv")
	if err != nil {
		return model.MapData{}, err
	}
	file, err := data.Data.Open(filePath)
	if err != nil {
		return model.MapData{}, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return model.MapData{}, err
	}

	items := []model.Country{}
	for index, row := range records {
		if index == 0 {
			continue
		}

		yearIndex := getYearIndex(year)
		if yearIndex < 0 {
			log.Println("Year is less than 0")
			return model.MapData{}, errors.New("Year is less than 0")
		}
		rowYear := row[yearIndex]
		if rowYear == "" {
			rowYear = "0"
		}
		yearAsNumber, err := strconv.ParseFloat(rowYear, 64)
		if err != nil {
			log.Println("Error converting year from string to float64", err)
			return model.MapData{}, err
		}
		items = append(items, model.Country{
			ID:      row[0],
			Iso3:    row[3],
			Country: row[1],
			Value:   yearAsNumber,
		})
	}

	points := []model.Point{}
	for _, column := range records[0] {
		if strings.HasPrefix(column, "F") {
			year := column[1:]
			points = append(points, model.Point{
				Label: year,
				Value: year,
			})
		}
	}

	return model.MapData{
		Items: items,
		Points: model.PointData{
			DefaultPoint: points[len(points)-1],
			Points:       points,
		},
		Range:      getTemperatureRangeAcrossYears(records),
		ColorRange: []string{"#0000ff", "#ff0000"},
	}, nil
}

func getTemperatureRangeAcrossYears(records [][]string) []float64 {
	min := 0.0
	max := 0.0
	firstYearIndex := getYearIndex("1961")
	lastYearIndex := getYearIndex("2022")
	for index, row := range records {
		if index == 0 {
			continue
		}
		for i := firstYearIndex; i <= lastYearIndex; i++ {
			rowYear := row[i]
			if rowYear == "" {
				rowYear = "0"
			}
			yearAsNumber, err := strconv.ParseFloat(rowYear, 64)
			if err != nil {
				log.Println("Error converting year from string to float64 to find min and max", err)
				return []float64{}
			}
			if index == 1 {
				min = yearAsNumber
				max = yearAsNumber
				continue
			}
			min = math.Min(yearAsNumber, min)
			max = math.Max(yearAsNumber, max)
		}
	}
	return []float64{min, max}
}

func getYearIndex(year string) int {
	yearAsNumber, err := strconv.Atoi(year)
	if err != nil {
		log.Println("Error converting year from string to int", err)
		return -1
	}
	return (yearAsNumber - 2022) + 71
}

func isYearInTemperaturesMapRange(year string) bool {
	yearAsNumber, err := strconv.Atoi(year)
	if err != nil {
		return false
	}
	return yearAsNumber >= 1961 && yearAsNumber <= 2022
}

func getWorldTemperatureChartData() (model.SimpleLineChartData, error) {
	filePath, err := util.GetEmbeddedFilePathByName("surface_temperature_change_due_co2.csv")
	if err != nil {
		return model.SimpleLineChartData{}, err
	}
	file, err := data.Data.Open(filePath)
	if err != nil {
		return model.SimpleLineChartData{}, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return model.SimpleLineChartData{}, err
	}

	yearAverages := make(map[string]model.YearAverage)
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
			yearAverages[year] = model.YearAverage{
				Sum:     yearAverage.Sum + valueAsNumber,
				Divider: yearAverage.Divider + 1,
			}
			continue
		}
		yearAverages[year] = model.YearAverage{
			Sum:     valueAsNumber,
			Divider: 1,
		}
	}

	lines := []model.SimpleLineChartValue{}
	for key, value := range yearAverages {
		lines = append(lines, model.SimpleLineChartValue{
			Label: key,
			St:    value.Sum / float64(value.Divider),
		})
	}

	sort.Slice(lines, func(i, j int) bool {
		return lines[i].Label < lines[j].Label
	})

	return model.SimpleLineChartData{
		LineData: lines,
		LineColor: model.LineColor{
			St: "#ff5233",
		},
	}, nil
}

var temperaturesMap = map[string][]model.Temperature{
	"1": {
		{
			ID:          "2ce53d5c-4bd4-4f02-89cc-d5b8f551770c",
			Name:        "World",
			Description: "Annual surface temperature change in the world",
		},
	},
	"2": {
		{
			ID:          "43d5bc63-4f2f-4098-9e97-5df06149a218",
			Name:        "World",
			Description: "Surface temperature change due to CO2 in the world",
		},
	},
}

var temperaturesDescriptionMap = map[string]string{
	"2ce53d5c-4bd4-4f02-89cc-d5b8f551770c": `This data presents the annual surface temperature change by country during the period 1961-2022.
	This data is provided by the Food and Agriculture Organization Corporate Statistical Database (FAOSTAT) and is based on publicly available GISTEMP data from the National Aeronautics and Space Administration Goddard Institute for Space Studies (NASA GISS).`,
	"43d5bc63-4f2f-4098-9e97-5df06149a218": `This data comes from a submission to the data science competition "Data Science vs Fake News" put on by Data World, KDNuggets, and Data4Democracy.
	The extracted data highlights the land average temperature per year due to C02 from 1958 to 2015.`,
}

var temperaturesPaginatedTableData = model.PaginatedTableData{
	Columns: []model.PaginatedTableColumn{
		{
			Title: "Name",
		},
		{
			Title: "Description",
		},
	},
	Rows: []model.PaginatedTableRow{},
}
