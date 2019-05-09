#!/bin/bash
set -e
ARG1=$1
if [ -z "$ARG1" ]; then
    ARG1="."
fi
ARG2=$2
if [ -z "$ARG2" ]; then
    ARG2="__graph.png"
fi

go get github.com/kisielk/godepgraph
godepgraph -s -p github.com,go.uber.org,golang.org,google.golang.org,gopkg.in $ARG1 | dot -Tpng -o $ARG2

open $ARG2
