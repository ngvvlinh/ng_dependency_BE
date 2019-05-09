#!/bin/bash
set -e

mode=atomic
profile=coverage.out

set --
set -- *.cover
if [[ $# -gt 0 ]]; then
  printf 'mode: %s\n' "$mode" > "$profile"
  grep   -h -v -- "^mode:" *.cover \
  | grep -h -v '\.pb\.go'   \
  | grep -h -v '.twirp\.go' \
  | grep -h -v '.gen\.go'   \
  >> "$profile"
  go tool cover -func="$profile" | grep -E "^total:"
fi

mkdir -p artifacts/coverage
mv *.cover $profile artifacts/coverage
