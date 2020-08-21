#!/usr/bin/env bash
set -e

: ${PROJECT_DIR?Must set PROJECT_DIR}
BACKEND="${PROJECT_DIR}/backend"
source "${BACKEND}/scripts/lib/init.sh"

deps=$(::get cmd o.o/backend/tools/cmd/deps)
for path in "$PROJECT_DIR"/backend/cmd/*; do
    files=$(ls $path/*.go 2>/dev/null || true)
    if [[ -n $files ]]; then
        >&2 echo $path
        "$deps" "$path" > "${path}/__deps"
    fi
done

printf "\n✔ done!\n"
