package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sort"
	"sync"
	"time"

	"github.com/oliverbenns/uk-housing-developments/scraper/internal/location"
	"github.com/oliverbenns/uk-housing-developments/scraper/internal/scraper"
	"googlemaps.github.io/maps"
)

type Service struct {
	GoogleMapsClient *maps.Client
	Scrapers         []scraper.Scraper
	LocationClient   *location.Client
}

func (s *Service) Run() ([]byte, error) {
	results, err := s.runScrapers()
	if err != nil {
		return nil, err
	}

	// Scrapers run concurrently. Sort so
	// that diff can easily be seen in the json out.
	sort.Sort(ByUrl(results))

	results, err = s.addLatLngs(results)
	if err != nil {
		return nil, err
	}

	out := Out{
		ScrapedAt: time.Now().UTC(),
		Results:   results,
	}

	return json.MarshalIndent(out, "", "  ")
}

func (s *Service) runScrapers() ([]Result, error) {
	serviceResults := []Result{}
	mu := sync.Mutex{}
	var wg sync.WaitGroup
	wg.Add(len(s.Scrapers))
	errs := []error{}

	for _, sc := range s.Scrapers {
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

	return serviceResults, nil
}

func (s *Service) addLatLngs(results []Result) ([]Result, error) {
	resultsWithLocs := []Result{}

	for _, result := range results {
		res, err := s.LocationClient.GetFromAddress(context.Background(), result.Location)
		if err != nil {
			log.Printf("could not get lat/lng for for place %v", result)
			continue
		}

		result.Lat = &res.Lat
		result.Lng = &res.Lng
		log.Print("result with loc", result)
		resultsWithLocs = append(resultsWithLocs, result)
	}

	return resultsWithLocs, nil
}
