#!/bin/bash
set -e

COMMIT=$(git log -1 --pretty='%h %B')
COMMIT=$(echo "${COMMIT}" | tr '\n' ' ' | sed 's/\s/·/g')
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
