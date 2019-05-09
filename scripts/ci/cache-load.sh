#!/bin/bash
set -e
mkdir -p /go/pkg
if [[ -d /go/pkg/mod ]]; then rm -r /go/pkg/mod ; fi
if [[ -d .mod ]]; then
  mv .mod /go/pkg/mod
  echo "go mod: cache loaded"
else
  echo "go mod: no cache"
fi
