#!/bin/bash
DN=$(dirname -- "$0")
BN=$(basename -- "$0")
source "${DN}/common"

usage() {
  cat <<EOT
usage: ${BN} OPTIONS <version>

OPTIONS:
  -h|--help|help  ... display this usage information and exit

environment variables:
  FORCE ... if set to non-empty, the tag will be recreated if it exists
EOT
  exit 1
}

[ "$1" != '-h' -a "$1" != '--help' -a "$1" != 'help'  ] || usage

VER="$1"
[ -n "${VER}" -a "$#" -eq 1 ] || usage

echo "release ${VER}..."

git diff --cached --quiet &>/dev/null || die "repo is not clean (changes in staging area)"
git diff --quiet &>/dev/null || die "repo is not clean (changes not staged)"

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

echo "release ${VER}: SUCCEEDED"
