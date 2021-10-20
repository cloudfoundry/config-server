#!/bin/bash

set -e

if which certstrap ; then
  folder=$(mktemp -d)

  certstrap --depot-path "$folder" init --cn config-server --expires "5 years" --passphrase "" > /dev/null 2>&1
  certstrap --depot-path "$folder" request-cert --cn integration --ip 127.0.0.1 --domain localhost --passphrase "" > /dev/null 2>&1
  certstrap --depot-path "$folder" sign --expires "5 years" --CA config-server integration > /dev/null 2>&1
  cp "$folder"/config-server.crt integration/assets/ssl_root_ca.crt
  cp "$folder"/integration.crt integration/assets/ssl.crt
  cp "$folder"/integration.key integration/assets/ssl.key
  echo "Done! Your certs have been regenerated"
else
  echo "This requires certstrap, install it from here: https://github.com/square/certstrap"
fi