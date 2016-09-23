#!/bin/bash

set -e -x

kill -9 `cat config_server.pid`
