#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

vendor/k8s.io/code-generator/generate-groups.sh \
deepcopy \
github.com/pivotal-cf-experimental/mysql-operator/pkg/generated \
github.com/pivotal-cf-experimental/mysql-operator/pkg/apis \
binding:v1alpha1 \
--go-header-file "./tmp/codegen/boilerplate.go.txt"
