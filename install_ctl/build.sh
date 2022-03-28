#!/bin/bash
go mod tidy
go build -o cessctl ../cessctl/main.go