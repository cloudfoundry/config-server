#!/bin/bash
set -e

staged_files=$(git diff --cached --name-only --diff-filter=ACMR | grep '.go$')
[ -z "$staged_files" ] && exit 0

unformatted_files=$(gofmt -l ${staged_files})
[ -z "$unformatted_files" ] && exit 0

echo -e "\n Go files are not formatted... "
for file in ${unformatted_files}; do
	echo "$file"
done

echo
exit 1
