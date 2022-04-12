module o.o/backend

go 1.15

replace (
	github.com/google/wire => ./up/github.com/google/wire
	o.o/api => ./up/o.o/api
	o.o/capi => ./up/o.o/capi
	o.o/common => ./up/o.o/common
)

require (
	cloud.google.com/go/storage v1.6.0
	github.com/360EntSecGroup-Skylar/excelize v1.4.1
	github.com/360EntSecGroup-Skylar/excelize/v2 v2.4.0
	github.com/GoogleCloudPlatform/cloudsql-proxy v1.17.0
	github.com/Luzifer/go-openssl/v4 v4.1.0
	github.com/PuerkitoBio/goquery v1.5.1
	github.com/SebastiaanKlippert/go-wkhtmltopdf v1.6.0
	github.com/Shopify/sarama v1.26.4
	github.com/andybalholm/cascadia v1.2.0 // indirect
	github.com/asaskevich/govalidator v0.0.0-20200428143746-21a406dcc535
	github.com/aws/aws-sdk-go v1.40.1
	github.com/axw/gocov v1.0.0
	github.com/casbin/casbin v1.9.1
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/dpapathanasiou/go-recaptcha v0.0.0-20190121160230-be5090b17804
	github.com/dustin/go-humanize v1.0.0
	github.com/garyburd/redigo v1.6.2
	github.com/gin-gonic/gin v1.6.3
	github.com/go-bindata/go-bindata/v3 v3.1.3
	github.com/go-openapi/jsonreference v0.19.4
	github.com/go-openapi/spec v0.19.9
	github.com/go-openapi/swag v0.19.9 // indirect
	github.com/go-playground/assert/v2 v2.0.1
	github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.0.0-rc1
	github.com/go-test/deep v1.0.7
	github.com/google/wire v0.4.1-0.20200530155014-5633f1675aa9
	github.com/gopherjs/gopherjs v0.0.0-20200217142428-fce0ec30dd00 // indirect
	github.com/gorilla/schema v1.1.0
	github.com/julienschmidt/httprouter v1.3.0
	github.com/k0kubun/pp v3.0.1+incompatible
	github.com/klauspost/compress v1.10.10 // indirect
	github.com/labstack/echo/v4 v4.1.16
	github.com/lfittl/pg_query_go v1.0.0
	github.com/lib/pq v1.8.0
	github.com/mailru/easyjson v0.7.2 // indirect
	github.com/matm/gocov-html v0.0.0-20200509184451-71874e2e203b
	github.com/mattn/go-colorable v0.1.7 // indirect
	github.com/microcosm-cc/bluemonday v1.0.3
	github.com/pierrec/lz4 v2.5.2+incompatible // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.7.1
	github.com/r3labs/sse v0.0.0-20210224172625-26fe804710bc
	github.com/rcrowley/go-metrics v0.0.0-20200313005456-10cdbea86bc0 // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/smartystreets/assertions v1.1.1 // indirect
	github.com/smartystreets/goconvey v1.6.4
	github.com/stretchr/testify v1.6.1
	github.com/valyala/fasttemplate v1.2.0 // indirect
	github.com/valyala/tsvreader v1.0.0
	go.uber.org/atomic v1.6.0
	go.uber.org/zap v1.15.0
	golang.org/x/mod v0.3.0 // indirect
	golang.org/x/net v0.0.0-20210614182718-04defd469f4e
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208
	golang.org/x/text v0.3.6
	golang.org/x/tools v0.0.0-20200519205726-57a9e4404bf7
	google.golang.org/api v0.21.0
	google.golang.org/appengine v1.6.6 // indirect
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/olivere/elastic.v5 v5.0.86
	gopkg.in/resty.v1 v1.12.0
	gopkg.in/robfig/cron.v2 v2.0.0-20150107220207-be2e0b0deed5
	gopkg.in/yaml.v2 v2.3.0
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
	o.o/api v0.0.0-00010101000000-000000000000
	o.o/capi v0.0.0-00010101000000-000000000000
	o.o/common v0.0.0-00010101000000-000000000000
)
