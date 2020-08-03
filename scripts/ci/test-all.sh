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

go install o.o/...

# # execute tests in coverage mode in one of those conditions
# #   arguments run test-all.sh = simple
# #   1. environment variable COVER is set
# #   2. branch name contains "master"
# #   3. commit message contains: "Flags: coverage"
set +e
cover="${COVER}$(check_branch master)$(check_flag coverage)"
set -e

exec_coverage() {
  echo "Executing tests in coverage mode"
  mode=atomic
  file=coverage.cover
  go test -v -race -covermode="$mode" -coverprofile="$file" ${packages}
  scripts/ci/show-coverage.sh
}

if [[ -z "$1" ]]; then
  # run test simple
  packages="$(go list o.o/... | grep -v '/tests/')"
  if [[ -z "$cover" ]]; then
    go test -v -race ${packages}
  else
    exec_coverage
  fi
  # run test e2e
  packages="$(go list o.o/... | grep '/tests/')"

  go test -v -race ${packages}
else
  # only run simple
  if [[ -n "$1" && $1 = "simple" ]]; then
    packages="$(go list o.o/... | grep -v '/tests/')"
    if [[ -z "$cover" ]]; then
      go test -v -race ${packages}
    else
      exec_coverage
    fi
  fi

  # execute tests in coverage mode in one of those conditions
  # #   . arguments run test-all.sh = e2e

  if [[ -n "$1" && $1 = "e2e" ]]; then
    packages="$(go list o.o/... | grep '/tests/')"

    go test -v -race ${packages}
  fi
fi
