package scraper

import (
	"fmt"

	"github.com/gocolly/colly/v2"
)

type Bellway struct {
}

var _ Scraper = &Bellway{}

func (b *Bellway) Name() string {
	return "Bellway"
}

func (b *Bellway) Scrape() ([]ScrapeResult, error) {
	c := colly.NewCollector()
	results := []ScrapeResult{}
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
		locationResults, err := b.scrapeLocationPage(pageUrl)
		if err != nil {
			return nil, err
		}
		results = append(results, locationResults...)

	}

	return results, nil
}

func (b *Bellway) scrapeLocationPage(pageUrl string) ([]ScrapeResult, error) {
	c := colly.NewCollector()
	results := []ScrapeResult{}

	c.OnHTML(".search__results__list .tile", func(e *colly.HTMLElement) {
		result := ScrapeResult{
			Name:     e.ChildText(".heading"),
			Url:      e.ChildAttr("tile_content > a", "href"),
			Location: e.ChildText(".description"),
		}

		results = append(results, result)
	})

	err := c.Visit(pageUrl)
	if err != nil {
		return nil, fmt.Errorf("could not visit %s: %w", pageUrl, err)
	}

	return results, nil

}
