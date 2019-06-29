#!/usr/bin/env bash
set -eo pipefail

: ${ETOPDIR?Must set ETOPDIR}
BACKEND="$ETOPDIR/backend"

::get() { ${BACKEND}/scripts/lib/get.sh $@ ; }

::clean() {
    for FILES in $@; do
        if ls "${FILES}" 1>/dev/null 2>/dev/null; then
            rm "${FILES}"
        fi
    done
}
