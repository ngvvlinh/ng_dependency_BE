#!/bin/bash

CURRENT_DIR="$(pwd)"
BACKEND_DIR="$PROJECT_DIR/backend"

cd $BACKEND_DIR

rm -rf fabo
mkdir fabo
mkdir fabo/doc


echo "run deps tools"
go run ./tools/cmd/deps -copy-files=./fabo ./cmd/fabo-server
go run ./tools/cmd/deps -copy-files=./fabo ./cmd/fabo-sync-service
go run ./tools/cmd/deps -copy-files=./fabo ./cmd/fabo-pgevent-forwarder
go run ./tools/cmd/deps -copy-files=./fabo ./cmd/fabo-event-handler
go run ./tools/cmd/deps -copy-files=./fabo ./cmd/etop-uploader


echo "restructure packages"
cp -R ./fabo/$BACKEND_DIR/* ./fabo

rm -rf ./fabo/Users
rm -rf ./fabo/o.o


echo "add up package"
cp -R ./up ./fabo


echo "add go mod"
cp ./go.mod ./fabo/go.mod
cp ./go.sum ./fabo/go.sum


echo "add necessary files and packages"
cp -R ./doc/fabo ./fabo/doc
cp -R ./doc/bindata.go ./fabo/doc
cp -R ./tools/pkg/gen ./fabo/tools/pkg
cp -Rf ./res/dl/fabo/* ./fabo/res/dl/fabo


echo " -> copy scripts package"
mkdir ./fabo/scripts
mkdir ./fabo/scripts/once
mkdir ./fabo/scripts/ci
mkdir ./fabo/scripts/lib

cp ./scripts/once/config-kafka.sh ./fabo/scripts/once
cp ./scripts/generate-all.sh ./fabo/scripts
cp ./scripts/ci/verify-generate.sh ./fabo/scripts/ci
cp ./scripts/install-tools.sh ./fabo/scripts
cp ./scripts/clean-imports.sh ./fabo/scripts
cp ./scripts/lib/init.sh ./fabo/scripts/lib
cp ./scripts/lib/get.sh ./fabo/scripts/lib


echo " -> copy tools package"
mkdir ./fabo/tools/cmd
mkdir ./fabo/tools/cmd/generator
cp -R ./tools/cmd/generator/main.go ./fabo/tools/cmd/generator

mkdir ./fabo/tools/pkg/generator
cp -R ./tools/pkg/generator ./fabo/tools/pkg
cp -R ./tools/pkg/generators ./fabo/tools/pkg
cp -R ./tools/pkg/genutil ./fabo/tools/pkg
cp -R ./tools/pkg/typedesc ./fabo/tools/pkg
cp -R ./tools/cmd/clean-imports ./fabo/tools/cmd


echo " -> copy missing packages"
cp -R ./com/main/connectioning/pm ./fabo/com/main/connectioning
cp -R ./com/shopping/customering/pm ./fabo/com/shopping/customering
cp -R ./com/main/receipting/aggregate ./fabo/com/main/receipting
cp -R ./com/main/receipting/pm ./fabo/com/main/receipting
cp -R ./com/main/stocktaking/aggregate ./fabo/com/main/stocktaking
cp -R ./com/main/authorization/query ./fabo/com/main/authorization
cp -R ./com/main/invitation/pm ./fabo/com/main/invitation
cp -R ./com/main/moneytx/aggregate ./fabo/com/main/moneytx
cp -R ./com/main/moneytx/pm ./fabo/com/main/moneytx
cp -R ./com/main/moneytx/query ./fabo/com/main/moneytx
cp -R ./com/main/inventory/pm ./fabo/com/main/inventory
cp -R ./com/main/moneytx/sqlstore ./fabo/com/main/moneytx
cp -R ./com/main/moneytx/convert ./fabo/com/main/moneytx
cp -R ./com/main/shipmentpricing/pricelist/pm ./fabo/com/main/shipmentpricing/pricelist
cp -R ./pkg/etop/api/shop/brand ./fabo/pkg/etop/api/shop
cp -R ./pkg/common/sql/migration ./fabo/pkg/common/sql


echo " -> copy cogs and wire files"
cp ./cmd/fabo-server/build/wire.go ./fabo/cmd/fabo-server/build
cp ./cmd/fabo-sync-service/build/wire.go ./fabo/cmd/fabo-sync-service/build
cp ./cmd/fabo-event-handler/build/wire.go ./fabo/cmd/fabo-event-handler/build
cp ./cogs/base/wire.go ./fabo/cogs/base
cp -R ./com/fabo/cogs ./fabo/com/fabo
cp ./com/main/address/wire.go ./fabo/com/main/address
cp ./com/main/authorization/wire.go ./fabo/com/main/authorization
cp ./com/main/catalog/wire.go ./fabo/com/main/catalog
cp ./com/main/connectioning/wire.go ./fabo/com/main/connectioning
cp ./com/main/identity/wire.go ./fabo/com/main/identity
cp ./com/main/invitation/wire.go ./fabo/com/main/invitation
cp ./com/main/location/wire.go ./fabo/com/main/location
cp ./com/main/moneytx/wire.go ./fabo/com/main/moneytx
cp ./com/main/ordering/wire.go ./fabo/com/main/ordering
cp ./com/main/receipting/wire.go ./fabo/com/main/receipting
cp ./com/main/shipmentpricing/wire.go ./fabo/com/main/shipmentpricing
cp ./com/main/shipnow/wire.go ./fabo/com/main/shipnow
cp ./com/main/shipping/wire.go ./fabo/com/main/shipping
cp ./com/main/shippingcode/wire.go ./fabo/com/main/shippingcode
cp ./com/main/stocktaking/wire.go ./fabo/com/main/stocktaking
cp ./com/shopping/customering/wire.go ./fabo/com/shopping/customering
cp ./com/shopping/setting/wire.go ./fabo/com/shopping/setting
cp ./com/fabo/main/fbcustomerconversationsearch/wire.go ./fabo/com/fabo/main/fbcustomerconversationsearch
cp ./com/fabo/main/fbmessagetemplate/wire.go ./fabo/com/fabo/main/fbmessagetemplate
cp ./com/fabo/main/fbmessaging/wire.go ./fabo/com/fabo/main/fbmessaging
cp ./com/fabo/main/fbpage/wire.go ./fabo/com/fabo/main/fbpage
cp ./com/fabo/main/fbuser/wire.go ./fabo/com/fabo/main/fbuser
cp ./com/eventhandler/handler/api/wire.go ./fabo/com/eventhandler/handler/api
cp ./com/eventhandler/webhook/sender/wire.go ./fabo/com/eventhandler/webhook/sender
cp ./com/eventhandler/webhook/storage/wire.go ./fabo/com/eventhandler/webhook/storage
cp ./pkg/etop/sqlstore/wire.go ./fabo/pkg/etop/sqlstore
cp ./cogs/shipment/wire.go ./fabo/cogs/shipment
cp ./cogs/shipment/_fabo/wire.go ./fabo/cogs/shipment/_fabo
cp ./cogs/shipment/ghn/v2/wire.go ./fabo/cogs/shipment/ghn/v2
mkdir ./fabo/cogs/core
cp ./cogs/core/wire.go ./fabo/cogs/core
cp ./com/summary/fabo/wire.go ./fabo/com/summary/fabo
mkdir ./fabo/pkg/etop/api/shop/_wire
mkdir ./fabo/pkg/etop/api/shop/_wire/fabo
cp ./pkg/etop/api/shop/_wire/fabo/wire.go ./fabo/pkg/etop/api/shop/_wire/fabo
cp ./cogs/database/_min/wire.go ./fabo/cogs/database/_min
cp ./cogs/sms/_min/wire.go ./fabo/cogs/sms/_min
cp ./cogs/config/_server/wire.go ./fabo/cogs/config/_server
cp ./cogs/uploader/wire.go ./fabo/cogs/uploader
cp ./cogs/server/shop/wire.go ./fabo/cogs/server/shop
cp ./cogs/server/fabo/wire.go ./fabo/cogs/server/fabo
cp ./pkg/etop/api/shop/_min/fabo/wire.go ./fabo/pkg/etop/api/shop/_min/fabo
cp ./cogs/storage/_all/wire.go ./fabo/cogs/storage/_all
cp ./pkg/integration/email/wire.go ./fabo/pkg/integration/email
cp ./pkg/etop/logic/orders/wire.go ./fabo/pkg/etop/logic/orders
cp ./pkg/etop/logic/orders/imcsv/wire.go ./fabo/pkg/etop/logic/orders/imcsv
cp ./pkg/etop/logic/products/imcsv/wire.go ./fabo/pkg/etop/logic/products/imcsv
cp ./pkg/etop/eventstream/wire.go ./fabo/pkg/etop/eventstream
cp ./com/eventhandler/notifier/wire.go ./fabo/com/eventhandler/notifier
cp ./com/main/inventory/aggregatex/wire.go ./fabo/com/main/inventory/aggregatex
cp ./pkg/etop/api/export/wire.go ./fabo/pkg/etop/api/export
cp ./pkg/etop/authorize/middleware/wire.go ./fabo/pkg/etop/authorize/middleware
cp ./pkg/etop/api/sadmin/_fabo/wire.go ./fabo/pkg/etop/api/sadmin/_fabo
cp ./pkg/common/apifw/captcha/wire.go ./fabo/pkg/common/apifw/captcha
cp ./pkg/etop/authorize/auth/wire.go ./fabo/pkg/etop/authorize/auth
cp ./com/eventhandler/handler/wire.go ./fabo/com/eventhandler/handler
cp ./com/eventhandler/fabo/publisher/wire.go ./fabo/com/eventhandler/fabo/publisher
cp ./cogs/server/fabo/wire.go ./fabo/cogs/server/fabo
cp ./pkg/fabo/wire.go ./fabo/pkg/fabo
cp ./com/main/shipmentpricing/pricelist/wire.go ./fabo/com/main/shipmentpricing/pricelist
cp ./com/main/shipmentpricing/pricelistpromotion/wire.go ./fabo/com/main/shipmentpricing/pricelistpromotion
cp ./com/main/shipmentpricing/shipmentprice/wire.go ./fabo/com/main/shipmentpricing/shipmentprice
cp ./com/main/shipmentpricing/shipmentservice/wire.go ./fabo/com/main/shipmentpricing/shipmentservice
cp ./com/main/shipmentpricing/shopshipmentpricelist/wire.go ./fabo/com/main/shipmentpricing/shopshipmentpricelist
cp ./com/main/shipping/carrier/wire.go ./fabo/com/main/shipping/carrier
cp ./com/etc/logging/shippingwebhook/wire.go ./fabo/com/etc/logging/shippingwebhook
cp ./pkg/integration/shipping/ghn/webhook/v2/wire.go ./fabo/pkg/integration/shipping/ghn/webhook/v2
cp ./cogs/sms/wire.go ./fabo/cogs/sms
cp ./com/etc/logging/smslog/wire.go ./fabo/com/etc/logging/smslog
cp ./pkg/integration/sms/wire.go ./fabo/pkg/integration/sms


echo " -> copy docker-compose.yml"
cp ./docker-compose.yml ./fabo


echo " -> replace PROJECT_DIR to FABO_PROJECT_DIR"
script_files=("./fabo/scripts/generate-all.sh" "./fabo/docker-compose.yml" "./fabo/scripts/clean-imports.sh" "./fabo/scripts/lib/init.sh" "./fabo/tools/pkg/gen/gen.go" "./fabo/scripts/lib/get.sh")
for file in "${script_files[@]}"
do
    tmp=$(sed 's/PROJECT_DIR/FABO_PROJECT_DIR/g' $file); echo "$tmp" > $file
done


mv fabo backend
mkdir fabo
mkdir fabo/o.o

mv ./backend ./fabo/o.o
cd ./fabo/o.o/backend

mv $BACKEND_DIR/fabo/ $CURRENT_DIR

echo "generate-all"
./scripts/generate-all.sh


echo "DONE"