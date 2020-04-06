#!/bin/bash
set -e

: ${ETOPDIR?Must set ETOPDIR}
BACKEND="${ETOPDIR}/backend"
USAGE="Usage: build.sh [docker]"

replace() { echo "$1" | sed "s/$2/$3/g"; }

preprocess() {
    rm -r "${BACKEND:?}/bin" || true
    mkdir -p "$BACKEND/bin"
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
    if docker ps -a | grep 'etop_golang$' | grep Exited ; then
        docker start etop_golang
    fi
    if ! docker ps | grep 'etop_golang$' ; then
        docker run -d --name etop_golang \
            -e 'ETOPDIR=/etop.vn' \
            -v "$PWD":/etop.vn/backend \
            -w /etop.vn/backend olvrng/golang-toolbox \
            sleep 3600
    fi

    if [[ -n $ENV_FILE ]]; then _env_file="-e=ENV_FILE=$ENV_FILE" ; fi
    docker exec -it -e COMMIT="$COMMIT" $_env_file \
        etop_golang scripts/build-inner.sh
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
