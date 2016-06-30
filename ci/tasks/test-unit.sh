#!/bin/sh
set -e -x

export GOPATH=$(pwd)/config-server
export PATH=/usr/local/go/bin:$GOPATH/bin:$PATH

go clean -r config_server

go get github.com/onsi/ginkgo/ginkgo
go get github.com/onsi/gomega

cd config-server/src/config_server
ginkgo -r -trace 