#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

export GO111MODULE="on"

typeset -A licences
licences=(
    [github.com/sky-uk/cassandra-operator/cassandra-operator]=BSD-3-Clause
    [sigs.k8s.io/yaml]=BSD-3-Clause
    [github.com/hashicorp/golang-lru]=MPL-2.0
)

function pkg_dir() {
    echo $(go list -m -f "{{.Dir}}" $1)
}

overrides=""
for pkg in "${!licences[@]}"; do
    overrides=$(printf '%s -o %s=%s ' "$overrides" "$(pkg_dir $pkg)" "${licences[$pkg]}")
done

REPO_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

restricted=$(paste -s -d ',' $REPO_ROOT/restricted-licences.txt)
projects=$(go list -m -f "{{.Dir}}" all)

licence-compliance-checker -L error -E -r $restricted $overrides $projects
