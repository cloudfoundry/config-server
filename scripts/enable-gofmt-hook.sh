#!/bin/sh
set -e

ln -s -f ../../scripts/gofmt-git-hook.sh .git/hooks/pre-commit
