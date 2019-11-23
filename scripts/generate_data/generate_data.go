package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"etop.vn/backend/cmd/etop-server/config"
	"etop.vn/backend/com/main/ordering/model"
	"etop.vn/backend/com/main/ordering/modelx"
	"etop.vn/backend/pkg/common/cmsql"
	cc "etop.vn/backend/pkg/common/config"
	"etop.vn/backend/pkg/etop/sqlstore"
)

func main() {
	order := &modelx.CreateOrderCommand{
		Order: &model.Order{
			ID:     0,
			ShopID: 1017379750412823345,
			Code:   "",
			EdCode: "",
			ProductIDs: []int64{
				1096813038197039889,
			},
			VariantIDs: []int64{
				1096813038199788142,
			},
			PartnerID:     0,
			Currency:      "",
			PaymentMethod: "cod",
			Customer: &model.OrderCustomer{
				FirstName:     "",
				LastName:      "",
				FullName:      "huynh hai nam",
				Email:         "",
				Phone:         "",
				Gender:        "",
				Birthday:      "",
				VerifiedEmail: false,
				ExternalID:    "",
			},
			CustomerAddress: &model.OrderAddress{
				FullName:     "huynh hai nam",
				FirstName:    "",
				LastName:     "",
				Phone:        "0346644185",
				Email:        "",
				Country:      "",
				City:         "",
				Province:     "Thành phố Đà Nẵng",
				District:     "Quận Hải Châu",
				Ward:         "Phường Hải Châu II",
				Zip:          "",
				DistrictCode: "492",
				ProvinceCode: "48",
				WardCode:     "20239",
				Company:      "",
				Address1:     "3424",
				Address2:     "",
				Coordinates:  (*model.Coordinates)(nil),
			},
			BillingAddress: &model.OrderAddress{
				FullName:     "huynh hai nam",
				FirstName:    "",
				LastName:     "",
				Phone:        "0346644185",
				Email:        "",
				Country:      "",
				City:         "",
				Province:     "Thành phố Đà Nẵng",
				District:     "Quận Hải Châu",
				Ward:         "Phường Hải Châu II",
				Zip:          "",
				DistrictCode: "492",
				ProvinceCode: "48",
				WardCode:     "20239",
				Company:      "",
				Address1:     "3424",
				Address2:     "",
				Coordinates:  (*model.Coordinates)(nil),
			},
			ShippingAddress: &model.OrderAddress{
				FullName:     "huynh hai nam",
				FirstName:    "",
				LastName:     "",
				Phone:        "0346644185",
				Email:        "",
				Country:      "",
				City:         "",
				Province:     "Thành phố Đà Nẵng",
				District:     "Quận Hải Châu",
				Ward:         "Phường Hải Châu II",
				Zip:          "",
				DistrictCode: "492",
				ProvinceCode: "48",
				WardCode:     "20239",
				Company:      "",
				Address1:     "3424",
				Address2:     "",
				Coordinates:  (*model.Coordinates)(nil),
			},
			CustomerName:               "",
			CustomerPhone:              "",
			CustomerEmail:              "",
			CreatedAt:                  time.Now(),
			ProcessedAt:                time.Now(),
			UpdatedAt:                  time.Now(),
			ClosedAt:                   time.Now(),
			ConfirmedAt:                time.Now(),
			CancelledAt:                time.Now(),
			CancelReason:               "",
			CustomerConfirm:            0,
			ShopConfirm:                0,
			ConfirmStatus:              0,
			FulfillmentShippingStatus:  0,
			EtopPaymentStatus:          0,
			Status:                     0,
			FulfillmentShippingStates:  []string{},
			FulfillmentPaymentStatuses: []int{},
			Lines: model.OrderLinesList{
				&model.OrderLine{
					OrderID:         0,
					VariantID:       1096813038199788142,
					ProductName:     "234234",
					ProductID:       1096813038197039889,
					ShopID:          1017379750412823345,
					Weight:          0,
					Quantity:        1,
					ListPrice:       2323,
					RetailPrice:     2323,
					PaymentPrice:    2323,
					LineAmount:      2323,
					TotalDiscount:   0,
					TotalLineAmount: 2323,
					ImageURL:        "",
					IsOutsideEtop:   false,
					Code:            "",
				},
			},
			Discounts:       []*model.OrderDiscount{},
			TotalItems:      1,
			BasketValue:     2323,
			TotalWeight:     0,
			TotalTax:        0,
			OrderDiscount:   0,
			TotalDiscount:   0,
			ShopShippingFee: 0,
			TotalFee:        0,
			FeeLines: model.OrderFeeLines{
				model.OrderFeeLine{
					Amount: 0,
					Desc:   "Phí vận chuyển tính với khách",
					Code:   "",
					Name:   "Phí vận chuyển tính với khách",
					Type:   "shipping",
				},
			},
			ShopCOD:          0,
			TotalAmount:      2323,
			OrderNote:        "",
			ShopNote:         "",
			ShippingNote:     "",
			OrderSourceType:  "etop_pos",
			OrderSourceID:    0,
			ExternalOrderID:  "",
			ReferenceURL:     "",
			ExternalURL:      "",
			ShopShipping:     (*model.OrderShipping)(nil),
			IsOutsideEtop:    false,
			GhnNoteCode:      "",
			TryOn:            "",
			CustomerNameNorm: "",
			ProductNameNorm:  "",
			FulfillmentType:  1,
			FulfillmentIDs:   []int64{},
			ExternalMeta: json.RawMessage{
				0x6e, 0x75, 0x6c, 0x6c,
			},
			TradingShopID: 0,
			PaymentStatus: 0,
			PaymentID:     0,
			ReferralMeta: json.RawMessage{
				0x6e, 0x75, 0x6c, 0x6c,
			},
		},
	}
	ctx := context.Background()
	configpg := config.Config{
		Postgres: cc.Postgres{
			Protocol:        "",
			Host:            "localhost",
			Port:            5432,
			Username:        "postgres",
			Password:        "postgres",
			Database:        "etop_dev",
			SSLMode:         "",
			Timeout:         15,
			MaxOpenConns:    0,
			MaxIdleConns:    0,
			MaxConnLifetime: 0,
			GoogleAuthFile:  "",
		},
	}
	db, _ := cmsql.Connect(configpg.Postgres)
	sqlstore.Init(db)
	for i := 1; i < 1000000; i++ {
		err := sqlstore.CreateOrder(ctx, order)
		if err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println("done")
}
