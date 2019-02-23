#!/bin/bash

export OCTOMON_OCTOPUS_ADDRESS=https://demo.octopusdeploy.com
export OCTOMON_OCTOPUS_APIKEY=API-GUEST
export OCTOMON_HEALTHCHECK_INTERVAL=10s
go run main.go

