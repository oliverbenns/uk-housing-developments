package scraper

import (
	"errors"
	"net/url"
)

type Scraper interface {
	Name() string
	Scrape() ([]Result, error)
}

type Result struct {
	Name     string
	Url      string
	Location string
}

var errInvalidName = errors.New("invalid name")
var errInvalidUrl = errors.New("invalid url")
var errInvalidLocation = errors.New("invalid location")

func (r *Result) Validate() error {
	if r.Name == "" {
		return errInvalidName
	}

	if r.Url == "" {
		return errInvalidUrl
	}

	_, err := url.Parse(r.Url)
	if err != nil {
		return errInvalidUrl
	}

	if r.Location == "" {
		return errInvalidLocation
	}

	return nil
}
