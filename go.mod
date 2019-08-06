module etop.vn/backend

go 1.12

replace etop.vn/api => ./up/etop.vn/api

replace etop.vn/apis => ./up/etop.vn/apis

replace etop.vn/apix => ./up/etop.vn/apix

replace etop.vn/common => ./up/etop.vn/common

require (
	cloud.google.com/go v0.35.1 // indirect
	etop.vn/api v0.0.0-00010101000000-000000000000
	etop.vn/common v0.0.0-00010101000000-000000000000
	github.com/360EntSecGroup-Skylar/excelize v1.4.1
	github.com/DataDog/zstd v1.3.5 // indirect
	github.com/GoogleCloudPlatform/cloudsql-proxy v0.0.0-20181215173202-6f1ecdcf9588
	github.com/PuerkitoBio/goquery v1.5.0
	github.com/Shopify/sarama v1.20.1
	github.com/Shopify/toxiproxy v2.1.4+incompatible // indirect
	github.com/asaskevich/govalidator v0.0.0-20180720115003-f9ffefc3facf
	github.com/awalterschulze/goderive v0.0.0-20180911165245-2b2ac3583a4c
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
	github.com/gogo/protobuf v1.2.2-0.20190415061611-67e450fba694
	github.com/golang/protobuf v1.3.1
	github.com/golang/snappy v0.0.0-20180518054509-2e65f85255db // indirect
	github.com/gopherjs/gopherjs v0.0.0-20190328170749-bb2674552d8f // indirect
	github.com/gorilla/schema v1.0.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.0
	github.com/grpc-ecosystem/grpc-gateway v1.8.5
	github.com/julienschmidt/httprouter v1.2.0
	github.com/lib/pq v1.0.0
	github.com/mongodb/mongo-go-driver v0.2.0
	github.com/pierrec/lz4 v2.0.5+incompatible // indirect
	github.com/pquerna/ffjson v0.0.0-20181028064349-e517b90714f7
	github.com/prometheus/client_golang v1.0.0
	github.com/rcrowley/go-metrics v0.0.0-20181016184325-3113b8401b8a // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/smartystreets/assertions v0.0.0-20190401211740-f487f9de1cd3 // indirect
	github.com/smartystreets/goconvey v0.0.0-20190330032615-68dc04aab96a
	github.com/spf13/pflag v1.0.3
	github.com/stretchr/testify v1.3.0
	github.com/technoweenie/multipartstreamer v1.0.1 // indirect
	github.com/tidwall/pretty v0.0.0-20180105212114-65a9db5fad51 // indirect
	github.com/twitchtv/twirp v5.7.0+incompatible
	github.com/valyala/tsvreader v1.0.0
	github.com/xdg/scram v0.0.0-20180814205039-7eeb5667e42c // indirect
	github.com/xdg/stringprep v1.0.0 // indirect
	go.uber.org/atomic v1.4.0
	go.uber.org/zap v1.10.0
	golang.org/x/crypto v0.0.0-20190325154230-a5d413f7728c // indirect
	golang.org/x/net v0.0.0-20190328230028-74de082e2cca // indirect
	golang.org/x/oauth2 v0.0.0-20190115181402-5dab4167f31c
	golang.org/x/sync v0.0.0-20190423024810-112230192c58
	golang.org/x/sys v0.0.0-20190402142545-baf5eb976a8c // indirect
	golang.org/x/text v0.3.1-0.20181227161524-e6919f6577db
	golang.org/x/tools v0.0.0-20190617190820-da514acc4774
	google.golang.org/genproto v0.0.0-20190128161407-8ac453e89fca
	google.golang.org/grpc v1.22.0
	gopkg.in/resty.v1 v1.12.0
	gopkg.in/yaml.v2 v2.2.2
	k8s.io/apimachinery v0.0.0-20190425132440-17f84483f500
	k8s.io/code-generator v0.0.0-20190419212335-ff26e7842f9d
	k8s.io/klog v0.3.0
)
