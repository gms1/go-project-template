#!/bin/bash
DN=$(dirname -- "$0")
BN=$(basename -- "$0")
source "${DN}/common"

usage() {
  cat <<EOT
usage: ${BN} OPTIONS [-- [gofumpt arguments...]]

all arguments after '--' will be forwarded to 'gofumpt'

OPTIONS:
  -h|--help  ... display this usage information and exit
EOT
  exit 1
}

OPTS=()
getopt "$@"

[ "${#MAIN_ARGS[@]}" -eq 0 ] || usage

if [ "${#MORE_ARGS[@]}" -ge 1 ]; then
  exec gofumpt -l -w "${MORE_ARGS[@]}"
fi

info "format..."

gofumpt -l -w .

succeeded "format"
