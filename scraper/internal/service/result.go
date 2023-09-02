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
