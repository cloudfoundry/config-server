#!/usr/bin/env bash
set -eu -o pipefail

export PATH=/usr/local/go/bin:${PATH}

cd config-server
go run github.com/onsi/ginkgo/ginkgo -r -trace -skipPackage="integration,vendor"
