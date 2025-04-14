#!/usr/bin/env bash

fly -t "${CONCOURSE_TARGET:-bosh}" \
    set-pipeline -p config-server \
    -c ci/pipeline.yml
