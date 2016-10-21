#!/bin/bash

set -e -x

setup_test_dir (){
  rm -rf ./tmp
  mkdir ./tmp

  pushd ./tmp

  cp ../scripts/start_server.sh .
  cp ../scripts/stop_server.sh .
  cp ../scripts/setup_db.sh .
  cp ../scripts/test-integration.sh .
  chmod +x ./*

  cp -r ../src/integration_tests/assets ./assets

  go build -o config_server.bin ../src/config_server/main.go
  ginkgo build -r ../src/integration_tests

  mv ../src/integration_tests/integration_tests.test .
}

run_tests (){
  export DB=$1
  echo -e "\n\nRunning Integration tests for $DB\n\n"
  #Copy DB file to config.json
  cp ./assets/config.$DB.json ./assets/config.json

  ./integration_tests.test
}

setup_test_dir

case $1 in
    memory  )
     run_tests memory ;;
    mysql  )
     run_tests mysql ;;
    postgresql  )
     run_tests postgresql ;;
    * )
     if [ -n "$1" ]; then
        echo "usage test-integration.sh [memory|postgresql|mysql] Defaults to all"
        exit 1
     fi
     run_tests memory
     run_tests mysql
     run_tests postgresql
     ;;
esac



