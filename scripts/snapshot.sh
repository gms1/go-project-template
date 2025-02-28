#!/usr/bin/bash
DN=$(dirname -- "$0")
BN=$(basename -- "$0")
source "${DN}/common"

usage() {
  cat <<EOT
usage: ${BN} OPTIONS

creates a snapshot locally without creating and pushing a release tag

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

[ -n "${DOCKER_CLI}" ] || DOCKER_CLI="docker"

export GORELEASER_CURRENT_TAG="v$(getVersion)"
NEXT_VERSION="$(getNextVersion)"
DOCKER_LOGGED_IN=$(cat ~/."${DOCKER_CLI}"/config.json 2>/dev/null | jq -r '.auths."ghcr.io".auth' 2>/dev/null)

info "snapshot..."

# NOTE: goreleaser increments the current patch version from the '${GORELEASER_CURRENT_TAG}' to get the new version
# temporary chaning the version in code, so that the packaged version is in sync
gitMustBeClean "${ABOOUT_GO}"
sed -i "s|\(Version\s*=\s*\"\)[^\"]*\"|\1${NEXT_VERSION}-SNAPSHOT\"|" "${ABOOUT_GO}"

goreleaser release --snapshot --clean

git checkout "${ABOOUT_GO}"

if [ -z "${push}" -o -z "${DOCKER_LOGGED_IN}" -o "${DOCKER_LOGGED_IN}" = "null" ]; then
  "${DOCKER_CLI}" tag "ghcr.io/gms1/go-project-template:${NEXT_VERSION}-SNAPSHOT-amd64" "ghcr.io/gms1/go-project-template:${NEXT_VERSION}-SNAPSHOT"
  succeeded ""${DOCKER_CLI}" images are not pushed"
fi

"${DOCKER_CLI}" push "ghcr.io/gms1/go-project-template:${NEXT_VERSION}-SNAPSHOT-arm64"
"${DOCKER_CLI}" push "ghcr.io/gms1/go-project-template:${NEXT_VERSION}-SNAPSHOT-amd64"

"${DOCKER_CLI}" manifest rm "ghcr.io/gms1/go-project-template:${NEXT_VERSION}-SNAPSHOT" &>/dev/null || true

"${DOCKER_CLI}" manifest create "ghcr.io/gms1/go-project-template:${NEXT_VERSION}-SNAPSHOT" \
  "ghcr.io/gms1/go-project-template:${NEXT_VERSION}-SNAPSHOT-arm64" \
  "ghcr.io/gms1/go-project-template:${NEXT_VERSION}-SNAPSHOT-amd64" \

"${DOCKER_CLI}" manifest push "ghcr.io/gms1/go-project-template:${NEXT_VERSION}-SNAPSHOT"

"${DOCKER_CLI}" manifest rm "ghcr.io/gms1/go-project-template:${NEXT_VERSION}-SNAPSHOT" &>/dev/null || true

"${DOCKER_CLI}" pull "ghcr.io/gms1/go-project-template:${NEXT_VERSION}-SNAPSHOT"

"${DOCKER_CLI}" image prune -f

succeeded "snapshot"
