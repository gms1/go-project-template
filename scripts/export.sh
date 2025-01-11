#!/bin/bash
DN=$(dirname -- "$0")
BN=$(basename -- "$0")
source "${DN}/common"

usage() {
  cat <<EOT
usage: ${BN} OPTIONS <target-project-name> [target-project-owner]

OPTIONS:
  -h|--help|help  ... display this usage information and exit

EOT
  exit 1
}


[ "$1" != '-h' -a "$1" != '--help' -a "$1" != 'help'  ] || usage

PROJECT_NAME="$1"
PROJECT_OWNER="$2"
[ -n "${PROJECT_NAME}" -a "$#" -ge 1 -a "$#" -le 2 ] || usage

echo "export..."

PARENT_PREFIX=$(realpath ..)

rm -rf "${PARENT_PREFIX}/${PROJECT_NAME}"
git checkout-index --prefix="${PARENT_PREFIX}/${PROJECT_NAME}/" -a

cd "${PARENT_PREFIX}/${PROJECT_NAME}/"

rm -rf "docs"
mkdir -p "docs"

sed -i "s|\(Version\s*=\s*\"\)[^\"]*\"|\10.0.1\"|" pkg/common/about.go

find . -type f -exec sed -i "s|go-project-template|${PROJECT_NAME}|g" {} \;

if [ -n "${PROJECT_OWNER}" ]; then
  find . -type f -exec sed -i "s|gms1|${PROJECT_OWNER}|g" {} \;
fi

cd -

echo "export: SUCCEEDED to '${PARENT_PREFIX}/${PROJECT_NAME}'"
