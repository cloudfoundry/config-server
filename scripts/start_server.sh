#!/bin/bash

ASSETS_DIR="./assets"
CONFIG_FILE="$ASSETS_DIR/config.json"

cat > $CONFIG_FILE <<- EOM
{
  "port": 9000,
  "store": "memory",
  "certificate_file_path": "$ASSETS_DIR/ssl.crt",
  "private_key_file_path": "$ASSETS_DIR/ssl.key",
  "ca_certificate_file_path": "$ASSETS_DIR/ca.crt",
  "ca_private_key_file_path": "$ASSETS_DIR/crt.key",
  "jwt_verification_key_path": "$ASSETS_DIR/uaa.pub"
}
EOM

./config_server.bin $CONFIG_FILE &
SERVER_PID=$!

echo "$SERVER_PID" > config_server.pid