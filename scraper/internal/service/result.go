package service

import (
	"time"
)

type Out struct {
	ScrapedAt time.Time `json:"scraped_at"`
	Results   []Result  `json:"results"`
}

type Result struct {
	Builder  string   `json:"builder"`
	Name     string   `json:"name"`
	Url      string   `json:"url"`
	Location string   `json:"location"`
	Lat      *float64 `json:"lat,omitempty"`
	Lng      *float64 `json:"lng,omitempty"`
}

type ByUrl []Result

func (a ByUrl) Len() int {
	return len(a)
}
func (a ByUrl) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a ByUrl) Less(i, j int) bool {
	return a[i].Url < a[j].Url
}
