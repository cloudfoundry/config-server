#!/bin/bash
set -e

go clean -r config_server

echo -e "\n Formatting Go files..."
gofmt -l -w src/config_server/

echo -e "\n Testing package..."
ginkgo -r --keepGoing src/config_server/
