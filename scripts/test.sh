#!/bin/bash
DN=$(dirname -- "$0")
BN=$(basename -- "$0")

usage() {
  cat <<EOT
usage: ${BN} OPTIONS [go-test arguments...]

all arguments will be forwarded to 'go test'
  e.g. for running the tests multiple times (this will take about 30s):
    VERBOSE=1 ./scripts/test.sh -count 100 -shuffle on
  e.g. for selecting a specific test you can use:
    ./scripts/test.sh -run '^TestName$'
  NOTE: setting "-parallel n" does not seem to have any effect, it looks like we have to
    call "t.Parallel()" inside the tests too, to mark the tests as parallizable

OPTIONS:
  -h|--help|help  ... display this usage information and exit

environment variables:
  VERBOSE
    when not set, the log output of the tests will not be shown as long as the tests are succeeding
    when set, the log output will be always shown
  OTEL_SDK_DISABLED
    when 'true', otel will be instrumented using a noop exporter
    else otel will use the default otlptracegrpc exporter
    when OTEL_SDK_DISABLED and OTEL_EXPORTER_OTLP_ENDPOINT are not set, OTEL_SDK_DISABLED will be set to "true" before running the tests
EOT
  exit 1
}

[ "$1" != '-h' -a "$1" != '--help' -a "$1" != 'help'  ] || usage


echo "test..."

if [ -z "${OTEL_SDK_DISABLED}" -a -z "${OTEL_EXPORTER_OTLP_ENDPOINT}" ]; then
  export OTEL_SDK_DISABLED="true"
fi

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

[ "${RC}" -eq 0 ] || die "test: FAILED in ${seconds}."

echo "test: COVERAGE:"
go tool cover -func=./tmp/coverage.out | sed '/100\.0/d; s/^/  /'
echo "test: SUCCEEDED in ${seconds} seconds."
