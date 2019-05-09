#!/bin/bash
set -e
ARG=$1
if [ -z "$ARG" ]; then
    ARG="."
fi

rm derived.gen.go || true

go install etop.vn/backend/up/gogen/cmd/goderive
goderive $ARG
goimports -w $ARG
