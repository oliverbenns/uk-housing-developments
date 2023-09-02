package scraper

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/gocolly/colly/v2"
)

type Berkeley struct {
}

var _ Scraper = &Berkeley{}

func (b *Berkeley) Name() string {
	return "Berkeley"
}

func (b *Berkeley) Scrape() ([]Result, error) {
	c := colly.NewCollector()
	results := []Result{}
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

	// Scraper is slow due to ajax scroll functionality and
	// loading in browser so run concurrently.
	mu := sync.Mutex{}
	var wg sync.WaitGroup
	wg.Add(len(locationPageUrls))
	errs := []error{}

	for _, pageUrl := range locationPageUrls {
		go func(closurePageUrl string) {
			locationResults, err := b.scrapeLocationPage(baseUrl, closurePageUrl)
			mu.Lock()
			defer mu.Unlock()

			if err != nil {
				errs = append(errs, err)
			} else {
				results = append(results, locationResults...)
			}
			wg.Done()
		}(pageUrl)
	}

	wg.Wait()

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return results, nil
}

func (b *Berkeley) scrapeLocationPage(baseUrl, pageUrl string) ([]Result, error) {
	// Site uses ajax loading for their developments so colly not suitable.
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	results := []Result{}

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

	err := chromedp.Run(ctx,
		chromedp.Navigate(pageUrl),
		chromedp.ActionFunc(scrollDown),
		chromedp.OuterHTML("html", &html, chromedp.ByQuery),
	)
	if err != nil {
		return nil, err
	}

	reader := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}

	doc.Find(".result-wrapper").Each(func(i int, s *goquery.Selection) {
		result := Result{
			Name:     s.Find("h2").Text(),
			Url:      baseUrl + s.Find(".button--primary").AttrOr("href", ""),
			Location: s.Find(".address").Text(),
		}

		err := result.Validate()
		if err != nil {
			log.Printf("invalid result so omitting %v: %v", result, err)
			return
		}

		results = append(results, result)
	})

	return results, nil

}
