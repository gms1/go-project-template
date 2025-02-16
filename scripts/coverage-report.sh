#!/usr/bin/bash
DN=$(dirname -- "$0")
BN=$(basename -- "$0")
source "${DN}/common"

usage() {
  cat <<EOT
usage: ${BN} OPTIONS

opens the coverage report in the browser

OPTIONS:
  -h|--help  ... display this usage information and exit
EOT
  exit 1
}

OPTS=()
getopt "$@"

[ "${#ARGS[@]}" -eq 0 ] || usage

[ -f tmp/coverage.out ] || die "'tmp/coverage.out' is not available, please run the tests first"
go tool cover -html=tmp/coverage.out
