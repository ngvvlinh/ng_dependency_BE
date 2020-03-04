#!/bin/bash
set -e

: "${COMMIT?Must set COMMIT}"

build() {
    FILE=$1
    NAME=$(echo $FILE | rev | cut -f1 -d'/' | rev)
    echo $NAME
    CGO_ENABLED=1 go build \
        -tags release \
        -ldflags "-X \"etop.vn/backend/pkg/common.commit=${COMMIT}\"" \
        -o bin/$NAME $FILE
}

# build
go version

build ./cmd/etop-server
build ./cmd/etop-event-handler
build ./cmd/etop-uploader
build ./cmd/pgevent-forwarder
build ./cmd/shipping-sync-service
build ./cmd/etop-notifier
build ./cmd/etop-etl

mkdir -p bin/com/web/ecom
   cp -R     com/web/ecom/assets    bin/com/web/ecom/
   cp -R     com/web/ecom/templates bin/com/web/ecom/

mkdir -p bin/zexp/etl
   cp -R     zexp/etl/db   bin/zexp/etl/
