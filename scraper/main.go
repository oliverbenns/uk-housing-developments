package main

import (
	"log"

	"github.com/oliverbenns/uk-housing-developments/scraper/internal/service"
)

func main() {
	svc := service.Service{}
	err := svc.Run()
	if err != nil {
		log.Fatalf("could not run service: %v", err)
	}
}
