#!/bin/bash
set -e

mode=atomic
profile=artifacts/coverage/coverage.out
output=artifacts/COVERAGE
artifacts_url=artifacts/ARTIFACTS_URL
output_html=artifacts/coverage/coverage.html
output_html_gocov=artifacts/coverage/gocov.html

mkdir -p artifacts/coverage

set --
set -- *.cover
if [[ $# -gt 0 ]]; then
  printf 'mode: %s\n' "$mode" > "$profile"

  grep   -h -v -- "^mode:" *.cover \
  | grep -h -v '\.pb\.go'   \
  | grep -h -v '.twirp\.go' \
  | grep -h -v '.gen\.go'   \
  >> "$profile"

  go tool cover -func="$profile" \
  | grep "^total:" \
  | grep -o "[0-9.]\+%" \
  > "$output"

  go tool cover -html="$profile" -o "$output_html"

  if [[ -n "$CI_JOB_URL" ]]; then
    echo "${CI_JOB_URL}/artifacts" > "$artifacts_url"
  fi
  echo COVERAGE: $(cat "$output")

## gocov
#  : ${PROJECT_DIR?Must set PROJECT_DIR}
#  BACKEND="${PROJECT_DIR}/backend"
#  source "${BACKEND}/scripts/lib/init.sh"
#
#  gocov=$(::get cmd github.com/axw/gocov/gocov)
#  gocov_html=$(::get cmd github.com/matm/gocov-html)
#  "$gocov" convert "$profile" | "$gocov_html" > "$output_html_gocov"
fi

mv *.cover artifacts/coverage
