#!/bin/bash
DN=$(dirname -- "$0")
BN=$(basename -- "$0")

echo "release..."

# environment variables:
#   FORCE ... if set to non-empty, the tag will be recreated if it exists

errorHandler() {
  local script="$0"
  local line="$1"
  local command="$2"
  local code="$3"
  echo "${script}: line: ${line}: command \"${command}\" exited with ${code}" >&2
}

trap 'errorHandler "$LINENO" "$BASH_COMMAND" "$?"' "ERR"

set -e
set -o pipefail
cd "${DN}/.."

VER="$1"
if [ -z "${VER}" -o "$#" -ne 1 ]; then
  echo "usage error: usage: ${BN} <version>" >&2
  exit 1
fi

if ! git diff --cached --quiet &>/dev/null || ! git diff --quiet &>/dev/null; then
  echo "repo is not clean" >&2
  exit 1
fi

sed -i "s|\(Version\s*=\s*\"\)[^\"]*\"|\1${VER}\"|" pkg/common/about.go

go build -v ./...

git add pkg/common/about.go

./go-project-template docs ./docs/
git add docs

if ! git diff --cached --quiet; then

  git commit -m "release: ${VER}"
  git push

fi

if [ -n "${FORCE}" ]; then
  git push origin ":refs/tags/v${VER}"
fi
git tag "v${VER}" -f -a -m "release: ${VER}"
git push origin --tags

echo "release: SUCCEEDED"
