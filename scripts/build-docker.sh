#!/bin/bash

: ${ETOPDIR?Must set ETOPDIR}
BACKEND="${ETOPDIR}/backend"

"$BACKEND"/scripts/build.sh docker
