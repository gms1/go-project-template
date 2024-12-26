#!/bin/bash
DN=$(dirname -- "$0")
BN=$(basename -- "$0")

set -e
set -o pipefail
cd "${DN}/.."

./scripts/clean.sh
./scripts/format.sh
./scripts/lint.sh
./scripts/build.sh
./scripts/test.sh
