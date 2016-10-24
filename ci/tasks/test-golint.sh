#!/bin/sh
set -e -x

export GOPATH=$(pwd)/config-server
export PATH=/usr/local/go/bin:$GOPATH/bin:$PATH

cd config-server

scripts/golint.sh
