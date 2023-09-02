package main

import (
	"fmt"
	"log"
	"os"

	"github.com/oliverbenns/uk-housing-developments/scraper/internal/scraper"
	"github.com/oliverbenns/uk-housing-developments/scraper/internal/service"
	"googlemaps.github.io/maps"
)

func main() {
	googleMapsApiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	if googleMapsApiKey == "" {
		log.Fatal("missing google maps api key")
	}

	googleMapsClient, err := maps.NewClient(maps.WithAPIKey(googleMapsApiKey))
	if err != nil {
		log.Fatalf("could not create google maps client: %v", err)
	}

	svc := service.Service{
		GoogleMapsClient: googleMapsClient,
		Scrapers: []scraper.Scraper{
			&scraper.Barratt{},
			&scraper.Persimmon{},
			&scraper.Bellway{},
			&scraper.TaylorWimpey{},
			&scraper.Berkeley{},
		},
	}
	data, err := svc.Run()
	if err != nil {
		log.Fatalf("could not run service: %v", err)
	}

	fmt.Println(string(data))
}
