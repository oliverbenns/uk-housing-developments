package service

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/oliverbenns/uk-housing-developments/scraper/internal/scraper"
)

type Service struct{}

type ServiceOut struct {
	ScrapedAt time.Time       `json:"scraped_at"`
	Results   []ServiceResult `json:"results"`
}

type ServiceResult struct {
	Builder  string `json:"builder"`
	Name     string `json:"name"`
	Url      string `json:"url"`
	Location string `json:"location"`
}

func (s *Service) Run() ([]byte, error) {
	scrapers := []scraper.Scraper{
		&scraper.Barratt{},
		//&scraper.Persimmon{},
		//&scraper.Bellway{},
		//&scraper.TaylorWimpey{},
		//	&scraper.Berkeley{},
	}

	serviceResults := []ServiceResult{}

	for _, scraper := range scrapers {
		results, err := scraper.Scrape()
		if err != nil {
			return nil, fmt.Errorf("could not scrape %s: %w", scraper.Name(), err)
		}

		for _, result := range results {
			serviceResult := ServiceResult{
				Builder:  scraper.Name(),
				Name:     result.Name,
				Url:      result.Url,
				Location: result.Location,
			}

			serviceResults = append(serviceResults, serviceResult)
		}

		log.Printf("Finished scraping %s", scraper.Name())
	}

	log.Printf("all: %v", serviceResults)

	// @TODO: save to json
	out := ServiceOut{
		ScrapedAt: time.Now().UTC(),
		Results:   serviceResults,
	}

	return json.Marshal(out)
}
