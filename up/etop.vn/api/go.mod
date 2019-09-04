module etop.vn/api

go 1.13

replace etop.vn/capi => ../capi

replace etop.vn/common => ../common

require (
	etop.vn/capi v0.0.0-00010101000000-000000000000
	etop.vn/common v0.0.0-00010101000000-000000000000
	github.com/golang/protobuf v1.3.2
	github.com/kr/pretty v0.1.0 // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/stretchr/testify v1.4.0 // indirect
	github.com/twitchtv/twirp v5.8.0+incompatible // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
)
