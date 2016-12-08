#!/bin/sh
set -e -x

export GOPATH=$(pwd)/config-server
export PATH=/usr/local/go/bin:$GOPATH/bin:$PATH

semver=`cat version-semver/number`
filename="config-server-${semver}-${GOOS}-${GOARCH}"

cd config-server
go build config_server

mv config_server ../compiled-${GOOS}/${filename}

openssl sha -sha256 ../compiled-${GOOS}/${filename}