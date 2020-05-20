#!/bin/bash
set -eo pipefail

fail() { echo "$@" ; exit 1 ; }

blacklist='("github.com/k0kubun/pp)\s*$'

if git show HEAD | grep -E "$blacklist" | grep '^+' ; then
  fail Must remove debug instructions
fi
