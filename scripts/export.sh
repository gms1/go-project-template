#!/usr/bin/bash
DN=$(dirname -- "$0")
BN=$(basename -- "$0")
source "${DN}/common"

usage() {
  cat <<EOT
usage: ${BN} OPTIONS <target-project-name> [target-project-owner]

exports this project

OPTIONS:
  -h|--help  ... display this usage information and exit

EOT
  exit 1
}

OPTS=()
getopt "$@"

[ "${#ARGS[@]}" -ge 1 -a "${#ARGS[@]}" -le 2 ] || usage

PROJECT_NAME="${ARGS[0]}"
PROJECT_OWNER="${ARGS[1]}"
[ -n "${PROJECT_NAME}" ] || usage

info "export..."

PARENT_PREFIX=$(realpath ..)

rm -rf "${PARENT_PREFIX}/${PROJECT_NAME}"
git checkout-index --prefix="${PARENT_PREFIX}/${PROJECT_NAME}/" -a

cd "${PARENT_PREFIX}/${PROJECT_NAME}/"

rm -rf "docs"
mkdir -p "docs"

sed -i "s|\(Version\s*=\s*\"\)[^\"]*\"|\10.0.1\"|" "${ABOOUT_GO}"

find . -type f -exec sed -i "s|go-project-template|${PROJECT_NAME}|g" {} \;

if [ -n "${PROJECT_OWNER}" ]; then
  find . -type f -exec sed -i "s|gms1|${PROJECT_OWNER}|g" {} \;
fi

cd -

succeeded "export" "to export to '${PARENT_PREFIX}/${PROJECT_NAME}'"
