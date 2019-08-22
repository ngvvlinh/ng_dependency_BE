#!/usr/bin/env bash
set -eo pipefail

: ${ETOPDIR?Must set ETOPDIR}
BACKEND="${ETOPDIR}/backend"
source "${BACKEND}/scripts/lib/init.sh"

apiDir="${BACKEND}/up/etop.vn/api"
packages="$(find "${apiDir}" -type d | grep -v __)"

# generate commands & queries
packagesPath=""
for pkg in ${packages}; do
  pkgPath=${pkg#"$BACKEND/up/"}
  if ls "${pkg}"/*.go    1>/dev/null 2>/dev/null && \
   ! ls "${pkg}"/*.proto 1>/dev/null 2>/dev/null ; then
    packagesPath="$packagesPath $pkgPath"
  fi
done

generator=$(::get cmd etop.vn/backend/tools/cmd/generator)
"${generator}" -plugins=cq $packagesPath
