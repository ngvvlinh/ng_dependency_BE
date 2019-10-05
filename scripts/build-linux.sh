#!/bin/bash
set -e

COMMIT=$(git log -5 --pretty='%h %B' | grep ':' | grep -E '^([0-9a-f]{6,10} )|(Change-Id:)')
COMMIT=$(echo "${COMMIT}" | tr '\n' '⮐' | sed 's/\s/·/g')
function build() {
    FILE=$1
    NAME=$(echo $FILE | rev | cut -f1 -d'/' | rev)
    echo $NAME
    CGO_ENABLED=0 GOOS=linux go build \
        -ldflags "-X etop.vn/backend/pkg/common.commit='${COMMIT}'" \
        -o bin/$NAME $FILE
}

build ./cmd/etop-server
build ./cmd/etop-event-handler
build ./cmd/etop-uploader
build ./cmd/pgevent-forwarder
build ./cmd/shipping-sync-service
build ./cmd/etop-notifier
build ./cmd/haravan-gateway
build ./cmd/supporting/crm-sync-service
