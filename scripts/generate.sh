#!/bin/bash
set -e
go install ./scripts/cmd/...
if [ -z "$1" ]; then
    go generate ./...
else
    go generate $1
fi
