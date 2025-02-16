#!/bin/bash
DN=$(dirname -- "$0")
BN=$(basename -- "$0")
source "${DN}/common"

usage() {
  cat <<EOT
usage: ${BN} OPTIONS

creates a snapshot locally

OPTIONS:
  -p|--push ... push images
  -h|--help ... display this usage information and exit

EOT
  exit 1
}

OPTS=(p push)
getopt "$@"

[ "${#ARGS[@]}" -eq 0 ] || usage

[ -z "${p}" ] || push="${p}"

info "snapshot..."

DOCKER_LOGGED_IN=$(cat ~/.docker/config.json 2>/dev/null | jq -r '.auths."ghcr.io".auth' 2>/dev/null)

gitMustBeClean "${ABOOUT_GO}"

NEXT_VERSION="$(getNextVersion)"
sed -i "s|\(Version\s*=\s*\"\)[^\"]*\"|\1${NEXT_VERSION}\"|" "${ABOOUT_GO}"

export GORELEASER_CURRENT_TAG="v$(getVersion)"
# NOTE: goreleaser increments the current patch version from the '${GORELEASER_CURRENT_TAG}' to get the new version
goreleaser release --snapshot --clean

git checkout "${ABOOUT_GO}"

if [ -z "${push}" -o -z "${DOCKER_LOGGED_IN}" -o "${DOCKER_LOGGED_IN}" = "null" ]; then
  docker tag "ghcr.io/gms1/go-project-template:${NEXT_VERSION}-SNAPSHOT-amd64" "ghcr.io/gms1/go-project-template:${NEXT_VERSION}-SNAPSHOT"
  succeeded "docker images are not pushed"
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

succeeded "snapshot"
