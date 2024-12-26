#!/bin/bash
DN=$(dirname -- "$0")
BN=$(basename -- "$0")

echo "test..."

set -e
set -o pipefail
cd "${DN}/.."

mkdir -p tmp
go test ./... "$@" -cover -coverprofile=tmp/coverage.out

echo "test: SUCCEEDED"
