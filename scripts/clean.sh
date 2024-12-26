#!/bin/bash
DN=$(dirname -- "$0")
BN=$(basename -- "$0")

echo "clean..."

set -e
set -o pipefail
cd "${DN}/.."

rm -rf dist tmp go-project-template

echo "clean: SUCCEEDED"
