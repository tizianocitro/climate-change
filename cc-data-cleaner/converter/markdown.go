package converter

import (
	"fmt"
	"io/ioutil"

	"github.com/tizianocitro/climate-change/cc-data-cleaner/model"
)

type MarkdownConverter struct{}

func NewMarkdownConverter() *MarkdownConverter {
	return &MarkdownConverter{}
}

func (mc *MarkdownConverter) ToTable(data model.WritableData, outputPath string) (bool, error) {
	table := "|"
	for _, header := range data.Header {
		table += fmt.Sprintf(" %s |", header)
	}
	table += "\n|"
	for range data.Header {
		table += " --- |"
	}

	for _, row := range data.Rows {
		table += "\n|"
		for _, cell := range row {
			table += " " + cell + " |"
		}
	}

	err := ioutil.WriteFile(outputPath, []byte(table), 0644)
	if err != nil {
		fmt.Printf("Error writing table to file due to %s", err.Error())
		return false, err
	}

	return true, nil
}
