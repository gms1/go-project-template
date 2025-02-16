#!/usr/bin/bash
DN=$(dirname -- "$0")
BN=$(basename -- "$0")
source "${DN}/common"

usage() {
  cat <<EOT
usage: ${BN} OPTIONS

runs: clean, format, lint, build and test

OPTIONS:
  -h|--help  ... display this usage information and exit
EOT
  exit 1
}

OPTS=()
getopt "$@"

[ "${#ARGS[@]}" -eq 0 ] || usage

./scripts/clean.sh
./scripts/format.sh
./scripts/lint.sh
./scripts/build.sh
./scripts/test.sh
