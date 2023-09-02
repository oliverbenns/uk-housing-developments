package scraper

type ScrapeResult struct {
	Name     string
	Url      string
	Location string
}

type Scraper interface {
	Name() string
	Scrape() ([]ScrapeResult, error)
}
