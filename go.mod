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
	github.com/DataDog/zstd v1.4.1 // indirect
	github.com/GoogleCloudPlatform/cloudsql-proxy v0.0.0-20190828224159-d93c53a4824c
	github.com/PuerkitoBio/goquery v1.5.0
	github.com/Shopify/sarama v1.23.1
	github.com/asaskevich/govalidator v0.0.0-20190424111038-f61b66f89f4a
	github.com/awalterschulze/goderive v0.0.0-20190728081913-2613afbe1240
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/dpapathanasiou/go-recaptcha v0.0.0-20190121160230-be5090b17804
	github.com/dustin/go-humanize v1.0.0
	github.com/eapache/go-resiliency v1.2.0 // indirect
	github.com/elliots/protoc-gen-twirp_swagger v0.0.0-20190707132119-fda01df7c6eb
	github.com/frankban/quicktest v1.4.2 // indirect
	github.com/garyburd/redigo v1.6.0
	github.com/go-telegram-bot-api/telegram-bot-api v4.6.4+incompatible
	github.com/go-test/deep v1.0.3
	github.com/go-xorm/builder v0.3.4
	github.com/gogo/protobuf v1.3.0
	github.com/golang/protobuf v1.3.2
	github.com/google/go-cmp v0.3.1 // indirect
	github.com/gopherjs/gopherjs v0.0.0-20190812055157-5d271430af9f // indirect
	github.com/gorilla/schema v1.1.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.0
	github.com/grpc-ecosystem/grpc-gateway v1.11.1
	github.com/hashicorp/golang-lru v0.5.3 // indirect
	github.com/jcmturner/gofork v1.0.0 // indirect
	github.com/jteeuwen/go-bindata v3.0.7+incompatible
	github.com/julienschmidt/httprouter v1.2.0
	github.com/lib/pq v1.2.0
	github.com/mattn/go-isatty v0.0.9 // indirect
	github.com/pierrec/lz4 v2.2.7+incompatible // indirect
	github.com/pquerna/ffjson v0.0.0-20190813045741-dac163c6c0a9
	github.com/prometheus/client_golang v1.1.0
	github.com/rcrowley/go-metrics v0.0.0-20190826022208-cac0b30c2563 // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/smartystreets/assertions v1.0.1 // indirect
	github.com/smartystreets/goconvey v0.0.0-20190731233626-505e41936337
	github.com/spf13/pflag v1.0.3
	github.com/stretchr/testify v1.4.0
	github.com/technoweenie/multipartstreamer v1.0.1 // indirect
	github.com/twitchtv/twirp v5.8.0+incompatible
	github.com/valyala/tsvreader v1.0.0
	go.uber.org/atomic v1.4.0
	go.uber.org/zap v1.10.0
	golang.org/x/crypto v0.0.0-20190829043050-9756ffdc2472 // indirect
	golang.org/x/net v0.0.0-20190827160401-ba9fcec4b297 // indirect
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	golang.org/x/sync v0.0.0-20190423024810-112230192c58
	golang.org/x/sys v0.0.0-20190904005037-43c01164e931 // indirect
	golang.org/x/text v0.3.2
	golang.org/x/tools v0.0.0-20190907020128-2ca718005c18
	google.golang.org/appengine v1.6.2 // indirect
	google.golang.org/genproto v0.0.0-20190905072037-92dd089d5514
	google.golang.org/grpc v1.23.0
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/jcmturner/goidentity.v3 v3.0.0 // indirect
	gopkg.in/jcmturner/gokrb5.v7 v7.3.0 // indirect
	gopkg.in/resty.v1 v1.12.0
	gopkg.in/yaml.v2 v2.2.2
	k8s.io/apimachinery v0.0.0-20190831074630-461753078381
	k8s.io/code-generator v0.0.0-20190831074504-732c9ca86353
	k8s.io/klog v0.4.0
)
