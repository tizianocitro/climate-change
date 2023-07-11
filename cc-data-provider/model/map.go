package model

type MapData struct {
	Items      []Country `json:"items"`
	Points     PointData `json:"points"`
	Range      []float64 `json:"range"`
	ColorRange []string  `json:"colorRange"`
	WorldEnv   WorldEnv  `json:"worldEnv"`
	SeaEnv     SeaEnv    `json:"seaEnv"`
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

type WorldEnv struct {
	Value      float64   `json:"value"`
	Range      []float64 `json:"range"`
	ColorRange []string  `json:"colorRange"`
}

type SeaEnv struct {
	Label            string    `json:"label"`
	Value            float64   `json:"value"`
	CountriesColor   string    `json:"countriesColor"`
	NoCountriesValue bool      `json:"noCountriesValue"`
	Range            []float64 `json:"range"`
	ColorRange       []string  `json:"colorRange"`
}
