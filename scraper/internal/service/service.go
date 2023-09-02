package service

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/oliverbenns/uk-housing-developments/scraper/internal/scraper"
)

type Service struct{}

func (s *Service) Run() ([]byte, error) {
	scrapers := []scraper.Scraper{
		//&scraper.Barratt{},
		//&scraper.Persimmon{},
		&scraper.Bellway{},
		//&scraper.TaylorWimpey{},
		//	&scraper.Berkeley{},
	}

	serviceResults := []Result{}

	for _, scraper := range scrapers {
		startTime := time.Now()
		results, err := scraper.Scrape()
		if err != nil {
			return nil, fmt.Errorf("could not scrape %s: %w", scraper.Name(), err)
		}

		for _, result := range results {
			serviceResult := Result{
				Builder:  scraper.Name(),
				Name:     result.Name,
				Url:      result.Url,
				Location: result.Location,
			}

			serviceResults = append(serviceResults, serviceResult)
		}

		elapsed := time.Now().Sub(startTime)
		log.Printf("finished scraping %s, took %s", scraper.Name(), elapsed.String())

	}

	out := Out{
		ScrapedAt: time.Now().UTC(),
		Results:   serviceResults,
	}

	return json.Marshal(out)
}
