#!/bin/bash
set -e

if diff --brief --recursive --new-file .mod /go/pkg/mod > /dev/null ; then
  echo "go mod: skip saving"
  exit 0
fi

# because Go cache is immutable, we can safely skip existing files
cp -R /go/pkg/mod .mod
if [[ -d /go/bin ]]; then cp -R /go/bin .mod/ ; fi
