#!/bin/bash
DN=$(dirname -- "$0")
BN=$(basename -- "$0")

echo "test..."

set -e
set -o pipefail
cd "${DN}/.."

LOG=$(mktemp)

exitHandler() {
  trap - "EXIT"
  rm -f "${LOG}"
}

trap exitHandler "EXIT"

mkdir -p tmp
if ! go test ./... -cover -coverprofile=tmp/coverage.out -v &>"${LOG}"; then
  awk '
    function color(c,s) {
            printf("\033[%dm%s\033[0m\n",30+c,s)
    }
    /(^ok|coverage:)\s/ {next}
    /FAIL:/ {color(1,$0);next}
    /PASS:/ {color(2,$0);next}
    {print}
  ' "${LOG}"
  echo "test: FAILED"
  exit 1
fi

echo "test: SUCCEEDED"
