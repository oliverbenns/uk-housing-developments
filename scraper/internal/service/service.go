package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
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
	mu := sync.Mutex{}
	var wg sync.WaitGroup
	wg.Add(len(scrapers))
	errs := []error{}

	for _, sc := range scrapers {
		go func(closureScraper scraper.Scraper) {
			defer wg.Done()

			startTime := time.Now()
			log.Printf("scraping %s", closureScraper.Name())
			results, err := closureScraper.Scrape()

			mu.Lock()
			defer mu.Unlock()

			if err != nil {
				wrappedErr := fmt.Errorf("could not scrape %s: %w", closureScraper.Name(), err)
				errs = append(errs, wrappedErr)
				return
			}

			// Currently locked through this, could improve.
			for _, result := range results {
				serviceResult := Result{
					Builder:  closureScraper.Name(),
					Name:     result.Name,
					Url:      result.Url,
					Location: result.Location,
				}

				serviceResults = append(serviceResults, serviceResult)
			}

			elapsed := time.Now().Sub(startTime)
			log.Printf("finished scraping %s, took %s", closureScraper.Name(), elapsed.String())

		}(sc)
	}

	wg.Wait()

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	out := Out{
		ScrapedAt: time.Now().UTC(),
		Results:   serviceResults,
	}

	return json.Marshal(out)
}
