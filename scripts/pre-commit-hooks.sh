#!/bin/bash

staged_files=$(git diff --cached --name-only --diff-filter=ACMR | grep '.go$')
[ -z "$staged_files" ] && exit 0

success=0

./scripts/check_gofmt.sh "${staged_files}"
[ $? -ne 0 ] && success=1

./scripts/check_golint.sh "${staged_files}"
[ $? -ne 0 ] && success=1

./scripts/check_govet.sh
[ $? -ne 0 ] && success=1

exit ${success}
