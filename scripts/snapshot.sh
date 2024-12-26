#!/bin/bash
DN=$(dirname -- "$0")
BN=$(basename -- "$0")

echo "snapshot..."

# environment variables:
#   NO_PUSH ... if set to non-empty, the images will not be pushed

DOCKER_LOGGED_IN=$(cat ~/.docker/config.json 2>/dev/null | jq -r '.auths."ghcr.io".auth' 2>/dev/null)

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

VERSION=$(awk -F\" '/Version/{print $2}' pkg/common/about.go)

export GORELEASER_CURRENT_TAG="v${VERSION}"

goreleaser release --snapshot --clean

# NOTE: goreleaser increments the current patch version from the '${GORELEASER_CURRENT_TAG}' to get the new version

VERSION_PARTS=( $(echo ${VERSION//./ } ))

NEW_PATCH_VERSION="${VERSION_PARTS[2]}"
NEW_PATCH_VERSION="$((NEW_PATCH_VERSION + 1))"
NEW_VERSION="${VERSION_PARTS[0]}.${VERSION_PARTS[1]}.${NEW_PATCH_VERSION}"

if [ -n "${NO_PUSH}" -o -z "${DOCKER_LOGGED_IN}" -o "${DOCKER_LOGGED_IN}" = "null" ]; then
  echo "SUCCEEDED (docker images not pushed)"
  docker tag "ghcr.io/gms1/go-project-template:${NEW_VERSION}-SNAPSHOT-amd64" "ghcr.io/gms1/go-project-template:${NEW_VERSION}-SNAPSHOT"
  exit 0
fi

docker push "ghcr.io/gms1/go-project-template:${NEW_VERSION}-SNAPSHOT-arm64"
docker push "ghcr.io/gms1/go-project-template:${NEW_VERSION}-SNAPSHOT-amd64"

docker manifest rm "ghcr.io/gms1/go-project-template:${NEW_VERSION}-SNAPSHOT" &>/dev/null || true

docker manifest create "ghcr.io/gms1/go-project-template:${NEW_VERSION}-SNAPSHOT" \
  "ghcr.io/gms1/go-project-template:${NEW_VERSION}-SNAPSHOT-arm64" \
  "ghcr.io/gms1/go-project-template:${NEW_VERSION}-SNAPSHOT-amd64" \

docker manifest push "ghcr.io/gms1/go-project-template:${NEW_VERSION}-SNAPSHOT"

docker manifest rm "ghcr.io/gms1/go-project-template:${NEW_VERSION}-SNAPSHOT" &>/dev/null || true

docker pull "ghcr.io/gms1/go-project-template:${NEW_VERSION}-SNAPSHOT"

docker image prune -f

echo "snapshot: SUCCEEDED"
