package main

import (
	"log"

	"github.com/tizianocitro/climate-change/cc-data-cleaner/converter"
)

func main() {
	CSVConverter := converter.NewCSVConverter()
	markdownConverter := converter.NewMarkdownConverter()

	temperatureConverter := converter.NewTemperatureConverter()
	annualSurfaceTemperature, err := temperatureConverter.
		ConvertAnnualSurfaceTemperature("../cc-data-provider/data/annual_surface_temperature_change.csv")
	if err != nil {
		log.Fatalf("Could not convert annual surface temperature data due to %s", err.Error())
	}
	log.Println("Writing converted annual surface temperature data to CSV")
	result, err := CSVConverter.ToCSV(annualSurfaceTemperature, "./data/annual_surface_temperature_change.csv")
	if err != nil {
		log.Fatalf("Could not write converted annual surface temperature to csv due to %s", err.Error())
	}
	log.Printf("Writing converted annual surface temperature data to CSV completed with result %v", result)
	result, err = markdownConverter.ToTable(annualSurfaceTemperature, "./data/annual_surface_temperature_change.md")
	if err != nil {
		log.Fatalf("Could not write converted annual surface temperature data to table due to %s", err.Error())
	}
	log.Printf("Writing converted annual surface temperature to table completed with result %v", result)

	surfaceTemperatureCO2, err := temperatureConverter.ConvertSurfaceTemperatureCO2("../cc-data-provider/data/surface_temperature_change_due_co2.csv")
	if err != nil {
		log.Fatalf("Could not convert surface temperature CO2 data due to %s", err.Error())
	}
	log.Println("Writing converted surface temperature CO2 data to CSV")
	result, err = CSVConverter.ToCSV(surfaceTemperatureCO2, "./data/surface_temperature_change_due_co2.csv")
	if err != nil {
		log.Fatalf("Could not write converted surface temperature CO2 data to csv due to %s", err.Error())
	}
	log.Printf("Writing converted surface temperature CO2 data to CSV completed with result %v", result)

	log.Println("Writing converted surface temperature CO2 data to table")
	result, err = markdownConverter.ToTable(surfaceTemperatureCO2, "./data/surface_temperature_change_due_co2.md")
	if err != nil {
		log.Fatalf("Could not write converted surface temperature CO2 data to table due to %s", err.Error())
	}
	log.Printf("Writing converted surface temperature CO2 data to table completed with result %v", result)

	seaConverter := converter.NewSeaConverter()
	seaLevels, err := seaConverter.ConvertSeaLevels("../cc-data-provider/data/change_in_mean_sea_levels.csv")
	if err != nil {
		log.Fatalf("Could not convert sea levels data due to %s", err.Error())
	}
	log.Println("Writing converted sea levels data to CSV")
	result, err = CSVConverter.ToCSV(seaLevels, "./data/change_in_mean_sea_levels.csv")
	if err != nil {
		log.Fatalf("Could not write converted sea levels data to csv due to %s", err.Error())
	}
	log.Printf("Writing converted sea levels data to CSV completed with result %v", result)

	log.Println("Writing converted sea levels data to table")
	result, err = markdownConverter.ToTable(seaLevels, "./data/change_in_mean_sea_levels.md")
	if err != nil {
		log.Fatalf("Could not write converted sea levels data to table due to %s", err.Error())
	}
	log.Printf("Writing converted sea levels data to table completed with result %v", result)

	dioxideConverter := converter.NewDioxideConverter()
	dioxideConcentrations, err := dioxideConverter.ConvertDioxideConcentrations("../cc-data-provider/data/atmospheric_dioxide_concentrations.csv")
	if err != nil {
		log.Fatalf("Could not convert dioxide concentrations data due to %s", err.Error())
	}
	log.Println("Writing converted dioxide concentrations data to CSV")
	result, err = CSVConverter.ToCSV(dioxideConcentrations, "./data/atmospheric_dioxide_concentrations.csv")
	if err != nil {
		log.Fatalf("Could not write converted dioxide concentrations data to csv due to %s", err.Error())
	}
	log.Printf("Writing converted dioxide concentrations data to CSV completed with result %v", result)

	log.Println("Writing converted dioxide concentrations data to table")
	result, err = markdownConverter.ToTable(dioxideConcentrations, "./data/atmospheric_dioxide_concentrations.md")
	if err != nil {
		log.Fatalf("Could not write converted dioxide concentrations data to table due to %s", err.Error())
	}
	log.Printf("Writing converted dioxide concentrations data to table completed with result %v", result)
}
