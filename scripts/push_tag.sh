#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

VERSION="${1}"

if [[ -z "${VERSION}" ]]; then
  echo "Usage: push_tag.sh <version>"
  exit 1
fi

git tag -a "v${VERSION}" -m "version ${VERSION}"
git push origin "v${VERSION}"
