package shop

import (
	"testing"

	"github.com/stretchr/testify/require"

	"etop.vn/api/main/ordering"
	"etop.vn/api/main/receipting"
	"etop.vn/api/shopping/tradering"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
)

func TestValidateReceiptLines(t *testing.T) {
	t.Run("Create order valid", func(t *testing.T) {
		mockBus := bus.New()
		mockBus.MockHandler(func(query *tradering.GetTraderByIDQuery) error {
			query.Result = &tradering.ShopTrader{
				ID:   1098696748569336868,
				Type: tradering.CustomerType,
			}
			return nil
		})
		traderQuery = tradering.NewQueryBus(mockBus)

		mockBus.MockHandler(func(query *receipting.GetReceiptByCodeQuery) error {
			return cm.Errorf(cm.NotFound, nil, "receipt not found")
		})
		mockBus.MockHandler(func(query *receipting.ListReceiptsByOrderIDsQuery) error {
			query.Result = &receipting.ReceiptsResponse{
				Receipts: []*receipting.Receipt{},
			}
			return nil
		})
		receiptQuery = receipting.NewQueryBus(mockBus)

		mockBus.MockHandler(func(query *ordering.GetOrdersQuery) error {
			query.Result = &ordering.OrdersResponse{
				Orders: []*ordering.Order{
					&ordering.Order{
						ID:          1096051083316051903,
						TotalAmount: 14200,
					},
					&ordering.Order{
						ID:          1096052621074709928,
						TotalAmount: 10000,
					},
				}}
			return nil
		})
		mockBus.MockHandler(func(query *ordering.GetOrdersByIDsAndCustomerIDQuery) error {
			query.Result = &ordering.OrdersResponse{
				Orders: []*ordering.Order{
					&ordering.Order{
						ID:          1096051083316051903,
						TotalAmount: 14200,
					},
					&ordering.Order{
						ID:          1096052621074709928,
						TotalAmount: 10000,
					},
				}}
			return nil
		})
		orderQuery = ordering.NewQueryBus(mockBus)

		receipt := &receipting.Receipt{
			TraderID:    1098696748569336868,
			Code:        "code1",
			Title:       "Receipt 1",
			Type:        "receipt",
			Description: "receipt 1",
			Amount:      1100,
			Lines: []*receipting.ReceiptLine{
				&receipting.ReceiptLine{
					OrderID: 1096051083316051903,
					Title:   "receiptLine 1",
					Amount:  600,
				},
				&receipting.ReceiptLine{
					OrderID: 1096052621074709928,
					Title:   "receiptLine 2",
					Amount:  500,
				},
			},
		}
		err := receiptService.validateReceiptForCreateOrUpdate(bus.Ctx(), 0, receipt)
		require.NoError(t, err)
	})
	t.Run("No OrderID", func(t *testing.T) {
		mockBus := bus.New()
		mockBus.MockHandler(func(*ordering.GetOrdersQuery) error {
			panic("unexpected")
		})
		orderQuery = ordering.NewQueryBus(mockBus)
		lines := []*receipting.ReceiptLine{
			{
				OrderID:        0,
				Amount:         100000,
				ReceivedAmount: 100000,
			},
			{
				OrderID:        0,
				Amount:         200000,
				ReceivedAmount: 200000,
			},
		}
		receipt := &receipting.Receipt{
			Lines:  lines,
			Amount: 300000,
		}
		err := receiptService.validateReceiptLines(bus.Ctx(), tradering.CustomerType, receipt)
		require.NoError(t, err)
	})

	t.Run("Duplicated OrderID (error)", func(t *testing.T) {
		mockBus := bus.New()

		mockBus.MockHandler(func(*ordering.GetOrdersQuery) error {
			panic("unexpected")
		})
		orderQuery = ordering.NewQueryBus(mockBus)
		lines := []*receipting.ReceiptLine{
			{
				OrderID:        123456,
				Amount:         100000,
				ReceivedAmount: 100000,
			},
			{
				OrderID:        123456,
				Amount:         200000,
				ReceivedAmount: 200000,
			},
		}
		receipt := &receipting.Receipt{
			Lines:  lines,
			Amount: 300000,
		}
		err := receiptService.validateReceiptLines(bus.Ctx(), tradering.CustomerType, receipt)
		require.EqualError(t, err, "Duplicated OrderId 123456 in receipt")
	})

	t.Run("OrderID does not exist (error)", func(t *testing.T) {
		mockBus := bus.New()
		mockBus.MockHandler(func(query *ordering.GetOrdersByIDsAndCustomerIDQuery) error {
			query.Result = &ordering.OrdersResponse{
				Orders: []*ordering.Order{
					&ordering.Order{
						ID: 1001,
					},
				}}
			return nil
		})
		orderQuery = ordering.NewQueryBus(mockBus)
		lines := []*receipting.ReceiptLine{
			{
				OrderID:        1001,
				Amount:         100000,
				ReceivedAmount: 100000,
			},
			{
				OrderID:        1002,
				Amount:         200000,
				ReceivedAmount: 200000,
			},
		}
		receipt := &receipting.Receipt{
			Lines:  lines,
			Amount: 300000,
		}
		err := receiptService.validateReceiptLines(bus.Ctx(), tradering.CustomerType, receipt)
		require.EqualError(t, err, "OrderID 1002 not found")
	})

	t.Run("Amount > TotalAmount - ReceivedAmount (orderID)", func(t *testing.T) {
		// TODO: ...
	})
}
