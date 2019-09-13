#!/usr/bin/env bash
set -e

check_flag() {
  : {"$1"?Missing required param}
  git log -1 | grep "\bFlags:" | grep "\b$1\b"
}
check_branch() {
  : {"$1"?Missing required param}
  echo "$CI_COMMIT_REF_NAME" | grep "\b$1\b"
}

repos="
etop.vn/api/...
etop.vn/backend/...
etop.vn/capi/...
etop.vn/common/...
"

go install ${repos}
packages="$(go list ${repos} | grep -v '/tests/')"

# execute tests in coverage mode in one of those conditions
#   1. environment variable COVER is set
#   2. branch name contains "master"
#   3. commit message contains: "Flags: coverage"
set +e
cover="${COVER}$(check_branch master)$(check_flag coverage)"
set -e

if [[ -z "$cover" ]]; then
  go test -v -race ${packages}
else
  echo "Executing tests in coverage mode"
  mode=atomic
  file=coverage.cover
  go test -v -race -covermode="$mode" -coverprofile="$file" ${packages}
  scripts/ci/show-coverage.sh
fi
