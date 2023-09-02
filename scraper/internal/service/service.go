package service

import (
	"log"

	"github.com/oliverbenns/uk-housing-developments/scraper/internal/scraper"
)

type Service struct{}

func (s *Service) Run() error {
	scrapers := []scraper.Scraper{
		&scraper.Barratt{},
		&scraper.Persimmon{},
		&scraper.Bellway{},
		//&scraper.TaylorWimpey{},
		//&scraper.Berkeley{},
	}

	all := []scraper.ScrapeResult{}

	for _, scraper := range scrapers {
		results, err := scraper.Scrape()
		if err != nil {
			return err
		}

		all = append(all, results...)
	}

	log.Printf("all: %v", all)

	// @TODO: save to json

	return nil
}
