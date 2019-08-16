#!/bin/bash
set -eo pipefail

: "${ETOPDIR?Must set ETOPDIR}"
BACKEND="${ETOPDIR}/backend"
source "${BACKEND}/scripts/lib/init.sh"

# install tools
go install \
    $(::get mod path github.com/gogo/protobuf)/protoc-gen-gogo   \
    $(::get mod path github.com/twitchtv/twirp)/protoc-gen-twirp \
    $(::get mod path github.com/golang/protobuf)/protoc-gen-go

ROOTDIR=$ETOPDIR/.. # the root of import path
IMPORT="-I${BACKEND}/pb \
    -I$(::get mod path github.com/gogo/protobuf) \
    -I$(::get mod path github.com/grpc-ecosystem/grpc-gateway)/third_party/googleapis \
    -I$(::get mod path github.com/grpc-ecosystem/grpc-gateway)"

if [ ! -z "$1" ] ; then
    filter="/$1"
    echo "Reading from directory: $1"
    if [ ! -d ${BACKEND}/pb/$1 ]; then
        echo "Invalid directory: ${BACKEND}/pb/$1"
        exit 1
    fi
fi

clean() {
    FILES=$1
    if ls $FILES 1>/dev/null 2>/dev/null; then
        rm $FILES
    fi
}

for PKG in $(find "${BACKEND}/pb${filter}" -type d); do
    clean $PKG/*.pb.go
    clean $PKG/*.pb.gw.go
    clean $PKG/*.pb.gen.go
    clean $PKG/*.pb.twirp.go
    clean $PKG/*.swagger.json
done
for PKG in $(find "${BACKEND}/doc${filter}" -type d); do
    clean $PKG/*.swagger.json
done

# Replace /twirp/ with /api/ (or /v1/)
sedtwirp() {
    FILES=$1
    NEWPATH=$2
    if ls $FILES 1>/dev/null 2>/dev/null; then
        sed -i -e 's/"\/twirp\//"\/'$NEWPATH'\//g' \
               -e 's/Service\//\//g' \
               -e 's/Service"/"/g' \
               -e 's/jsonpb\.Marshaler{OrigName: true}/jsonpb.Marshaler{OrigName: true, EmitDefaults: true}/g' \
               $FILES
    fi
}

prefixpath() {
  if [[ $1 == *"external/"* ]]; then
    echo v1
  else
    echo api
  fi
}

prefixext() {
  if [[ $1 == *"external/"* ]]; then echo "-p ext" ; fi
}

GENERATED_FILES=
for PKG in $(find "${BACKEND}/pb${filter}" -type d); do
    PKGNAME=$(basename $PKG)
    PROTO=$PKG/*.proto
    if ls $PROTO 1>/dev/null 2>/dev/null; then
        protoc $IMPORT --twirp_out=$ROOTDIR --gogo_out=$ROOTDIR --twirp_swagger_out=${BACKEND}/doc $PROTO

        sedtwirp $PKG/*.twirp.go $(prefixpath $PKG)
        echo "Generated from: $PKG"
    fi
    if ls $PKG/*.twirp.go 1>/dev/null 2>/dev/null; then
        GENERATED_FILES="$GENERATED_FILES $PKG/*.twirp.go"
    fi
done

if [[ -n "$GENERATED_FILES" ]]; then goimports -local etop.vn -w $GENERATED_FILES ; fi

# Sort swagger tags and parse @required fields
for FILE in $(find "${BACKEND}/doc" -name *.swagger.json); do
    if [ $(cat $FILE | jq '.paths | length') -eq 0 ]; then
        rm $FILE
        continue
    fi

    cat $FILE \
        | jq '.tags = ([.paths[][].tags[0]] | unique | [.[]|{name:.}])' \
        | jq '
        def reqs:
            .properties
            | to_entries
            | map(select(.value.title // "" | startswith("@required")))
            | map(.key);
        def req:
            if (. | reqs | length) > 0 then
               .required = (. | reqs)
            else . end;
        (.definitions[] | select(.properties != null)) |= req' \
        | sed 's/@required\s*//g' \
        | tr '\n' '~' \
        | sed 's/,~\s*"title":\s*""//g' \
        | tr '~' '\n' \
        > $FILE.jq

    sedtwirp $FILE.jq $(prefixpath $FILE)
    mv $FILE.jq $FILE
    echo "Updated doc:    $FILE"
done

# Generate go-bindata
cd ${BACKEND}/doc
go-bindata -pkg doc -o bindata.gen.go ./...

# Generate wrapper files
WRAPPER_ARGS=
for PKG in $(find "${BACKEND}/pb${filter}" -type d | grep -v common); do
    PKGNAME=$(basename $PKG)
    FILES=$PKG/*.twirp.go
    if ls $FILES 1>/dev/null 2>/dev/null; then
        WRAPPER_ARGS="$WRAPPER_ARGS $(prefixext $PKG) -s pb -o wrapper $FILES"
    fi
done

wrapper_gen=$(::get cmd etop.vn/backend/tools/cmd/wrapper_gen)
wrapper_gen $WRAPPER_ARGS

printf "\n✔ Done\n"
