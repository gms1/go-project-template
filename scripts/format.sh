#!/bin/bash
DN=$(dirname -- "$0")
BN=$(basename -- "$0")

if [ $# -ge 1 ]; then
  exec gofumpt -l -w "$@"
fi

echo "format..."

set -e
set -o pipefail
cd "${DN}/.."

gofumpt -l -w .

echo "format: SUCCEEDED"
