#!/bin/bash
set -e
if [[ -d .mod ]]; then rm -r .mod ; fi
mv /go/pkg/mod .mod
