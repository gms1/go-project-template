#!/usr/bin/bash
DN=$(dirname -- "$0")
BN=$(basename -- "$0")
source "${DN}/common"

usage() {
  cat <<EOT
usage: ${BN} OPTIONS

lints the project

OPTIONS:
  -h|--help  ... display this usage information and exit
EOT
  exit 1
}


OPTS=()
getopt "$@"

[ "${#ARGS[@]}" -eq 0 ] || usage

info "lint..."

go mod tidy
golangci-lint run

succeeded "lint"
