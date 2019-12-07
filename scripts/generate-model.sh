#!/bin/bash
set -ex

list=$(grep --include=\*.go -rnw . -e 'go:generate' | grep derive | cut -d':' -f1)
list=$(echo $list | sed -E 's/\/[a-z]+\.go//g')

go install etop.vn/backend/tools/cmd/goderive
goderive $list

genfiles=$(find . -name 'derived.gen.go' -o -name 'filters.gen.go')
goimports -w $genfiles
