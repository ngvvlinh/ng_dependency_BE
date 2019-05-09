#!/usr/bin/env bash
set -eo pipefail

BACKEND=$(dirname "${BASH_SOURCE}")/..
source "${BACKEND}/scripts/lib/init.sh"

importverifier=$(::get cmd k8s.io/kubernetes/cmd/importverifier)
verifyimports=$(::get cmd etop.vn/backend/scripts/cmd/verify-imports)

"$importverifier" "etop.vn/backend/" "${BACKEND}/scripts/import-restrictions.yaml"

"$verifyimports" -base "etop.vn/backend/" -dir "${BACKEND}" \
	"etop.vn/backend/pkg/..."
