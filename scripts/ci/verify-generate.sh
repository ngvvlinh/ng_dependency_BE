#!/bin/bash
set -eo pipefail

scripts/protobuf-gen.sh
scripts/generate-all.sh
CHANGES="$(git status -s)"
if [[ ! -z "$CHANGES" ]]; then
  echo "$CHANGES"
  echo "\nGenerated files are not up to date!"
  exit 1
fi
