#!/bin/bash
set -eo pipefail
wd=$(pwd)

tidy() {
  cd $1 && go mod tidy && cd $wd
}

scripts/protobuf-gen.sh
scripts/generate-all.sh

# Generate go-bindata
cd "${ETOPDIR}/backend/doc"
go-bindata -pkg doc -o bindata.gen.go -ignore '\.(md|go|xlsx)$' ./...
cd $wd

go mod tidy
tidy $wd/up/etop.vn/api
tidy $wd/up/etop.vn/capi
tidy $wd/up/etop.vn/common

CHANGES="$(git status -s)"
if [[ ! -z "$CHANGES" ]]; then
  printf "$CHANGES\n\n"
  printf "Generated files are not up to date!\n\n"
  git diff
  exit 1
fi
