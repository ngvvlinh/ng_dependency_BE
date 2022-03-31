package main

import (
	shippingtypes "o.o/api/main/shipping/types"
	"o.o/api/top/types/etc/shipping"
	"o.o/api/top/types/etc/shipping_fee_type"
	"o.o/backend/cmd/etop-server/config"
	shippingconvert "o.o/backend/com/main/shipping/convert"
	shippingmodel "o.o/backend/com/main/shipping/model"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/common/l"
)

var (
	cfg config.Config
	db  *cmsql.Database
	ll  = l.New()
)

func main() {
	var err error
	if cfg, err = config.Load(false); err != nil {
		ll.Fatal("Error while loading config", l.Error(err))
	}

	postgres := cfg.Databases.Postgres
	postgres.Host = "127.0.0.1"
	postgres.Port = 5432
	postgres.Username = "postgres"
	postgres.Password = "postgres"
	postgres.Database = "etopv1.12"
	if db, err = cmsql.Connect(postgres); err != nil {
		ll.Fatal("Error while connecting database", l.Error(err))
	}

	ffmIDs := []string{"ETO1588107", "ETO1587524", "ETO1578599", "ETO1578535", "457557213", "450781092", "357405538", "338309749", "233704933", "295933589", "DJA9L3FR", "DHAKYRY6", "DHN5N56U", "ETO1569087", "802024215169", "802024138215", "SNP623050928", "802023570006", "ETO1560989", "802022519147", "802021207094", "802021007738", "802020022451", "802019224182", "802019138203", "802018915111", "801000642695", "802018680109", "TOPS1550052", "802017229901", "851244400", "802015784883", "802015784572", "802015081718", "802013965602", "690845916", "802013553863", "802013189141", "ETO1538591", "802011927467", "5220128482033960", "ETO1512140", "ETO1489000", "ETO1478925", "TOPS1462650", "EL649US7X", "ELY6QSQ3N", "EL7LRS7U5", "ELNFX1AK5", "EL91RQK1D", "ELD6ADKRU", "ELX4KNUHX", "EL5D56669", "EL5D735D5", "ELLR6L66", "ELLR6L66", "ETO1121081", "ETO1120994", "329999268", "802020957063", "801000647077", "802015636722", "802015198657", "801000542614", "802011152426", "DIXDUN54", "DISAY733", "DIX4736H", "DHFNXAXK7", "DHDHKULF", "DHQK114H"}
	ffms, err := scanFulfillmentsByIDs(ffmIDs)
	if err != nil {
		ll.Fatal("Error while list fulfillments")
	}

	updatedCount, errCount := 0, 0
	count := len(ffmIDs)
	for _, ffm := range ffms {
		// chuyển đơn về trạng thái `Đã trả hàng`
		// và returned fee = 0
		returnFeeLine := &shippingtypes.ShippingFeeLine{
			ShippingFeeType: shipping_fee_type.Return,
			Cost:            0,
		}

		shippingFeeShopLines := shippingconvert.Convert_sharemodel_ShippingFeeLines_shippingtypes_ShippingFeeLines(ffm.ShippingFeeShopLines)
		shippingFeeLinesUpdate := shippingtypes.ApplyShippingFeeLine(shippingFeeShopLines, returnFeeLine)

		totalFee := shippingtypes.GetTotalShippingFee(shippingFeeLinesUpdate)
		shippingFeeShop := shippingmodel.CalcShopShippingFee(totalFee, ffm)

		err := db.
			Table("fulfillment").
			Where("id = ?", ffm.ID).
			ShouldUpdate(&shippingmodel.Fulfillment{
				ShippingFeeShopLines: shippingconvert.Convert_shippingtypes_ShippingFeeLines_sharemodel_ShippingFeeLines(shippingFeeLinesUpdate),
				ShippingState:        shipping.Returned,
				ShippingFeeShop:      shippingFeeShop,
			})
		if err != nil {
			ll.Error("Update fulfillment error", l.String("shipping_code", ffm.ShippingCode), l.Error(err))
			errCount++
			continue
		}
		updatedCount++
	}

	ll.S.Infof("Update fulfillment shipping fee line: success %v/%v, error %v/%v", updatedCount, count, errCount, count)
}

func scanFulfillmentsByIDs(shippingCodes []string) (ffms shippingmodel.Fulfillments, err error) {
	err = db.
		In("shipping_code", shippingCodes).
		Find(&ffms)
	return
}
