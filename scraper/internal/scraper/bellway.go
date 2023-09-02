package scraper

import (
	"fmt"
	"log"

	"github.com/gocolly/colly/v2"
)

type Bellway struct {
}

var _ Scraper = &Bellway{}

func (b *Bellway) Name() string {
	return "Bellway"
}

func (b *Bellway) Scrape() ([]Result, error) {
	c := colly.NewCollector()
	results := []Result{}
	locationPageUrls := []string{}
	baseUrl := "https://www.bellway.co.uk"

	c.OnHTML("a.map-point", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		locationPageUrls = append(locationPageUrls, baseUrl+link)
	})

	listUrl := baseUrl + "/new-homes"
	err := c.Visit(listUrl)
	if err != nil {
		return nil, fmt.Errorf("could not visit %s: %w", listUrl, err)
	}

	for _, pageUrl := range locationPageUrls {
		locationResults, err := b.scrapeLocationPage(baseUrl, pageUrl)
		if err != nil {
			return nil, err
		}
		results = append(results, locationResults...)

	}

	return results, nil
}

func (b *Bellway) scrapeLocationPage(baseUrl, pageUrl string) ([]Result, error) {
	c := colly.NewCollector()
	results := []Result{}

	c.OnHTML(".search__results__list .tile", func(e *colly.HTMLElement) {
		result := Result{
			Name:     e.ChildText(".heading"),
			Url:      baseUrl + e.ChildAttr(".tile__content > a", "href"),
			Location: e.ChildText(".description"),
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
