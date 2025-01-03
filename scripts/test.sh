#!/bin/bash
DN=$(dirname -- "$0")
BN=$(basename -- "$0")

# VERBOSE environment variable
#   when not set, the log from running the tests will not be shown
#   when set, the log output will always shown
#
# additional arguments can be forwarded to 'go test'
#   e.g. for running the tests multiple times (this will take about 30s):
#     VERBOSE=1 ./scripts/test.sh -count 100 -shuffle on
#   e.g. for selecting a specific test you can use:
#     ./scripts/test.sh -run '^TestName$'
#   NOTE: setting "-parallel n" does not seem to have any effect, it looks like we have to
#     call "t.Parallel()" inside the tests too, to mark the tests as parallizable


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

start=$(date "+%s")

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

end=$(date "+%s")
seconds=$((end - start))

if [ "${RC}" -ne 0 ]; then
  [ -t 2 ] && echo -n -e "\e[00;31m" >&2
  printf 'test: FAILED in %d seconds.\n' $((end - start)) >&2
  [ -t 2 ] && echo -n -e '\e[00m' >&2
  exit 1
fi

printf 'test: SUCCEEDED in %d seconds.\n' $((end - start))
