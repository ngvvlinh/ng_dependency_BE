module etop.vn/backend

go 1.12

replace etop.vn/api => ./up/etop.vn/api

replace etop.vn/capi => ./up/etop.vn/capi

replace etop.vn/common => ./up/etop.vn/common

require (
	etop.vn/api v0.0.0-00010101000000-000000000000
	etop.vn/capi v0.0.0-00010101000000-000000000000
	etop.vn/common v0.0.0-00010101000000-000000000000
	github.com/360EntSecGroup-Skylar/excelize v1.4.1
	github.com/DataDog/zstd v1.3.5 // indirect
	github.com/GoogleCloudPlatform/cloudsql-proxy v0.0.0-20190725230627-253d1edd4416
	github.com/PuerkitoBio/goquery v1.5.0
	github.com/Shopify/sarama v1.20.1
	github.com/Shopify/toxiproxy v2.1.4+incompatible // indirect
	github.com/asaskevich/govalidator v0.0.0-20180720115003-f9ffefc3facf
	github.com/awalterschulze/goderive v0.0.0-20190728081913-2613afbe1240
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/dpapathanasiou/go-recaptcha v0.0.0-20190121160230-be5090b17804
	github.com/dustin/go-humanize v1.0.0
	github.com/eapache/go-resiliency v1.1.0 // indirect
	github.com/eapache/go-xerial-snappy v0.0.0-20180814174437-776d5712da21 // indirect
	github.com/eapache/queue v1.1.0 // indirect
	github.com/garyburd/redigo v1.6.0
	github.com/go-telegram-bot-api/telegram-bot-api v4.6.4+incompatible
	github.com/go-test/deep v1.0.1
	github.com/go-xorm/builder v0.3.4
	github.com/gogo/protobuf v1.2.2-0.20190723190241-65acae22fc9d
	github.com/golang/protobuf v1.3.2
	github.com/golang/snappy v0.0.0-20180518054509-2e65f85255db // indirect
	github.com/gopherjs/gopherjs v0.0.0-20190328170749-bb2674552d8f // indirect
	github.com/gorilla/schema v1.0.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.0
	github.com/grpc-ecosystem/grpc-gateway v1.9.6
	github.com/julienschmidt/httprouter v1.2.0
	github.com/lib/pq v1.2.0
	github.com/mongodb/mongo-go-driver v0.2.0
	github.com/pierrec/lz4 v2.0.5+incompatible // indirect
	github.com/pquerna/ffjson v0.0.0-20181028064349-e517b90714f7
	github.com/prometheus/client_golang v1.0.0
	github.com/rcrowley/go-metrics v0.0.0-20181016184325-3113b8401b8a // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/smartystreets/assertions v0.0.0-20190401211740-f487f9de1cd3 // indirect
	github.com/smartystreets/goconvey v0.0.0-20190731233626-505e41936337
	github.com/spf13/pflag v1.0.3
	github.com/stretchr/testify v1.4.0
	github.com/technoweenie/multipartstreamer v1.0.1 // indirect
	github.com/tidwall/pretty v0.0.0-20180105212114-65a9db5fad51 // indirect
	github.com/twitchtv/twirp v5.8.0+incompatible
	github.com/valyala/tsvreader v1.0.0
	github.com/xdg/scram v0.0.0-20180814205039-7eeb5667e42c // indirect
	github.com/xdg/stringprep v1.0.0 // indirect
	go.uber.org/atomic v1.4.0
	go.uber.org/zap v1.10.0
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	golang.org/x/sync v0.0.0-20190423024810-112230192c58
	golang.org/x/text v0.3.2
	golang.org/x/tools v0.0.0-20190816200558-6889da9d5479
	google.golang.org/api v0.8.0 // indirect
	google.golang.org/genproto v0.0.0-20190502173448-54afdca5d873
	google.golang.org/grpc v1.22.0
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/resty.v1 v1.12.0
	gopkg.in/yaml.v2 v2.2.2
	k8s.io/apimachinery v0.0.0-20190814100815-533d101be9a6
	k8s.io/code-generator v0.0.0-20190814140513-6483f25b1faf
	k8s.io/klog v0.4.0
)
