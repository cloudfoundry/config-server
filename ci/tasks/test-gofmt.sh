#!/bin/sh
set -e -x

export GOPATH=$(pwd)/config-server
export PATH=/usr/local/go/bin:$GOPATH/bin:$PATH

cd config-server/src/config_server

unformatted_files=$(gofmt -l . | wc -l)
if [ ${unformatted_files} -gt 0 ]; then
  echo -e "\n\nGo files are not formatted... "
  gofmt -l .
  exit 1
fi
