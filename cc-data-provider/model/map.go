package model

type MapData struct {
	Items  []Country `json:"items"`
	Points PointData `json:"points"`
}

type Country struct {
	ID      string  `json:"id"`
	Iso3    string  `json:"iso3"`
	Country string  `json:"country"`
	Value   float64 `json:"value"`
}

type PointData struct {
	DefaultPoint Point   `json:"defaultPoint"`
	Points       []Point `json:"points"`
}

type Point struct {
	Label string `json:"label"`
	Value string `json:"value"`
}
