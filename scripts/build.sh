#!/bin/bash
set -e

: ${ETOPDIR?Must set ETOPDIR}
BACKEND="${ETOPDIR}/backend"
USAGE="Usage: build.sh [docker]"

replace() { echo "$1" | sed "s/$2/$3/g"; }

preprocess() {
    rm "$BACKEND"/bin/* || true
    "$BACKEND"/scripts/generate-release.sh

    COMMIT=$(git log -10 --pretty='¶%h <%ae> %B' | grep -E '^(¶[0-9a-f]{6,10} )|(Change-Id:)|(Issue:)')
    COMMIT=$(echo -e "${COMMIT}\n\n@thangtran268")
    COMMIT=$(echo "${COMMIT}" | tr '\n' '¶' | sed 's/\s/·/g')

    COMMIT=$(replace "$COMMIT" "<builamquangngoc91@gmail.com>" "@quangngoc430")
    COMMIT=$(replace "$COMMIT" "<congvan2498@gmail.com>" "@congvan2498")
    COMMIT=$(replace "$COMMIT" "<huynhhainam96qt@gmail.com>" "@hai_nam_qt")
    COMMIT=$(replace "$COMMIT" "<olvrng@gmail.com>" "@vunmq")
    COMMIT=$(replace "$COMMIT" "<tuan@eye-solution.vn>" "@tuanpn")
}

build_docker() {
    if docker ps -a | grep etop_golang_alpine | grep Exited ; then
        docker start etop_golang_alpine
    fi
    if ! docker ps | grep etop_golang_alpine ; then
        docker run -d --name etop_golang_alpine \
            -e 'ETOPDIR=/etop.vn' \
            -v "$PWD":/etop.vn/backend \
            -w /etop.vn/backend olvrng/golang-alpine-toolbox \
            sleep 3600
    fi

    docker exec -it -e COMMIT="$COMMIT" \
        etop_golang_alpine scripts/build-inner.sh
}

case "$1" in
"")
    preprocess
    COMMIT="$COMMIT" scripts/build-inner.sh
    ;;
docker)
    preprocess
    build_docker
    ;;
*)
    echo "$USAGE"
    exit 2
esac
