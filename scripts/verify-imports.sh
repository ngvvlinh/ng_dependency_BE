#!/usr/bin/env bash
set -eox pipefail

BACKEND=$(dirname "${BASH_SOURCE}")/..
source "${BACKEND}/scripts/lib/init.sh"

importverifier=$(::get cmd etop.vn/backend/tools/cmd/importverifier)
verifyimports=$(::get cmd etop.vn/backend/tools/cmd/verify-imports)

"$importverifier" "etop.vn/backend/" "${BACKEND}/scripts/import-restrictions.yaml"

"$verifyimports" -base "etop.vn/backend/" -dir "${BACKEND}" \
	"etop.vn/backend/pkg/..."
