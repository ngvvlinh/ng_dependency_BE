#!/usr/bin/env bash

packages="$(go list ./... | grep -v '/tests/')"

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
