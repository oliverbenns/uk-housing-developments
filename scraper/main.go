package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/oliverbenns/uk-housing-developments/scraper/internal/location"
	"github.com/oliverbenns/uk-housing-developments/scraper/internal/scraper"
	"github.com/oliverbenns/uk-housing-developments/scraper/internal/service"
	"github.com/redis/go-redis/v9"
	"googlemaps.github.io/maps"
)

func main() {
	ctx := context.Background()

	googleMapsClient, err := createGoogleMapsClient()
	if err != nil {
		log.Fatalf("could not create google maps client: %v", err)
	}

	redisClient, err := createRedisClient(ctx)
	if err != nil {
		log.Fatalf("could not create redis client: %v", err)
	}

	locationClient := &location.Client{
		GoogleMapsClient: googleMapsClient,
		RedisClient:      redisClient,
	}

	svc := service.Service{
		LocationClient: locationClient,
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

func createGoogleMapsClient() (*maps.Client, error) {
	googleMapsApiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	if googleMapsApiKey == "" {
		return nil, fmt.Errorf("missing google maps url")
	}

	return maps.NewClient(maps.WithAPIKey(googleMapsApiKey))
}

func createRedisClient(ctx context.Context) (*redis.Client, error) {
	redisUrl := os.Getenv("REDIS_URL")
	if redisUrl == "" {
		return nil, fmt.Errorf("missing redis url")
	}

	opt, err := redis.ParseURL(redisUrl)
	if err != nil {
		return nil, err
	}

	redisClient := redis.NewClient(opt)

	_, err = redisClient.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return redisClient, nil
}
