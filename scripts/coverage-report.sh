#!/bin/bash
DN=$(dirname -- "$0")
BN=$(basename -- "$0")
source "${DN}/common"

usage() {
  cat <<EOT
usage: ${BN} OPTIONS

OPTIONS:
  -h|--help|help  ... display this usage information and exit
EOT
  exit 1
}

[ "$1" != '-h' -a "$1" != '--help' -a "$1" != 'help'  ] || usage

[ -f tmp/coverage.out ] || die "'tmp/coverage.out' is not available, please run the tests first"
go tool cover -html=tmp/coverage.out
