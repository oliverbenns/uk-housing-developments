package scraper

import (
	"fmt"
	"log"

	"github.com/gocolly/colly/v2"
)

type Berkeley struct {
}

var _ Scraper = &Berkeley{}

func (b *Berkeley) Scrape() ([]ScrapeResult, error) {
	c := colly.NewCollector()
	results := []ScrapeResult{}
	countyPageUrls := []string{}
	baseUrl := "https://www.berkeleygroup.co.uk"

	c.OnHTML("#mainNav > li:first-child .menu-second-level--navigation .menu-third-level--wrap a", func(e *colly.HTMLElement) {

		link := e.Attr("href")
		countyPageUrls = append(countyPageUrls, baseUrl+link)
	})

	err := c.Visit(baseUrl)
	if err != nil {
		return nil, fmt.Errorf("could not visit %s: %w", baseUrl, err)
	}
	log.Print("a", countyPageUrls)

	for _, pageUrl := range countyPageUrls {
		locationResults, err := b.scrapeCountyPage(pageUrl)
		if err != nil {
			return nil, err
		}
		results = append(results, locationResults...)

	}

	return results, nil
}

func (b *Berkeley) scrapeCountyPage(pageUrl string) ([]ScrapeResult, error) {
	c := colly.NewCollector()
	results := []ScrapeResult{}

	c.OnHTML(".search-card", func(e *colly.HTMLElement) {
		result := ScrapeResult{
			Builder:  "Berkeley",
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
