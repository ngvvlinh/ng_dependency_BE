#!/bin/bash
set -eo pipefail

fail() { echo "$@" ; exit 1 ; }

if git show HEAD | grep 'github.com/k0kubun/''pp' | grep '+' ; then
  fail Must remove debug instructions
fi
