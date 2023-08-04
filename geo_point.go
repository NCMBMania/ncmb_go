package NCMB

type GeoPoint struct {
	Type      string  `json:"__type"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
