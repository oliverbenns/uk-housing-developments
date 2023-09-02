package scraper

import (
	"fmt"
	"log"

	"github.com/gocolly/colly/v2"
)

type Persimmon struct {
}

var _ Scraper = &Persimmon{}

func (p *Persimmon) Name() string {
	return "Persimmon"
}

func (p *Persimmon) Scrape() ([]Result, error) {
	c := colly.NewCollector()
	results := []Result{}
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

func (p *Persimmon) scrapeDevelopmentPage(pageUrl string) ([]Result, error) {
	c := colly.NewCollector()
	results := []Result{}

	c.OnHTML("#details", func(e *colly.HTMLElement) {
		result := Result{
			Name:     e.ChildText("h1"),
			Url:      pageUrl,
			Location: e.ChildText("h1 + h2"),
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
