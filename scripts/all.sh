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

./scripts/clean.sh
./scripts/format.sh
./scripts/lint.sh
./scripts/build.sh
./scripts/test.sh
