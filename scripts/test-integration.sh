#!/bin/bash

go build -o ./test_server src/config_server/main.go

ASSETS_DIR="$PWD/src/integration_tests/assets"
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

./test_server $CONFIG_FILE &
SERVER_PID=$!

ginkgo build -r src/integration_tests

# prevent ginkgo from running tests in random tmp directory
mv src/integration_tests/integration_tests.test .

./integration_tests.test

kill $SERVER_PID

rm ./test_server
rm ./integration_tests.test
rm $CONFIG_FILE