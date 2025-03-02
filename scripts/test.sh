#!/usr/bin/bash
DN=$(dirname -- "$0")
BN=$(basename -- "$0")
source "${DN}/common"

usage() {
  cat <<EOT
usage: ${BN} OPTIONS [-- [go-test arguments...]]

runs the tests for this project

all arguments after '--' will be forwarded to 'go test'
  e.g. for running the tests multiple times (this will take about 30s):
    VERBOSE=1 ./scripts/test.sh -- -count 100 -shuffle on
  e.g. for selecting a specific test you can use:
    ./scripts/test.sh -- -run '^TestName$'
  NOTE: setting "-parallel n" does not seem to have any effect, it looks like we have to
    call "t.Parallel()" inside the tests too, to mark the tests as parallizable

OPTIONS:
  -h|--help  ... display this usage information and exit

environment variables:
  VERBOSE
    when not set, the log output of the tests will not be shown as long as the tests are succeeding
    when set, the log output will be always shown
  OTEL_SDK_DISABLED
    when 'true', otel will be instrumented using a noop exporter
    else otel will use the default otlptracegrpc exporter
    when OTEL_SDK_DISABLED and OTEL_EXPORTER_OTLP_ENDPOINT are not set, OTEL_SDK_DISABLED will be set to "true" before running the tests
  LOG_LEVEL=DEBUG
    sets the logging level to debug
EOT
  exit 1
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


OPTS=()
getopt "$@"

[ "${#MAIN_ARGS[@]}" -eq 0 ] || usage

info "test..."

if [ -z "${OTEL_SDK_DISABLED}" -a -z "${OTEL_EXPORTER_OTLP_ENDPOINT}" ]; then
  export OTEL_SDK_DISABLED="true"
fi


mkdir -p tmp

start=$(date "+%s")

cmd=(go test -cover -coverprofile=tmp/coverage.out -v)
if [ "${#MORE_ARGS[@]}" -eq 0 ]; then
  cmd+=(-count=1 ./...)
else
  cmd+=( "${MORE_ARGS[@]}" )
fi

set +e
echo "\$ ${cmd[@]}"
if [ -n "${VERBOSE}" ]; then
  "${cmd[@]}" 2>&1 | colorize
  RC=$?
else
  createTempFile LOG
  "${cmd[@]}" &>"${LOG}"
  RC=$?
  [ "${RC}" -eq 0 ] || colorize "${LOG}"
fi
set -e

end=$(date "+%s")
seconds=$((end - start))

[ "${RC}" -eq 0 ] || die "test: FAILED in ${seconds} seconds."
./scripts/coverage-report.sh
succeeded "test" "in ${seconds} seconds."
