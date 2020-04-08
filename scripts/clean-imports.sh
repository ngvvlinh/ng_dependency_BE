#!/bin/bash
set -eo pipefail

: ${ETOPDIR?Must set ETOPDIR}
BACKEND="${ETOPDIR}/backend"
source "${BACKEND}/scripts/lib/init.sh"

cleanImports=$(::get cmd etop.vn/backend/tools/cmd/clean-imports)
"${cleanImports}" "${ETOPDIR}/backend"
"${cleanImports}" -name '^(zz_generated.+)|(.+\.gen)\.go$' -check-alias "${ETOPDIR}/backend"
