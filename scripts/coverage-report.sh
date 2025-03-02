#!/usr/bin/bash
DN=$(dirname -- "$0")
BN=$(basename -- "$0")
source "${DN}/common"

usage() {
  cat <<EOT
usage: ${BN} OPTIONS

generates the coverage report and optional opens it in the browser

OPTIONS:
  -h|--help  ... display this usage information and exit
  --open     ... do not open in browser
EOT
  exit 1
}

OPTS=(open)
getopt "$@"

[ "${#ARGS[@]}" -eq 0 ] || usage

[ -f tmp/coverage.out ] || die "'tmp/coverage.out' is not available, please run the tests first"
info "test: COVERAGE:"
go tool cover -func=./tmp/coverage.out | awk '$3 !~ /100\.0%/{ sub(/:$/,"",$1); print "  " $3 " " $1 }'

go tool cover -html=tmp/coverage.out -o tmp/coverage.html
[ -z "${open}" ] || xdg-open tmp/coverage.html
