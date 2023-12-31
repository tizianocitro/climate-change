package controller

import (
	"encoding/csv"
	"math"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/tizianocitro/climate-change/cc-data-provider/data"
	"github.com/tizianocitro/climate-change/cc-data-provider/model"
	"github.com/tizianocitro/climate-change/cc-data-provider/util"
	"golang.org/x/exp/slices"
)

type SeaController struct{}

func NewSeaController() *SeaController {
	return &SeaController{}
}

func (sc *SeaController) GetSeas(c *fiber.Ctx) error {
	organizationId := c.Params("organizationId")
	tableData := model.PaginatedTableData{
		Columns: temperaturesPaginatedTableData.Columns,
		Rows:    []model.PaginatedTableRow{},
	}
	for _, dioxide := range seasMap[organizationId] {
		tableData.Rows = append(tableData.Rows, model.PaginatedTableRow{
			ID:          dioxide.ID,
			Name:        dioxide.Name,
			Description: dioxide.Description,
		})
	}
	return c.JSON(tableData)
}

func (sc *SeaController) GetSea(c *fiber.Ctx) error {
	return c.JSON(getSeaByID(c))
}

func (sc *SeaController) GetSeaDescription(c *fiber.Ctx) error {
	seaId := c.Params("seaId")
	return c.JSON(fiber.Map{"text": seasDescriptionMap[seaId]})
}

func (sc *SeaController) GetSeaMap(c *fiber.Ctx) error {
	sea := getSeaByID(c)
	if sea == (model.Sea{}) {
		return c.JSON(model.MapData{})
	}
	year := c.Query("year")
	if !isYearInSeaRange(year) {
		return c.JSON(model.MapData{})
	}
	if sea.Name == "World" {
		mapData, err := getWorldSeaData(year)
		if err != nil {
			return err
		}
		return c.JSON(fiber.Map{
			"seaEnv": model.SeaEnv{
				Label:            mapData.SeaEnv.Label,
				Value:            mapData.SeaEnv.Value,
				NoCountriesValue: mapData.SeaEnv.NoCountriesValue,
				CountriesColor:   mapData.SeaEnv.CountriesColor,
				Range:            mapData.SeaEnv.Range,
				ColorRange:       mapData.SeaEnv.ColorRange,
			},
			"points": mapData.Points,
		})
	}
	return c.JSON(model.MapData{})
}

func getSeaByID(c *fiber.Ctx) model.Sea {
	organizationId := c.Params("organizationId")
	seaId := c.Params("seaId")
	for _, sea := range seasMap[organizationId] {
		if sea.ID == seaId {
			return sea
		}
	}
	return model.Sea{}
}

func getWorldSeaData(year string) (model.MapData, error) {
	filePath, err := util.GetEmbeddedFilePathByName("change_in_mean_sea_levels.csv")
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

	sum := 0.0
	divider := 0
	for index, row := range records {
		if index == 0 {
			continue
		}
		date := row[11]
		if !strings.Contains(date, year) {
			continue
		}
		value := row[12]
		valueAsNumber, err := strconv.ParseFloat(value, 64)
		if err != nil {
			continue
		}
		sum += valueAsNumber
		divider += 1
	}

	points := []model.Point{}
	for index, row := range records {
		if index == 0 {
			continue
		}
		date := row[11]
		year := date[7:]
		point := model.Point{
			Label: year,
			Value: year,
		}
		if !slices.Contains(points, point) {
			points = append(points, point)
		}
	}
	return model.MapData{
		SeaEnv: model.SeaEnv{
			Label:            "Sea",
			Value:            sum / float64(divider),
			Range:            getSeaRangeAcrossYears(records),
			CountriesColor:   "#8B4513",
			NoCountriesValue: true,
			ColorRange:       []string{"#87CEEB", "#000080"},
		},
		Points: model.PointData{
			DefaultPoint: points[len(points)-1],
			Points:       points,
		},
	}, nil
}

func getSeaRangeAcrossYears(records [][]string) []float64 {
	yearAverages := make(map[string]model.YearAverage)
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

	// create a slice of all keys-value pairs in map and append all them to the slice
	averages := make([]model.YearAverage, 0, len(yearAverages))
	for _, value := range yearAverages {
		averages = append(averages, value)
	}

	min := 0.0
	max := 0.0
	for index, average := range averages {
		value := average.Sum / float64(average.Divider)
		if index == 0 {
			min = value
			max = value
			continue
		}
		min = math.Min(value, min)
		max = math.Max(value, max)
	}
	return []float64{min, max}
}

func isYearInSeaRange(year string) bool {
	yearAsNumber, err := strconv.Atoi(year)
	if err != nil {
		return false
	}
	return yearAsNumber >= 1992 && yearAsNumber <= 2022
}

var seasMap = map[string][]model.Sea{
	"1": {
		{
			ID:          "04e49629-406f-401e-a650-f577b5b4a949",
			Name:        "World",
			Description: "Change in mean sea levels in the world",
		},
	},
}

var seasDescriptionMap = map[string]string{
	"04e49629-406f-401e-a650-f577b5b4a949": `This data gives estimates of the rise of global sea levels, based on measurements from satellite radar altimeters. These are produced by measuring the time it takes a radar pulse to make a round-trip from the satellite to the sea surface and back again.
	The data is obtained by four satellite altimeters (TOPEX/Poseidon, Jason-1, Jason-2, Jason-3), which have monitored the same surfaces since 1992.`,
}

var seaPaginatedTableData = model.PaginatedTableData{
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
