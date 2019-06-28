#!/usr/bin/env bash
set -e

repos="
etop.vn/api/...
etop.vn/apix/...
etop.vn/apis/...
etop.vn/backend/...
etop.vn/common/...
"

go install ${repos}

packages="$(go list ${repos} | grep -v '/tests/')"

set +e
master="$(echo $CI_COMMIT_REF_NAME | grep master)"
cover="${COVER}${master}"

if [[ -z "$cover" ]]; then
  set -e
  go test -v -race ${packages}
else
  set -e
  mode=atomic
  file=coverage.cover
  go test -v -race -covermode="$mode" -coverprofile="$file" ${packages}
  scripts/ci/show-coverage.sh
fi
