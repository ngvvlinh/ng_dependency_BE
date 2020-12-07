#!/bin/bash
set -e

: "${COMMIT?Must set COMMIT}"
: "${1?Usage: build-inner.sh TARGET [MOUNT_DIR]>}"

binDir=bin
target="$1"
mountDir="$2"
if [[ -n $mountDir ]]; then
    buildDir=/project
    if [[ "$mountDir" == "$buildDir" ]]; then
        echo "invalid MOUNT_DIR"
        exit 1
    fi

    # this will copy the source to buildDir and start building in the new location
    binDir="$mountDir/bin"
    mkdir -p "$binDir"
    cd "$mountDir"
    rm -rf "$buildDir" || true
    cp -r "$mountDir" "$buildDir"
    cd "$buildDir"
    echo "source copied, start building..."
fi

build() {
    FILE=$1 ; shift
    NAME=$(echo $FILE | rev | cut -f1 -d'/' | rev)
    echo $NAME "$@"

    # enable CGO for building github.com/lfittl/pg_query_go
    CGO_ENABLED=1 go build "$@" \
        -tags release \
        -ldflags "-X \"o.o/backend/pkg/common.commit=${COMMIT}\"" \
        -o "$binDir/$NAME" $FILE
}

if [[ -n $ENV_FILE ]]; then source "$ENV_FILE" ; fi

# build
go version
time go install ./...

case "$target" in
etop)
    build ./cmd/etop-server           $BUILD_SERVER
    build ./cmd/fabo-server           $BUILD_FABO_SERVER
    build ./cmd/etop-event-handler    $BUILD_EVENT_HANDLER
    build ./cmd/etop-uploader         $BUILD_UPLOADER
    build ./cmd/pgevent-forwarder     $BUILD_PGEVENT_FORWARDER
    build ./cmd/etop-etl              $BUILD_ETL
    build ./cmd/fabo-sync-service     $BUILD_FABO_SYNC_SERVICE
    build ./cmd/shipment-sync-service $BUILD_SHIPMENT_SYNCE_SERVICE
    build ./cmd/fabo-event-handler    $BUILD_FABO_EVENT_HANDLER
    build ./cmd/telecom-sync-service  $BUILD_TELECOM_SYNC_SERVICE

    mkdir -p "$binDir"/com/web/ecom
       cp -R     com/web/ecom/assets    "$binDir"/com/web/ecom/
       cp -R     com/web/ecom/templates "$binDir"/com/web/ecom/

    mkdir -p "$binDir"/com/report
       cp -R     com/report/templates   "$binDir"/com/report/

    mkdir -p "$binDir"/zexp/etl
       cp -R     zexp/etl/db   "$binDir"/zexp/etl/

    ;;
fabo)
    build ./cmd/fabo-server             $BUILD_FABO_SERVER
    build ./cmd/fabo-event-handler      $BUILD_EVENT_HANDLER
    build ./cmd/etop-uploader           $BUILD_UPLOADER
    build ./cmd/fabo-pgevent-forwarder  $BUILD_PGEVENT_FORWARDER
    build ./cmd/fabo-sync-service       $BUILD_FABO_SYNC_SERVICE
    ;;
*)
    echo unexpected
    exit 1
esac
