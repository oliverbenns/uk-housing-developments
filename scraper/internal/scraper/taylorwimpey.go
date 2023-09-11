package scraper

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/oliverbenns/uk-housing-developments/scraper/internal/cloudflare"
)

type TaylorWimpey struct {
}

var _ Scraper = &TaylorWimpey{}

func (tw *TaylorWimpey) Name() string {
	return "Taylor Wimpey"
}

func (tw *TaylorWimpey) Scrape() ([]Result, error) {
	results := []Result{}
	baseUrl := "https://www.taylorwimpey.co.uk"

	developments, err := tw.requestDevelopments(baseUrl + "/api/tw/DevelopmentSearch/GetDevelopments")
	if err != nil {
		return nil, err
	}

	for _, development := range developments {
		locationResult, err := tw.scrapeDevelopmentPageWithRetry(baseUrl, development, 5)
		if err != nil {
			return nil, err
		}

		results = append(results, locationResult)
	}

	return results, nil
}

var errParseSubHeading = errors.New("could not parse subheading")

// Due to terrible website timeouts.
func (tw *TaylorWimpey) scrapeDevelopmentPageWithRetry(baseUrl string, development TaylorWimpeyAPIDevelopment, retryCount int) (Result, error) {
	var result Result
	var err error
	for i := 0; i < retryCount; i++ {
		result, err = tw.scrapeDevelopmentPage(baseUrl, development)
		if err != nil {
			log.Printf("could not get development page %s, retries left: %d", development.Url, retryCount-1-i)
			numSecs := time.Duration(5 * (i + 1))
			time.Sleep(numSecs * time.Second)
			continue
		}

		return result, nil
	}

	return Result{}, err
}

// We actually get most of the info from the API.
// But we do not have the location so use the page src.
func (tw *TaylorWimpey) scrapeDevelopmentPage(baseUrl string, development TaylorWimpeyAPIDevelopment) (Result, error) {
	c := colly.NewCollector()
	cloudflareTransport := cloudflare.NewTransport()
	c.WithTransport(&cloudflareTransport)
	// Some development pages are slow. Or Cloudflare throttle?
	c.SetRequestTimeout(20 * time.Second)

	pageUrl := baseUrl + development.Url
	result := Result{}
	redirectUrl := ""

	c.OnResponse(func(r *colly.Response) {
		requestUrl := r.Request.URL.String()
		if requestUrl != pageUrl {
			redirectUrl = requestUrl
		}
	})

	var parseErr error
	c.OnHTML("html", func(e *colly.HTMLElement) {
		subHeading := e.ChildText(".landing-page-hero__standfirst")
		subHeadingParts := strings.Split(subHeading, "Prices from")
		if len(subHeadingParts) == 0 {
			parseErr = errParseSubHeading
			return
		}

		result = Result{
			Name:     e.ChildText("h1.landing-page-hero__title"),
			Url:      pageUrl,
			Location: subHeadingParts[0],
		}
	})

	err := c.Visit(pageUrl)
	if err != nil {
		return Result{}, fmt.Errorf("could not visit %s: %w", pageUrl, err)
	}

	// Some developments have their own website.
	// We can't scrape all so take a best guess with just the API data we have.
	if redirectUrl != "" {
		result = Result{
			Name: development.Name,
			Url:  redirectUrl,
			// Probably going to fail location lookup.
			Location: development.Name,
		}
	}

	// we know it's not redirect so page should be parsable
	if parseErr != nil {
		return Result{}, parseErr
	}

	err = result.Validate()
	if err != nil {
		log.Printf("invalid result so omitting %v: %v", result, err)
		return Result{}, err
	}

	return result, nil
}

type TaylorWimpeyAPIDevelopmentResponse struct {
	Results []TaylorWimpeyAPIDevelopment `json:"Results"`
	Status  string                       `json:"Status"`
}

type TaylorWimpeyAPIDevelopment struct {
	Name string `json:"Name"`
	Type int    `json:"Type"`
	Url  string `json:"Url"`
}

var errInvalidResponseCode = errors.New("invalid response code")

func (tw *TaylorWimpey) requestDevelopments(apiUrl string) ([]TaylorWimpeyAPIDevelopment, error) {
	client := cloudflare.CreateHttpClient()

	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("could not get %s, code %d: %w", apiUrl, resp.StatusCode, errInvalidResponseCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var res TaylorWimpeyAPIDevelopmentResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	return res.Results, nil
}
