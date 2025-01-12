#!/bin/bash
DN=$(dirname -- "$0")
BN=$(basename -- "$0")
source "${DN}/common"

usage() {
  cat <<EOT
usage: ${BN} OPTIONS

OPTIONS:
  -h|--help|help  ... display this usage information and exit

environment variables:
  NO_PUSH ...  when set to non-empty, the images will not be pushed
EOT
  exit 1
}

[ "$1" != '-h' -a "$1" != '--help' -a "$1" != 'help'  ] || usage


echo "snapshot..."

DOCKER_LOGGED_IN=$(cat ~/.docker/config.json 2>/dev/null | jq -r '.auths."ghcr.io".auth' 2>/dev/null)

setNextVersion
gitMustBeClean pkg/common/about.go

sed -i "s|\(Version\s*=\s*\"\)[^\"]*\"|\1${NEXT_VERSION}\"|" pkg/common/about.go

# environment variables:
#   NO_PUSH ... if set to non-empty, the images will not be pushed

export GORELEASER_CURRENT_TAG="v${VERSION}"
# NOTE: goreleaser increments the current patch version from the '${GORELEASER_CURRENT_TAG}' to get the new version
goreleaser release --snapshot --clean

git checkout pkg/common/about.go

if [ -n "${NO_PUSH}" -o -z "${DOCKER_LOGGED_IN}" -o "${DOCKER_LOGGED_IN}" = "null" ]; then
  echo "SUCCEEDED (docker images not pushed)"
  docker tag "ghcr.io/gms1/go-project-template:${NEXT_VERSION}-SNAPSHOT-amd64" "ghcr.io/gms1/go-project-template:${NEXT_VERSION}-SNAPSHOT"
  exit 0
fi

docker push "ghcr.io/gms1/go-project-template:${NEXT_VERSION}-SNAPSHOT-arm64"
docker push "ghcr.io/gms1/go-project-template:${NEXT_VERSION}-SNAPSHOT-amd64"

docker manifest rm "ghcr.io/gms1/go-project-template:${NEXT_VERSION}-SNAPSHOT" &>/dev/null || true

docker manifest create "ghcr.io/gms1/go-project-template:${NEXT_VERSION}-SNAPSHOT" \
  "ghcr.io/gms1/go-project-template:${NEXT_VERSION}-SNAPSHOT-arm64" \
  "ghcr.io/gms1/go-project-template:${NEXT_VERSION}-SNAPSHOT-amd64" \

docker manifest push "ghcr.io/gms1/go-project-template:${NEXT_VERSION}-SNAPSHOT"

docker manifest rm "ghcr.io/gms1/go-project-template:${NEXT_VERSION}-SNAPSHOT" &>/dev/null || true

docker pull "ghcr.io/gms1/go-project-template:${NEXT_VERSION}-SNAPSHOT"

docker image prune -f

echo "snapshot: SUCCEEDED"
