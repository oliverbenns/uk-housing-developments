package location

import (
	"context"
	"errors"
	"log"

	"github.com/redis/go-redis/v9"
	"googlemaps.github.io/maps"
)

type Client struct {
	GoogleMapsClient *maps.Client
	RedisClient      *redis.Client
}

type Geometry struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

func (c *Client) GetFromAddress(ctx context.Context, address string) (Geometry, error) {
	res, err := c.get(ctx, address)
	if err == nil {
		return res, nil
	}

	if err != nil && !errors.Is(err, errNotFound) {
		return Geometry{}, err
	}

	postcode, err := getPostcodeFromAddress(address)
	if err != nil {
		return Geometry{}, err
	}

	res, err = c.get(ctx, postcode)
	if err != nil {
		return Geometry{}, err
	}

	return res, nil

}

var errNotFound = errors.New("not found")

func (c *Client) get(ctx context.Context, address string) (Geometry, error) {
	geometry, err := c.getFromCache(ctx, address)
	if err == nil {
		log.Printf("got from cache: %s - %v", address, geometry)
		return geometry, nil
	}

	// genuine err - not cache miss
	if err != redis.Nil {
		return Geometry{}, err
	}

	params := &maps.FindPlaceFromTextRequest{
		Input:     address,
		InputType: maps.FindPlaceFromTextInputTypeTextQuery,
		Fields: []maps.PlaceSearchFieldMask{
			maps.PlaceSearchFieldMaskGeometryLocation,
		},
	}
	res, err := c.GoogleMapsClient.FindPlaceFromText(ctx, params)
	if err != nil {
		return Geometry{}, err
	}

	if len(res.Candidates) == 0 {
		return Geometry{}, errNotFound
	}

	candidate := res.Candidates[0]
	geometry.Lat = candidate.Geometry.Location.Lat
	geometry.Lng = candidate.Geometry.Location.Lng

	err = c.saveToCache(ctx, address, geometry)
	if err != nil {
		return Geometry{}, err
	}

	return geometry, nil
}
