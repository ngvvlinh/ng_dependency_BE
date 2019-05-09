#!/bin/bash
set -e

go version
packages=( $(go list ./... | grep -v '/tests/') )
mode=atomic

for pkg in "${packages[@]}"; do
  file=${pkg//\//--}.cover
  go test -v -race -covermode="$mode" -coverprofile="$file" "$pkg"
done
