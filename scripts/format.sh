#!/bin/bash
DN=$(dirname -- "$0")
BN=$(basename -- "$0")
source "${DN}/common"

usage() {
  cat <<EOT
usage: ${BN} OPTIONS [gofumpt arguments...]

all arguments will be forwarded to 'gofumpt'

OPTIONS:
  -h|--help|help  ... display this usage information and exit
EOT
  exit 1
}

[ "$1" != '-h' -a "$1" != '--help' -a "$1" != 'help'  ] || usage

if [ $# -ge 1 ]; then
  exec gofumpt -l -w "$@"
fi

echo "format..."

gofumpt -l -w .

echo "format: SUCCEEDED"
