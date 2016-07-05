#!/bin/sh
set -e -x

export GOPATH=$(pwd)/config-server
export PATH=/usr/local/go/bin:$GOPATH/bin:$PATH

go clean -r config_server

echo INTEGRATION_TESTS_NOT_IMPLEMENTED!