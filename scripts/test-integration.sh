#!/bin/bash

set -e -x

rm -rf ./tmp
mkdir ./tmp

pushd ./tmp

cp ../scripts/start_server.sh .
cp ../scripts/stop_server.sh .
cp ../scripts/test-integration.sh .
chmod +x ./*

cp -r ../src/integration_tests/assets ./assets

go build -o config_server.bin ../src/config_server/main.go
ginkgo build -r ../src/integration_tests

mv ../src/integration_tests/integration_tests.test .

./integration_tests.test
