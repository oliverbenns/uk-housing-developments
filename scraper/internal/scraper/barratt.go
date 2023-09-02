package scraper

import (
	"fmt"

	"github.com/gocolly/colly/v2"
)

type Barratt struct {
}

var _ Scraper = &Barratt{}

func (b *Barratt) Scrape() ([]ScrapeResult, error) {
	c := colly.NewCollector()
	results := []ScrapeResult{}
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

func (b *Barratt) scrapeLocationPage(pageUrl string) ([]ScrapeResult, error) {
	c := colly.NewCollector()
	results := []ScrapeResult{}

	c.OnHTML(".search-card", func(e *colly.HTMLElement) {
		result := ScrapeResult{
			Builder:  "Barratt",
			Name:     e.ChildText("h2.search-card__heading"),
			Url:      e.ChildAttr("a.search-card__thumbnail", "href"),
			Location: e.ChildText("div.search-card__address"),
		}

		results = append(results, result)
	})

	err := c.Visit(pageUrl)
	if err != nil {
		return nil, fmt.Errorf("could not visit %s: %w", pageUrl, err)
	}

	return results, nil

}
