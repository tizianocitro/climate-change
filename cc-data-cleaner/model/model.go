package model

type WritableData struct {
	Header []string
	Rows   [][]string
}

type Average struct {
	Sum     float64
	Divider int
}
