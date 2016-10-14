#!/bin/bash

set -e -x

if [ -z "$1" ]; then
 echo "usage test-integration.sh [memory|postgresql|mysql]"
 exit 1
fi

export DB=$1

rm -rf ./tmp
mkdir ./tmp

pushd ./tmp

cp ../scripts/start_server.sh .
cp ../scripts/stop_server.sh .
cp ../scripts/setup_db.sh .
cp ../scripts/test-integration.sh .
chmod +x ./*

cp -r ../src/integration_tests/assets ./assets

#Copy DB file to config.json
cp ./assets/config.$DB.json ./assets/config.json

go build -o config_server.bin ../src/config_server/main.go
ginkgo build -r ../src/integration_tests

mv ../src/integration_tests/integration_tests.test .

./integration_tests.test
