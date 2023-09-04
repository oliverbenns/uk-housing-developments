#!/bin/bash

pushd scraper
env $(cat .env | xargs) go run . > ../site/public/developments.json
popd
