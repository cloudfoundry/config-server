#!/bin/sh
set -e -x

export GOPATH=$(pwd)
export PATH=/usr/local/go/bin:$GOPATH/bin:$PATH

cd src/github.com/cloudfoundry/config-server
go run github.com/onsi/ginkgo/ginkgo -r -trace -skipPackage="integration,vendor"
