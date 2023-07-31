package controller

import (
	"encoding/csv"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/tizianocitro/climate-change/cc-data-provider/data"
	"github.com/tizianocitro/climate-change/cc-data-provider/model"
	"github.com/tizianocitro/climate-change/cc-data-provider/util"
	"golang.org/x/exp/slices"
)

type DioxideController struct{}

func NewDioxideController() *DioxideController {
	return &DioxideController{}
}

func (dc *DioxideController) GetAllDioxide(c *fiber.Ctx) error {
	organizationId := c.Params("organizationId")
	tableData := model.PaginatedTableData{
		Columns: temperaturesPaginatedTableData.Columns,
		Rows:    []model.PaginatedTableRow{},
	}
	for _, dioxide := range dioxideMap[organizationId] {
		tableData.Rows = append(tableData.Rows, model.PaginatedTableRow{
			ID:          dioxide.ID,
			Name:        dioxide.Name,
			Description: dioxide.Description,
		})
	}
	return c.JSON(tableData)
}

func (dc *DioxideController) GetDioxide(c *fiber.Ctx) error {
	return c.JSON(getDioxideByID(c))
}

func (dc *DioxideController) GetDioxideDescription(c *fiber.Ctx) error {
	dioxideId := c.Params("dioxideId")
	return c.JSON(fiber.Map{"text": dioxideDescriptionMap[dioxideId]})
}

func (dc *DioxideController) GetDioxideMap(c *fiber.Ctx) error {
	dioxide := getDioxideByID(c)
	if dioxide == (model.Dioxide{}) {
		return c.JSON(model.MapData{})
	}
	year := c.Query("year")
	if !isYearInDioxiteRange(year) {
		return c.JSON(model.MapData{})
	}
	if dioxide.Name == "World" {
		mapData, err := getWorldDioxideData(year)
		if err != nil {
			return err
		}
		return c.JSON(fiber.Map{
			"worldEnv": model.WorldEnv{
				Value:      mapData.WorldEnv.Value,
				Range:      mapData.WorldEnv.Range,
				ColorRange: mapData.WorldEnv.ColorRange,
			},
			"points": mapData.Points,
		})
	}
	return c.JSON(model.MapData{})
}

func (dc *DioxideController) GetDioxideChart(c *fiber.Ctx) error {
	dioxide := getDioxideByID(c)
	if dioxide == (model.Dioxide{}) {
		return c.JSON(model.SimpleLineChartData{})
	}
	if dioxide.Name == "World" {
		chartData, err := getDioxidedChartData()
		if err != nil {
			return err
		}
		return c.JSON(chartData)
	}
	return c.JSON(model.SimpleLineChartData{})
}

func getDioxideByID(c *fiber.Ctx) model.Dioxide {
	organizationId := c.Params("organizationId")
	dioxideId := c.Params("dioxideId")
	for _, dioxide := range dioxideMap[organizationId] {
		if dioxide.ID == dioxideId {
			return dioxide
		}
	}
	return model.Dioxide{}
}

func getWorldDioxideData(year string) (model.MapData, error) {
	filePath, err := util.GetEmbeddedFilePathByName("atmospheric_dioxide_concentrations.csv")
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
		unit := row[5]
		if index == 0 || unit == "Percent" {
			continue
		}
		date := row[10]
		if !strings.Contains(date, year) {
			continue
		}
		value := row[11]
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
		date := row[10]
		year := date[:4]
		point := model.Point{
			Label: year,
			Value: year,
		}
		if !slices.Contains(points, point) {
			points = append(points, point)
		}
	}
	return model.MapData{
		WorldEnv: model.WorldEnv{
			Value:      sum / float64(divider),
			Range:      getDioxideRangeAcrossYears(records),
			ColorRange: []string{"#CCCCCC", "#636363"},
		},
		Points: model.PointData{
			DefaultPoint: points[len(points)-2],
			Points:       points[:len(points)-1],
		},
	}, nil
}

func getDioxideRangeAcrossYears(records [][]string) []float64 {
	yearAverages := make(map[string]model.YearAverage)
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

func isYearInDioxiteRange(year string) bool {
	yearAsNumber, err := strconv.Atoi(year)
	if err != nil {
		return false
	}
	return yearAsNumber >= 1958 && yearAsNumber <= 2022
}

func getDioxidedChartData() (model.SimpleLineChartData, error) {
	filePath, err := util.GetEmbeddedFilePathByName("atmospheric_dioxide_concentrations.csv")
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
		if year > "2022" {
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
			St: "#636363",
		},
	}, nil
}

var dioxideMap = map[string][]model.Dioxide{
	"1": {
		{
			ID:          "64fc5461-b40c-49ca-a177-ddd2e121ffe1",
			Name:        "World",
			Description: "Atmospheric carbon dioxide concentrations in the world",
		},
	},
}

var dioxideDescriptionMap = map[string]string{
	"64fc5461-b40c-49ca-a177-ddd2e121ffe1": `This data presents the concentration of carbon dioxide in the atmosphere, dating back to 1958.
 	The source data for this data comes from the National Oceanic and Atmospheric Association Global Monitoring Laboratory.`,
}

var dioxidePaginatedTableData = model.PaginatedTableData{
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
