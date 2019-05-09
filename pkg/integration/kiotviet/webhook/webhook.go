package webhook

import (
	"context"
	"io/ioutil"
	"strconv"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/httpreq"
	"etop.vn/backend/pkg/common/httpx"
	"etop.vn/backend/pkg/common/l"
	"etop.vn/backend/pkg/etop/model"
)

var ll = l.New()

func Register(rt *httpx.Router) {
	rt.POST("/product.update/:source_id", ProductUpdated)
	rt.POST("/product.delete/:source_id", ProductDeleted)
	rt.POST("/stock.update/:source_id", StockUpdated)
}

type (
	String = httpreq.String
	Int    = httpreq.Int
)

type StockUpdateMessage struct {
	Notifications []struct {
		Action string   `json:"Action "`
		Data   []*Stock `json:"Data"`
	} `json:"Notifications"`
}

type Stock struct {
	ProductID   String `json:"ProductId"`
	ProductCode String `json:"ProductCode"`
	ProductName String `json:"ProductName"`
	BranchID    String `json:"BranchId"`
	BranchName  String `json:"BranchName"`
	Cost        Int    `json:"Cost"`
	OnHand      Int    `json:"OnHand"`
	Reserved    Int    `json:"Reserved"`
}

func ProductUpdated(c *httpx.Context) error {
	data, err := ioutil.ReadAll(c.Req.Body)
	ll.Info("product.update", l.String("msg", string(data)), l.Error(err))
	return nil
}

func ProductDeleted(c *httpx.Context) error {
	data, err := ioutil.ReadAll(c.Req.Body)
	ll.Info("product.delete", l.String("msg", string(data)), l.Error(err))
	return nil
}

func StockUpdated(c *httpx.Context) error {
	var msg StockUpdateMessage
	err := c.DecodeJson(&msg)
	if err != nil {
		return cm.Errorf(cm.InvalidArgument, nil, "...")
	}

	ll.Info("stock.update: Received", l.Any("msg", msg))
	if len(msg.Notifications) == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "No notifications")
	}

	updates := make([]*model.ExternalQuantity, 0, 16)
	for _, n := range msg.Notifications {
		for _, d := range n.Data {
			if d == nil {
				ll.Error("Nil data", l.Any("msg", msg))
				continue
			}
			update := &model.ExternalQuantity{
				ExternalProductID: string(d.ProductID),
				BranchID:          string(d.BranchID),
				QuantityOnHand:    int(d.OnHand),
				QuantityReserved:  int(d.Reserved),
			}
			updates = append(updates, update)
		}
	}

	sourceID := paramSourceID(c)
	cmd := &model.SyncUpdateProductsQuantityCommand{
		SourceID: sourceID,
		Updates:  updates,
	}
	ctx := context.Background()
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	return nil
}

func paramSourceID(c *httpx.Context) int64 {
	s := c.ByName("source_id")
	i, _ := strconv.Atoi(s)
	return int64(i)
}
