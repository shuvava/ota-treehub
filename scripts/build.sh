#!/bin/sh

# Copyright 2016 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o errexit
set -o nounset
set -o pipefail

if [ -z "${OS:-}" ]; then
    echo "OS must be set"
    exit 1
fi
if [ -z "${ARCH:-}" ]; then
    echo "ARCH must be set"
    exit 1
fi
if [ -z "${VERSION:-}" ]; then
    echo "VERSION must be set"
    exit 1
fi
if [ -z "${COMMIT_HASH:-}" ]; then
    echo "COMMIT_HASH must be set"
    exit 1
fi

export CGO_ENABLED=0
export GOARCH="${ARCH}"
export GOOS="${OS}"
export GO111MODULE=on
export GOFLAGS="${GOFLAGS:-} -mod=${MOD}"

PACKAGE="$(go list -m)/pkg/version"
BUILD_TIMESTAMP=$(date '+%Y-%m-%dT%H:%M:%S')


go install                                                      \
    -installsuffix "static"                                     \
    -ldflags "-X '${PACKAGE}.Version=${VERSION}' -X '${PACKAGE}.CommitHash=${COMMIT_HASH}' -X '${PACKAGE}.BuildDate=${BUILD_TIMESTAMP}'"  \
    "$@"
