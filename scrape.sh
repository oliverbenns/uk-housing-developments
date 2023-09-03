#!/bin/bash

pushd scraper
go run . > ../site/public/developments.json
popd
