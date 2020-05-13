module o.o/api

go 1.14

replace (
	o.o/capi => ../capi
	o.o/common => ../common
)

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/kr/text v0.2.0 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/stretchr/testify v1.5.1
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	o.o/capi v0.0.0-00010101000000-000000000000
	o.o/common v0.0.0-00010101000000-000000000000
)
