#!/bin/bash
DN=$(dirname -- "$0")
BN=$(basename -- "$0")

# VERBOSE environment variable
#   when not set, the log from running the tests will not be shown
#   when set,
# additional arguments will be forwarded to 'go test'
#   e.g. for running the tests multiple times (this will take about 5 min)
#     VERBOSE=1 ./scripts/test.sh -count 100 -shuffle on -parallel 4
#   e.g. for selecting a specific test you can use
#     ./scripts/test.sh -run '^<TestName>$'

echo "test..."

set -e
set -o pipefail
cd "${DN}/.."

exitHandler() {
  trap - "EXIT"
  rm -f "${LOG}"
}

colorize() {
  awk '
    function color(c,s) {
            printf("\033[%dm%s\033[0m\n",30+c,s)
    }
    /(^ok|coverage:)\s/ {next}
    /FAIL:/ {color(1,$0);next}
    /PASS:/ {color(2,$0);next}
    {print}
  '  "$@"
}

trap exitHandler "EXIT"

mkdir -p tmp

if [ -n "${VERBOSE}" ]; then
  set +e
  go test ./... -cover -coverprofile=tmp/coverage.out -v "$@" 2>&1 | colorize
  RC=$?
  set -e
else
  LOG=$(mktemp)
  set +e
  go test ./... -cover -coverprofile=tmp/coverage.out -v "$@" &>"${LOG}"
  RC=$?
  set -e
  [ "${RC}" -eq 0 ] || colorize "${LOG}"
fi

if [ "${RC}" -ne 0 ]; then
  [ -t 2 ] && echo -n -e "\e[00;31m" >&2
  echo "test: FAILED" >&2
  [ -t 2 ] && echo -n -e '\e[00m' >&2
  exit 1
fi

echo "test: SUCCEEDED"
