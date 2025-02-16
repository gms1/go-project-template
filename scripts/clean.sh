#!/usr/bin/bash
DN=$(dirname -- "$0")
BN=$(basename -- "$0")
source "${DN}/common"

usage() {
  cat <<EOT
usage: ${BN} OPTIONS

cleans the project from temporary files

OPTIONS:
  -h|--help  ... display this usage information and exit
EOT
  exit 1
}

OPTS=()
getopt "$@"

[ "${#ARGS[@]}" -eq 0 ] || usage

info "clean..."

rm -rf dist tmp go-project-template

succeeded "clean"
