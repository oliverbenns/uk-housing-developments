package main

import (
	"fmt"
	"log"

	"github.com/oliverbenns/uk-housing-developments/scraper/internal/service"
)

func main() {
	svc := service.Service{}
	data, err := svc.Run()
	if err != nil {
		log.Fatalf("could not run service: %v", err)
	}

	fmt.Println(string(data))
}
