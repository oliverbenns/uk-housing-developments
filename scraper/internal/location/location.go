package location

import (
	"context"
	"errors"
	"regexp"
	"strings"

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

var errNoPostcode = errors.New("could not get postcode")

func getPostcodeFromAddress(address string) (string, error) {
	parts := strings.Split(address, ",")
	for _, part := range parts {
		val := strings.TrimSpace(part)
		_isPostcode, err := isPostcode(val)
		if err != nil {
			return "", err
		}

		if _isPostcode {
			return val, nil
		}
	}

	return "", errNoPostcode
}

// https://stackoverflow.com/questions/164979/regex-for-matching-uk-postcodes#164994
const postcodePattern = "^([A-Za-z][A-Ha-hJ-Yj-y]?[0-9][A-Za-z0-9]? ?[0-9][A-Za-z]{2}|[Gg][Ii][Rr] ?0[Aa]{2})$"

func isPostcode(val string) (bool, error) {
	return regexp.Match(postcodePattern, []byte(val))
}
