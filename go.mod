module etop.vn/backend

go 1.13

replace (
	etop.vn/api => ./up/etop.vn/api
	etop.vn/capi => ./up/etop.vn/capi
	etop.vn/common => ./up/etop.vn/common
)

require (
	etop.vn/api v0.0.0-00010101000000-000000000000
	etop.vn/capi v0.0.0-00010101000000-000000000000
	etop.vn/common v0.0.0-00010101000000-000000000000
	github.com/360EntSecGroup-Skylar/excelize v1.4.1
	github.com/GoogleCloudPlatform/cloudsql-proxy v0.0.0-20191004194446-69852d3cd6a8
	github.com/PuerkitoBio/goquery v1.5.0
	github.com/Shopify/sarama v1.23.1
	github.com/asaskevich/govalidator v0.0.0-20190424111038-f61b66f89f4a
	github.com/awalterschulze/goderive v0.0.0-20190728081913-2613afbe1240
	github.com/casbin/casbin v1.9.1
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/dpapathanasiou/go-recaptcha v0.0.0-20190121160230-be5090b17804
	github.com/dustin/go-humanize v1.0.0
	github.com/garyburd/redigo v1.6.0
	github.com/go-openapi/jsonreference v0.19.3
	github.com/go-openapi/spec v0.19.4
	github.com/go-telegram-bot-api/telegram-bot-api v4.6.4+incompatible
	github.com/go-test/deep v1.0.4
	github.com/go-xorm/builder v0.3.4
	github.com/gorilla/schema v1.1.0
	github.com/jteeuwen/go-bindata v3.0.7+incompatible
	github.com/julienschmidt/httprouter v1.3.0
	github.com/k0kubun/pp v3.0.1+incompatible
	github.com/kisielk/gotool v1.0.0 // indirect
	github.com/labstack/echo v3.3.10+incompatible
	github.com/labstack/gommon v0.3.0 // indirect
	github.com/lfittl/pg_query_go v1.0.0
	github.com/lib/pq v1.2.0
	github.com/microcosm-cc/bluemonday v1.0.2
	github.com/pkg/errors v0.8.1
	github.com/pquerna/ffjson v0.0.0-20190930134022-aa0246cd15f7
	github.com/prometheus/client_golang v1.1.0
	github.com/satori/go.uuid v1.2.0
	github.com/smartystreets/goconvey v0.0.0-20190731233626-505e41936337
	github.com/stretchr/testify v1.4.0
	github.com/technoweenie/multipartstreamer v1.0.1 // indirect
	github.com/valyala/tsvreader v1.0.0
	go.uber.org/atomic v1.4.0
	go.uber.org/zap v1.10.0
	golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550 // indirect
	golang.org/x/net v0.0.0-20200222125558-5a598a2470a0 // indirect
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e
	golang.org/x/text v0.3.2
	golang.org/x/tools v0.0.0-20190821162956-65e3620a7ae7
	gopkg.in/jcmturner/goidentity.v3 v3.0.0 // indirect
	gopkg.in/resty.v1 v1.12.0
	gopkg.in/yaml.v2 v2.2.4
)
