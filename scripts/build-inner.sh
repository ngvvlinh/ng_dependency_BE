#!/bin/bash
set -e

: "${COMMIT?Must set COMMIT}"

# this will copy the source to $2 and start building in the new location
bindir=bin
if [[ -n "$1$2" ]]; then
    bindir="$1/bin"
    cd "$1"
    mkdir -p "$1/bin"
    mkdir -p "$2" && rm -rf "$2" && cp -r "$1" "$2"
    cd "$2"
    echo "source copied to $2"
fi

build() {
    FILE=$1 ; shift
    NAME=$(echo $FILE | rev | cut -f1 -d'/' | rev)
    echo $NAME "$@"

    # enable CGO for building github.com/lfittl/pg_query_go
    CGO_ENABLED=1 go build "$@" \
        -tags release \
        -ldflags "-X \"o.o/backend/pkg/common.commit=${COMMIT}\"" \
        -o "$bindir/$NAME" $FILE
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
build ./cmd/fabo-sync-service     $BUILD_FABO_SYNC_SERVICE

mkdir -p "$bindir"/com/web/ecom
   cp -R     com/web/ecom/assets    "$bindir"/com/web/ecom/
   cp -R     com/web/ecom/templates "$bindir"/com/web/ecom/

mkdir -p "$bindir"/zexp/etl
   cp -R     zexp/etl/db   "$bindir"/zexp/etl/
