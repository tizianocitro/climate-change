package app

type WritableData interface{}

type CSVData struct {
	WritableData
	Header []string
	Rows   [][]string
}
