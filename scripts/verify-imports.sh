#!/usr/bin/env bash
set -eox pipefail

BACKEND=$(dirname "${BASH_SOURCE}")/..
source "${BACKEND}/scripts/lib/init.sh"

importverifier=$(::get cmd o.o/backend/tools/cmd/importverifier)
verifyimports=$(::get cmd o.o/backend/tools/cmd/verify-imports)

"$importverifier" "o.o/backend/" "${BACKEND}/scripts/import-restrictions.yaml"

"$verifyimports" -base "o.o/backend/" -dir "${BACKEND}" \
	"o.o/backend/pkg/..."
