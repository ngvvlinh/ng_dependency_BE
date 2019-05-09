#!/usr/bin/env bash
set -eo pipefail

: ${ETOPDIR?Must set ETOPDIR}
BACKEND="${ETOPDIR}/backend"
BASEDIR="$(realpath ${ETOPDIR}/..)"
source "${BACKEND}/scripts/lib/init.sh"

genswaggergo=$(::get cmd etop.vn/backend/up/gogen/cmd/gen-swagger-go)

"${genswaggergo}" \
  "$@"
