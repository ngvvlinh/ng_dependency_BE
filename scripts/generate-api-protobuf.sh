#!/usr/bin/env bash
set -eo pipefail

: ${ETOPDIR?Must set ETOPDIR}
BACKEND="${ETOPDIR}/backend"
source "${BACKEND}/scripts/lib/init.sh"

go install github.com/gogo/protobuf/protoc-gen-gogo
go install github.com/twitchtv/twirp/protoc-gen-twirp

apiDir="${BACKEND}/up/etop.vn/api"
imports="\
    -I${BACKEND}/up \
    -I$(::get mod path github.com/gogo/protobuf) \
    -I$(::get mod path github.com/grpc-ecosystem/grpc-gateway) \
    -I$(::get mod path github.com/grpc-ecosystem/grpc-gateway)/third_party/googleapis \
    "

pbFiles=""
twFiles=""
packages="$(find ${apiDir} -type d | grep -v __)"
for pkg in ${packages} ; do
  ::clean "${pkg}/*.d.go"
  ::clean "${pkg}/*.pb.go"
  ::clean "${pkg}/*.twirp.go"

	protos="${pkg}/*.proto"
	if ls ${protos} 1>/dev/null 2>/dev/null ; then
		protoc ${imports} \
		    --gogo_out="${BACKEND}/up" \
		    --twirp_out="${BACKEND}/up" \
		    ${pkg}/*.proto

		# store all generated .pb.go file names for processing later
		pbFiles="$pbFiles ${pkg}/*.pb.go"

		if ls ${pkg}/*.twirp.go 1>/dev/null 2>/dev/null ; then
		  twFiles="$twFiles ${pkg}/*.twirp.go"
		fi
	fi
done

# move all type definitions in .pb.go to .def.go
splitter=$(::get cmd etop.vn/backend/tools/cmd/split-pbgo-def)
if [[ -n "${pbFiles}" ]] ; then
  "${splitter}" $pbFiles
fi

# generate commands & queries
packagesPath=""
for pkg in ${packages}; do
  pkgPath=${pkg#"$BACKEND/up/"}
  if ls ${pkg}/*.go    1>/dev/null 2>/dev/null && \
   ! ls ${pkg}/*.proto 1>/dev/null 2>/dev/null ; then
    packagesPath="$packagesPath $pkgPath"
  fi
done

gen_cmd_query=$(::get cmd etop.vn/backend/tools/cmd/gen-cmd-query)
"${gen_cmd_query}" $packagesPath
