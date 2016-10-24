#!/bin/bash

errors=$(
  find src/config_server -type f -name '*.go' | xargs -L 1 golint  \
    | grep -v 'should have comment.* or be unexported'      \
    | grep -v '/mocks/'                                     \
    | grep -v '/fakes/'                                     \
    | grep -v 'vendor/'                                     \
    | grep -v 'should not be capitalized'                   \
    | grep -v 'underscore in package name'                  \
    | grep -v 'bootstrapper/spec/'                          \
    | grep -v 'platform/cert/fakes/fake_manager.go'
)

if [ "$(echo -n "$errors")" != "" ]; then
  echo "$errors"
  exit 1
fi