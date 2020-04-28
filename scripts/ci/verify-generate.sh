#!/bin/bash
set -eo pipefail
wd=$(pwd)

tidy() {
  cd $1 && go mod tidy && cd $wd
}

scripts/install-tools.sh
scripts/generate-all.sh
scripts/clean-imports.sh

go mod tidy
tidy $wd/up/o.o/api
tidy $wd/up/o.o/capi
tidy $wd/up/o.o/common

CHANGES="$(git status -s)"
if [[ ! -z "$CHANGES" ]]; then
  printf "$CHANGES\n\n"
  printf "Generated files are not up to date!\n\n"
  git diff
  exit 1
fi
