#!/bin/bash
set -eo pipefail
wd=$(pwd)

tidy() {
  cd $1 && go mod tidy && cd $wd
}

scripts/protobuf-gen.sh
scripts/generate-all.sh

go mod tidy
tidy $wd/up/etop.vn/api
tidy $wd/up/etop.vn/capi
tidy $wd/up/etop.vn/common

CHANGES="$(git status -s)"
if [[ ! -z "$CHANGES" ]]; then
  echo "$CHANGES"
  echo "\nGenerated files are not up to date!"
  exit 1
fi
