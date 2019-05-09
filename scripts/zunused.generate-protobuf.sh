#!/usr/bin/env bash
set -eo pipefail

: ${ETOPDIR?Must set ETOPDIR}
BACKEND="${ETOPDIR}/backend"
BASEDIR="$(realpath ${ETOPDIR}/..)"
source "${BACKEND}/scripts/lib/init.sh"

gotoprotobuf=$(::get cmd etop.vn/backend/up/gogen/cmd/go-to-protobuf)

# assemble protobuf imports

PROTOBASE="${BACKEND}/_out/import_proto"
mkdir -p "${PROTOBASE}"

mkdir -p "${PROTOBASE}/github.com/gogo"
rm "${PROTOBASE}/github.com/gogo/protobuf" || true

ln -s \
  "$(::get mod path github.com/gogo/protobuf)" \
  "${PROTOBASE}/github.com/gogo/protobuf"

# generate .proto and .pb.go files

PACKAGES=(
	etop.vn/api/main/credit/v1
	etop.vn/api/meta/v1
)

"${gotoprotobuf}" \
  --output-base="${BASEDIR}" \
  --packages=$(IFS=, ; echo "${PACKAGES[*]}") \
  --proto-import="${PROTOBASE}" \
  --go-header-file "${BACKEND}/scripts/res/boilerplate.generate.txt" \
  "$@"
