package scraper

import (
	"fmt"
	"log"

	"github.com/gocolly/colly/v2"
)

type Barratt struct {
}

var _ Scraper = &Barratt{}

func (b *Barratt) Name() string {
	return "Barratt"
}

func (b *Barratt) Scrape() ([]Result, error) {
	c := colly.NewCollector()
	results := []Result{}
	locationPageUrls := []string{}
	baseUrl := "https://www.barratthomes.co.uk"

	c.OnHTML(".location-group > a", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		locationPageUrls = append(locationPageUrls, baseUrl+link)
	})

	listUrl := baseUrl + "/new-homes"
	err := c.Visit(listUrl)
	if err != nil {
		return nil, fmt.Errorf("could not visit %s: %w", listUrl, err)
	}

	for _, pageUrl := range locationPageUrls {
		locationResults, err := b.scrapeLocationPage(pageUrl)
		if err != nil {
			return nil, err
		}
		results = append(results, locationResults...)

	}

	return results, nil
}

func (b *Barratt) scrapeLocationPage(pageUrl string) ([]Result, error) {
	c := colly.NewCollector()
	results := []Result{}

	c.OnHTML(".search-card", func(e *colly.HTMLElement) {
		result := Result{
			Name:     e.ChildText("h2.search-card__heading"),
			Url:      e.ChildAttr("a.search-card__thumbnail", "href"),
			Location: e.ChildText("div.search-card__address"),
		}

		err := result.Validate()
		if err != nil {
			log.Printf("invalid result so omitting %v: %v", result, err)
		} else {
			results = append(results, result)
		}
	})

	err := c.Visit(pageUrl)
	if err != nil {
		return nil, fmt.Errorf("could not visit %s: %w", pageUrl, err)
	}

	return results, nil

}
