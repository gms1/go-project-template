#!/usr/bin/bash
DN=$(dirname -- "$0")
BN=$(basename -- "$0")
source "${DN}/common"

usage() {
  cat <<EOT
usage: ${BN} OPTIONS

install latest prerequisites

OPTIONS:
  -h|--help  ... display this usage information and exit
EOT
  exit 1
}

OPTS=()
getopt "$@"

[ "${#ARGS[@]}" -eq 0 ] || usage

go install mvdan.cc/gofumpt@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/goreleaser/goreleaser/v2@latest

if [ -n "$(which pip)" ]; then
  pip install pre-commit
fi
