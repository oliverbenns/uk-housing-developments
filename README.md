# UK Housing Developments

App that scrapes and maps locations of housing developments from some of the biggest UK developers.

Site was bootstrapped with [Create React App](https://github.com/facebook/create-react-app).

View the project here: [oliverbenns.github.io/uk-housing-developments](https://oliverbenns.github.io/uk-housing-developments)

## Getting started

### Site

-   `cp .env.example.local .env.local`
-   Add secrets to the newly copied env file
-   `npm start`

### Scraper

-   `docker compose up -d`
-   `cp .env.example .env`
-   Run the shell script `sh scrape.sh`
