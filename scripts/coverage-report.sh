#!/bin/bash
DN=$(dirname -- "$0")
BN=$(basename -- "$0")

set -e
set -o pipefail
cd "${DN}/.."

if [ ! -f tmp/coverage.out ]; then
  echo "'tmp/coverage.out' is not available, please run the tests first"
  exit 1
fi
go tool cover -html=tmp/coverage.out
