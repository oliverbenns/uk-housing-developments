package scraper

import (
	"fmt"

	"github.com/gocolly/colly/v2"
)

type Persimmon struct {
}

var _ Scraper = &Persimmon{}

func (p *Persimmon) Scrape() ([]ScrapeResult, error) {
	c := colly.NewCollector()
	results := []ScrapeResult{}
	developmentPageUrls := []string{}
	baseUrl := "https://www.persimmonhomes.com"

	c.OnHTML(".region-list a", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		developmentPageUrls = append(developmentPageUrls, baseUrl+link)
	})

	listUrl := baseUrl + "/find-your-new-home"
	err := c.Visit(listUrl)
	if err != nil {
		return nil, fmt.Errorf("could not visit %s: %w", listUrl, err)
	}

	for _, pageUrl := range developmentPageUrls {
		locationResults, err := p.scrapeDevelopmentPage(pageUrl)
		if err != nil {
			return nil, err
		}
		results = append(results, locationResults...)

	}

	return results, nil
}

func (p *Persimmon) scrapeDevelopmentPage(pageUrl string) ([]ScrapeResult, error) {
	c := colly.NewCollector()
	results := []ScrapeResult{}

	c.OnHTML("#details", func(e *colly.HTMLElement) {
		result := ScrapeResult{
			Builder:  "Persimmon",
			Name:     e.ChildText("h1"),
			Url:      pageUrl,
			Location: e.ChildText("h1 + h2"),
		}

		results = append(results, result)
	})

	err := c.Visit(pageUrl)
	if err != nil {
		return nil, fmt.Errorf("could not visit %s: %w", pageUrl, err)
	}

	return results, nil

}
