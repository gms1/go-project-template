#!/bin/bash
DN=$(dirname -- "$0")
BN=$(basename -- "$0")
source "${DN}/common"

usage() {
  cat <<EOT
usage: ${BN} OPTIONS <version>

creates a release using the specified <version>

OPTIONS:
  -f|--force ... force recreating of the tag if it already exists
  -h|--help  ... display this usage information and exit

EOT
  exit 1
}


OPTS=(f force)
getopt "$@"

[ "${#ARGS[@]}" -eq 1] || usage

VER="${ARGS[0]}"
[ -n "${VER}" ] || usage

[ -z "${f}" ] || force="${f}"

info "release ${VER}..."

gitMustBeClean

sed -i "s|\(Version\s*=\s*\"\)[^\"]*\"|\1${VER}\"|" "${ABOOUT_GO}"

go build -v ./...

git add "${ABOOUT_GO}"

./go-project-template docs ./docs/
git add docs

if ! git diff --cached --quiet; then

  git commit -m "release: ${VER}"
  git push

fi

if [ -n "${force}" ]; then
  git push origin ":refs/tags/v${VER}"
fi
git tag "v${VER}" -f -a -m "release: ${VER}"
git push origin --tags

succeeded "release" "creating version '${VER}'"
