#!/bin/bash
set -e

: "${COMMIT?Must set COMMIT}"

build() {
    FILE=$1 ; shift
    NAME=$(echo $FILE | rev | cut -f1 -d'/' | rev)
    echo $NAME "$@"

    # enable CGO for building github.com/lfittl/pg_query_go
    CGO_ENABLED=1 go build "$@" \
        -tags release \
        -ldflags "-X \"etop.vn/backend/pkg/common.commit=${COMMIT}\"" \
        -o bin/$NAME $FILE
}

if [[ -n $ENV_FILE ]]; then source "$ENV_FILE" ; fi

# build
go version
time go install ./...

build ./cmd/etop-server           $BUILD_SERVER
build ./cmd/fabo-server           $BUILD_FABO_SERVER
build ./cmd/etop-event-handler    $BUILD_EVENT_HANDLER
build ./cmd/etop-uploader         $BUILD_UPLOADER
build ./cmd/pgevent-forwarder     $BUILD_PGEVENT_FORWARDER
build ./cmd/shipping-sync-service $BUILD_SYNC_SERVICE
build ./cmd/etop-notifier         $BUILD_NOTIFIER
build ./cmd/etop-etl              $BUILD_ETL

mkdir -p bin/com/web/ecom
   cp -R     com/web/ecom/assets    bin/com/web/ecom/
   cp -R     com/web/ecom/templates bin/com/web/ecom/

mkdir -p bin/zexp/etl
   cp -R     zexp/etl/db   bin/zexp/etl/
