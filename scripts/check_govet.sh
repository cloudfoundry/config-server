#!/bin/bash
set -e

errors=$(
  go list ./... \
    | grep -v 'golang.org/' \
    | grep -v 'vendor/' \
    | grep -v 'github.com/' \
    | grep -v 'gopkg.in/'\
    | xargs -L 1 go vet 2>&1 \
    | grep -v 'exit status'
)

[ -z "${errors}" ] && exit 0

echo -e "\ngo vet failed... "
for err in "${errors}"; do
	echo "${err}"
done

echo
exit 1
