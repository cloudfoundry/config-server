#!/usr/bin/env bash
set -eu -o pipefail

cd config-server

bin/lint
