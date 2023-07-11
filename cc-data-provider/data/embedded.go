package data

import "embed"

//go:embed *.csv
var Data embed.FS
