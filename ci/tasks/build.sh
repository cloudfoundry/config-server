#!/usr/bin/env bash
set -eu -o pipefail

export PATH=/usr/local/go/bin:${PATH}

pushd config-server
  go build .
popd

semver=$(cat version-semver/number)
binary_name="config-server-${semver}-${GOOS}-${GOARCH}"
output_filename="compiled-${GOOS}/${binary_name}"

mv config-server/config-server "${output_filename}"

openssl sha256 "${output_filename}"