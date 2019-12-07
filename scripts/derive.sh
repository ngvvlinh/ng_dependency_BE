#!/bin/bash
set -ex
ARG=$1
if [ -z "$ARG" ]; then
    ARG="."
fi

rm derived.gen.go || true

go install etop.vn/backend/tools/cmd/goderive
goderive $ARG

genfiles=$(find .. -name 'derived.gen.go' -o -name 'filters.gen.go')
goimports -w $genfiles
