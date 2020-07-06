#!/usr/bin/env bash
set -eo pipefail

: ${PROJECT_DIR?Must set PROJECT_DIR}
BACKEND="${PROJECT_DIR}/backend"
source "${BACKEND}/scripts/lib/init.sh"

usage="generate-all.sh [simple|all]

    all       (default) run all generations
    simple    only run generator (without wire)
"

case "$1" in
""|"simple")
    ;;
*)
    printf "$usage"
    exit 2
esac

generator=$(::get cmd o.o/backend/tools/cmd/generator)
"${generator}" -ignored-plugins=sample o.o/...

if [[ $1 != "simple" ]]; then
    go install github.com/google/wire/cmd/wire
    wire o.o/...
fi
