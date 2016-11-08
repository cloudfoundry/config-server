#!/bin/bash
set -e

if [ $# -eq 0 ]; then
  gofiles="$(find src/config_server src/integration_tests -type f -name '*.go')"
else
  gofiles="$1"
fi

unformatted_files=$(gofmt -l ${gofiles})
[ -z "$unformatted_files" ] && exit 0

echo -e "\nGo files are not formatted... "
for file in ${unformatted_files}; do
	echo "$file"
done

echo -e "\nYou can run '<project_root_dir>/scripts/gofmt.sh' to format all files."

echo
exit 1
