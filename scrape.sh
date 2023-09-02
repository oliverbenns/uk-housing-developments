#!/bin/bash

pushd scraper
go run . > ../data.json
popd
