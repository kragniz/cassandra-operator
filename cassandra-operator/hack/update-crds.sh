#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

REPO_ROOT="$(cd "$(dirname "$0")" && pwd -P)"/..
cd "${REPO_ROOT}"

output="$(mktemp -d)"
controller-gen crd:trivialVersions=true paths=./pkg/apis/... output:crd:dir="${output}"
cp "${output}"/core.sky.uk_cassandras.yaml \
    "${REPO_ROOT}/kubernetes-resources/cassandra-operator-crd.yml"
