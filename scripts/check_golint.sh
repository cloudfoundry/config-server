#!/bin/bash
set -e

if [ $# -eq 0 ]; then
  gofiles="$(find src/config_server src/integration_tests -type f -name '*.go')"
else
  gofiles="$1"
fi

errors=$(
  echo "${gofiles}" | xargs -L 1 golint \
    | grep -v 'vendor/' \
    | grep -v '/mocks/' \
    | grep -v '/*fakes/' \
    | grep -v 'should have comment.* or be unexported' \
    | grep -v 'should not be capitalized' \
    | grep -v 'underscore in package name' \
    | xargs # If grep returns no results, it exits with exit code 1
)

[ -z "${errors}" ] && exit 0

echo -e "\ngolint failed with errors... "
for err in "${errors}"; do
	echo "${err}"
done

echo
exit 1
