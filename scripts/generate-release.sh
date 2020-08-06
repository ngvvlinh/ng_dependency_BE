#!/bin/bash
set -e

: ${PROJECT_DIR?Must set PROJECT_DIR}
BACKEND="${PROJECT_DIR}/backend"
wd=$(pwd)

"$BACKEND"/scripts/clean-release.sh

# generate go-bindata
cd "${PROJECT_DIR}/backend/doc"
go-bindata -pkg doc -o zz_release.bindata.go -tags release -ignore '\.(md|go|xlsx)$' ./...
cd $wd

cd "${PROJECT_DIR}/backend/res/dl/fabo"
go-bindata -pkg fabo -o zz_release.bindata.go -tags release -ignore '\.(md|go|xlsx)$' ./...
