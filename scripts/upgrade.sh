#!/usr/bin/bash
DN=$(dirname -- "$0")
BN=$(basename -- "$0")
source "${DN}/common"

usage() {
  cat <<EOT
usage: ${BN} OPTIONS

upgrade modules and prerequisites

OPTIONS:
  -h|--help  ... display this usage information and exit
EOT
  exit 1
}

OPTS=()
getopt "$@"

[ "${#ARGS[@]}" -eq 0 ] || usage

upgrade_flag="-u"
[ "${UPGRADE_MODE}" != "patch" ] || upgrade_flag="-u=patch"

go get ${upgrade_flag} -t ./...
go mod tidy
