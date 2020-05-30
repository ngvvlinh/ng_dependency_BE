#!/usr/bin/env bash
set -eo pipefail

: ${PROJECT_DIR?Must set PROJECT_DIR}
BACKEND="${PROJECT_DIR}/backend"
source "${BACKEND}/scripts/lib/init.sh"

generator=$(::get cmd o.o/backend/tools/cmd/generator)
"${generator}" -ignored-plugins=sample o.o/...

go install github.com/google/wire/cmd/wire
wire o.o/...
