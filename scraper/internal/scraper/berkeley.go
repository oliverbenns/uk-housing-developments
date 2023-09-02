package scraper

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/gocolly/colly/v2"
)

type Berkeley struct {
}

var _ Scraper = &Berkeley{}

func (b *Berkeley) Scrape() ([]ScrapeResult, error) {
	c := colly.NewCollector()
	results := []ScrapeResult{}
	locationPageUrls := []string{}
	baseUrl := "https://www.berkeleygroup.co.uk"

	c.OnHTML("#mainNav > li:first-child .menu-second-level--navigation .menu-third-level--wrap a", func(e *colly.HTMLElement) {

		link := e.Attr("href")
		locationPageUrls = append(locationPageUrls, baseUrl+link)
	})

	err := c.Visit(baseUrl)
	if err != nil {
		return nil, fmt.Errorf("could not visit %s: %w", baseUrl, err)
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

func (b *Berkeley) scrapeLocationPage(baseUrl, pageUrl string) ([]ScrapeResult, error) {
	// Site uses ajax loading for their developments so colly not suitable.
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	results := []ScrapeResult{}

	// Super hacky way of loading all results.
	// Should update this to poll for new elements, check if last list item is still last list item, etc.
	scrollDown := func(ctx context.Context) error {
		for i := 0; i < 5; i++ {
			// wait for network + render
			time.Sleep(2 * time.Second)

			_, exp, err := runtime.Evaluate(`window.scrollTo(0,document.body.scrollHeight);`).Do(ctx)
			if err != nil {
				return err
			}
			if exp != nil {
				return exp
			}
		}

		return nil
	}

	var html string

	chromedp.Run(ctx,
		chromedp.Navigate(pageUrl),
		chromedp.ActionFunc(scrollDown),
		chromedp.OuterHTML("html", &html, chromedp.ByQuery),
	)

	reader := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}

	doc.Find(".result-wrapper").Each(func(i int, s *goquery.Selection) {
		result := ScrapeResult{
			Builder:  "Berkeley",
			Name:     s.Find("h2").Text(),
			Url:      baseUrl + s.Find(".button--primary").AttrOr("href", ""),
			Location: s.Find(".address").Text(),
		}

		results = append(results, result)
	})

	return results, nil

}
