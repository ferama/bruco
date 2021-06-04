#!/usr/bin/env bash
set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname "${BASH_SOURCE[0]}")/../..
CODEGEN_PKG=${CODEGEN_PKG:-$(cd "${SCRIPT_ROOT}"; ls -d -1 ./vendor/k8s.io/code-generator 2>/dev/null || echo ../code-generator)}

bash "${CODEGEN_PKG}"/generate-groups.sh "deepcopy,client,informer,lister" \
  github.com/ferama/bruco/pkg/kube/generated github.com/ferama/bruco/pkg/kube/apis \
  brucocontroller:v1alpha1 \
  --output-base "$SCRIPT_ROOT" \
  --go-header-file "${SCRIPT_ROOT}"/hack/k8s/boilerplate.go.txt

rsync --recursive \
 --remove-source-files \
 --ignore-times \
 "${SCRIPT_ROOT}/github.com/ferama/bruco/pkg" "${SCRIPT_ROOT}"

rm -rf "${SCRIPT_ROOT}/github.com"
