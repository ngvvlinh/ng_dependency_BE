#!/bin/bash
set -e

: ${ETOPDIR?Must set ETOPDIR}
BACKEND="${ETOPDIR}/backend"

# remove generated files
find "$BACKEND" -name 'zz_release.*.go' -exec rm {} +
