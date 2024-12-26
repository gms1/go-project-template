#!/bin/bash
DN=$(dirname -- "$0")
BN=$(basename -- "$0")

echo "export..."

set -e
set -o pipefail
cd "${DN}/.."

PROJECT_NAME="$1"
PROJECT_OWNER="$2"
if [ -z "${PROJECT_NAME}" -o "$#" -lt 1 -o "$#" -gt 2 ]; then
  echo "usage error: usage: ${BN} <target-project-name> [target-project-owner]" >&2
  exit 1
fi

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
