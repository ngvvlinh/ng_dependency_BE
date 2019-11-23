#!/bin/bash
set -eo pipefail

: "${ETOPDIR?Must set ETOPDIR}"
BACKEND="${ETOPDIR}/backend"
source "${BACKEND}/scripts/lib/init.sh"

"${BACKEND}/scripts/install-tools.sh"

# install tools
go install \
    "$(::get mod path github.com/gogo/protobuf)/protoc-gen-gogo"   \
    "$(::get mod path github.com/golang/protobuf)/protoc-gen-go"

ROOTDIR=$BACKEND/up # the root of import path
PBDIR=$BACKEND/up/etop.vn/api/pb
IMPORT="-I$PBDIR \
    -I$(::get mod path github.com/gogo/protobuf)"

clean() {
    FILES=$1
    if ls $FILES 1>/dev/null 2>/dev/null; then
        rm $FILES
    fi
}

for PKG in $(find "$PBDIR" -type d); do
    clean $PKG/*.pb.go
done

for PKG in $(find "$PBDIR" -type d); do
    PROTO=$PKG/*.proto
    if ls $PROTO 1>/dev/null 2>/dev/null; then
        protoc $IMPORT --gogo_out=$ROOTDIR $PROTO
        echo "Generated from: $PKG"
    fi
done

printf "\n✔ Done\n"
