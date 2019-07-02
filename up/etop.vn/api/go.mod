module etop.vn/api

go 1.12

replace etop.vn/common => ../common

require (
	etop.vn/common v0.0.0-00010101000000-000000000000
	github.com/gogo/protobuf v1.2.2-0.20190415061611-67e450fba694
	github.com/golang/protobuf v1.3.1
	github.com/grpc-ecosystem/grpc-gateway v1.8.5
	github.com/pkg/errors v0.8.1 // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/stretchr/testify v1.3.0
	github.com/twitchtv/twirp v5.7.0+incompatible
)
