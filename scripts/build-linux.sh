#!/bin/bash
set -e

: ${ETOPDIR?Must set ETOPDIR}
BACKEND="${ETOPDIR}/backend"

replace() { echo "$1" | sed "s/$2/$3/g"; }

COMMIT=$(git log -10 --pretty='↵%h <%ae> %B' | grep -E '^(↵[0-9a-f]{6,10} )|(Change-Id:)|(Issue:)')
COMMIT=$(echo -e "${COMMIT}\n\n@thangtran268")
COMMIT=$(echo "${COMMIT}" | tr '\n' '↵' | sed 's/\s/·/g')

COMMIT=$(replace "$COMMIT" "<builamquangngoc91@gmail.com>" "@quangngoc430")
COMMIT=$(replace "$COMMIT" "<congvan2498@gmail.com>" "@congvan2498")
COMMIT=$(replace "$COMMIT" "<huynhhainam96qt@gmail.com>" "@hai_nam_qt")
COMMIT=$(replace "$COMMIT" "<olvrng@gmail.com>" "@vunmq")
COMMIT=$(replace "$COMMIT" "<tuan@eye-solution.vn>" "@tuanpn")

function build() {
    FILE=$1
    NAME=$(echo $FILE | rev | cut -f1 -d'/' | rev)
    echo $NAME
    CGO_ENABLED=0 GOOS=linux go build \
        -tags release \
        -ldflags "-X etop.vn/backend/pkg/common.commit=${COMMIT}" \
        -o bin/$NAME $FILE
}

# generate
"$BACKEND"/scripts/generate-release.sh

# build
build ./cmd/etop-server
build ./cmd/etop-event-handler
build ./cmd/etop-uploader
build ./cmd/pgevent-forwarder
build ./cmd/shipping-sync-service
build ./cmd/etop-notifier
build ./cmd/haravan-gateway
build ./cmd/supporting/crm-sync-service

# clean up
"$BACKEND"/scripts/clean-release.sh
