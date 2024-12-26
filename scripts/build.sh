#!/bin/bash
DN=$(dirname -- "$0")
BN=$(basename -- "$0")

echo "build..."

set -e
set -o pipefail
cd "${DN}/.."

go generate
go build -v

echo "build: SUCCEEDED"
