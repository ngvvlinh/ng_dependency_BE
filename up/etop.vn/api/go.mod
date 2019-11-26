module etop.vn/api

go 1.13

replace (
	etop.vn/capi => ../capi
	etop.vn/common => ../common
)

require (
	etop.vn/capi v0.0.0-00010101000000-000000000000
	etop.vn/common v0.0.0-00010101000000-000000000000
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/kr/pretty v0.1.0 // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/stretchr/testify v1.3.0
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
)
