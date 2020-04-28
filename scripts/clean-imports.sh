#!/bin/bash
set -eo pipefail

: ${PROJECT_DIR?Must set PROJECT_DIR}
BACKEND="${PROJECT_DIR}/backend"
source "${BACKEND}/scripts/lib/init.sh"

cleanImports=$(::get cmd o.o/backend/tools/cmd/clean-imports)
"${cleanImports}" "${PROJECT_DIR}/backend"
"${cleanImports}" -name '^(zz_generated.+)|(.+\.gen)\.go$' -check-alias "${PROJECT_DIR}/backend"
