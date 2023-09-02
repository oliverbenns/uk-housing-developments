package scraper

import (
	"fmt"
	"log"

	"github.com/gocolly/colly/v2"
)

type TaylorWimpey struct {
}

var _ Scraper = &TaylorWimpey{}

func (tw *TaylorWimpey) Name() string {
	return "Taylor Wimpey"
}

func (tw *TaylorWimpey) Scrape() ([]Result, error) {
	c := colly.NewCollector()
	results := []Result{}
	locationPageUrls := []string{}
	baseUrl := "https://www.taylorwimpey.co.uk"

	c.OnHTML("a.map-point", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		locationPageUrls = append(locationPageUrls, baseUrl+link)
	})

	sitemapUrl := baseUrl + "/sitemap"
	err := c.Visit(sitemapUrl)
	if err != nil {
		return nil, fmt.Errorf("could not visit %s: %w", sitemapUrl, err)
	}

	// Sitemap does not distinguish by location type, so easy to get dupes.
	// E.g. loading up Norwich page and Norfolk Page.
	// Use development url as unique id.
	resultsByUrl := map[string]Result{}
	for _, pageUrl := range locationPageUrls {
		locationResults, err := tw.scrapeLocationPage(pageUrl)
		if err != nil {
			return nil, err
		}

		for _, locationResult := range locationResults {
			_, ok := resultsByUrl[locationResult.Url]
			if !ok {
				resultsByUrl[locationResult.Url] = locationResult
				continue
			} else {
				log.Print("dupe found!", locationResult)
			}
		}
	}

	for _, result := range resultsByUrl {
		results = append(results, result)
	}

	return results, nil
}

func (tw *TaylorWimpey) scrapeLocationPage(pageUrl string) ([]Result, error) {
	c := colly.NewCollector()
	results := []Result{}

	c.OnHTML(".hf-dev-segment-content", func(e *colly.HTMLElement) {
		result := Result{
			Name:     e.ChildText(".hf-dev-segment-content__heading-title a"),
			Url:      e.ChildAttr(".hf-dev-segment-content__heading-title a", "href"),
			Location: e.ChildText(".hf-dev-segment-content__location--address"),
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
