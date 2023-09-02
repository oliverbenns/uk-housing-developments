package scraper

type ScrapeResult struct {
	Builder  string `json:"builder"`
	Name     string `json:"name"`
	Url      string `json:"url"`
	Location string `json:"location"`
}

type Scraper interface {
	Scrape() ([]ScrapeResult, error)
}
