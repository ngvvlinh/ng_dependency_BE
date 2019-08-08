module etop.vn/api

go 1.12

replace etop.vn/capi => ../capi

require (
	etop.vn/capi v0.0.0-00010101000000-000000000000
	github.com/gogo/protobuf v1.2.1
	github.com/golang/protobuf v1.3.2
	github.com/satori/go.uuid v1.2.0
)
