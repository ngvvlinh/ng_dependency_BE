#!/bin/bash
set -e

mode=atomic
profile=artifacts/coverage/coverage.out
output=artifacts/COVERAGE
output_url=artifacts/COVERAGE_URL
output_html=artifacts/coverage.html

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
    echo "${CI_JOB_URL}/artifacts/browse/${output_html}" > "$output_url"
  fi
  echo COVERAGE: $(cat "$output")
fi

mv *.cover artifacts/coverage
