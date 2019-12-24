#!/bin/bash
set -e

: ${ETOPDIR?Must set ETOPDIR}
BACKEND="${ETOPDIR}/backend"
wd=$(pwd)

"$BACKEND"/scripts/clean-release.sh

# generate go-bindata
cd "${ETOPDIR}/backend/doc"
go-bindata -pkg doc -o zz_release.bindata.go -ignore '\.(md|go|xlsx)$' ./...
cd $wd
