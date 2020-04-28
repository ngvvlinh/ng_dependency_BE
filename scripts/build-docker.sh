#!/bin/bash

: ${PROJECT_DIR?Must set PROJECT_DIR}
BACKEND="${PROJECT_DIR}/backend"

"$BACKEND"/scripts/build.sh docker
