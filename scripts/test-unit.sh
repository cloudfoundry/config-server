#!/bin/bash
set -e
BASEDIR=$(dirname "$0")

go clean -r config_server

echo -e "\n Generating Fakes..."
${BASEDIR}/generate_fakes.sh

echo -e "\n Formatting Go files..."
gofmt -l -w src/config_server/

echo -e "\n Testing package..."
ginkgo -r --keepGoing src/config_server/
