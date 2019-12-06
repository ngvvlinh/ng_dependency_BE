#!/usr/bin/env bash
set -eo pipefail

: ${ETOPDIR?Must set ETOPDIR}
BACKEND="${ETOPDIR}/backend"
source "${BACKEND}/scripts/lib/init.sh"

generator=$(::get cmd etop.vn/backend/tools/cmd/generator)
"${generator}"         -plugins=enum        etop.vn/...
"${generator}" -ignored-plugins=enum,sample etop.vn/...
