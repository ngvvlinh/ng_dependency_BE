module etop.vn/api

go 1.12

replace etop.vn/capi => ../capi

replace etop.vn/common => ../common

require (
	etop.vn/capi v0.0.0-00010101000000-000000000000
	etop.vn/common v0.0.0-00010101000000-000000000000
	github.com/gogo/protobuf v1.2.1
	github.com/golang/protobuf v1.3.2
	github.com/grpc-ecosystem/grpc-gateway v1.9.6
	github.com/satori/go.uuid v1.2.0
	github.com/stretchr/testify v1.4.0
	github.com/twitchtv/twirp v5.8.0+incompatible
)
