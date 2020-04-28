#!/bin/bash
set -e

: ${PROJECT_DIR?Must set PROJECT_DIR}
BACKEND="${PROJECT_DIR}/backend"

# remove generated files
find "$BACKEND" -name 'zz_release.*.go' -exec rm {} +
