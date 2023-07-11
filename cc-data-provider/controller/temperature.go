package controller

import (
	"encoding/csv"
	"errors"
	"io/fs"
	"log"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/tizianocitro/climate-change/cc-data-provider/data"
	"github.com/tizianocitro/climate-change/cc-data-provider/model"
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

func (tc *TemperatureController) GetTemperatureMap(c *fiber.Ctx) error {
	temperature := getTemperatureByID(c)
	if temperature == (model.Temperature{}) {
		return c.JSON(model.MapData{})
	}
	year := c.Query("year")
	if temperature.Name == "World" {
		mapData, err := getWorldTemperatureData(year)
		if err != nil {
			return err
		}
		return c.JSON(mapData)
	}
	return c.JSON(model.MapData{})
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

func getWorldTemperatureData(year string) (model.MapData, error) {
	filePaths, err := fs.Glob(data.Data, "*.csv")
	if err != nil {
		log.Fatal(err)
	}

	// filePath := filepath.Join("data", "annual_surface_temperature_change.csv")
	file, err := data.Data.Open(filePaths[0])
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
			log.Println("Error converting year from string to float64", err, "rowFromYearIndex", row[yearIndex])
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

var temperaturesMap = map[string][]model.Temperature{
	"1": {
		{
			ID:          "2ce53d5c-4bd4-4f02-89cc-d5b8f551770c",
			Name:        "World",
			Description: "Temperatures from all over the world",
		},
	},
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
