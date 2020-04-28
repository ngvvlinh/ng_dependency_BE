module o.o/api

go 1.14

replace (
	o.o/capi => ../capi
	o.o/common => ../common
)

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/kr/pretty v0.1.0 // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/stretchr/testify v1.3.0
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	o.o/capi v0.0.0-00010101000000-000000000000
	o.o/common v0.0.0-00010101000000-000000000000
)
