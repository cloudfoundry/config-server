#!/usr/bin/env bash

fly -t bosh-ecosystem set-pipeline -p config-server -c ci/pipeline.yml
