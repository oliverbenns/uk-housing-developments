package location

import (
	"context"
	"errors"

	"googlemaps.github.io/maps"
)

type Client struct {
	GoogleMapsClient *maps.Client
}

type Geometry struct {
	Lat float64
	Lng float64
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
	geometry := Geometry{
		Lat: candidate.Geometry.Location.Lat,
		Lng: candidate.Geometry.Location.Lng,
	}

	return geometry, nil

}
