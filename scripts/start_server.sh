#!/bin/bash

set -e -x

ASSETS_DIR="./assets"
CONFIG_FILE="$ASSETS_DIR/config.json"

./config_server.bin $CONFIG_FILE &
SERVER_PID=$!

echo "$SERVER_PID" > config_server.pid