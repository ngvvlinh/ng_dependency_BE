package sqlstore

import (
	"context"
	"time"

	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/sq"
	"etop.vn/backend/pkg/common/sq/core"
	"etop.vn/capi/dot"
)

type OrderStoreFactory func(context.Context) *OrderStore

func NewOrderStore(db *cmsql.Database) OrderStoreFactory {
	return func(ctx context.Context) *OrderStore {
		return &OrderStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type M map[string]interface{}

type OrderStore struct {
	query cmsql.QueryFactory
	ft    OrderFilters
	preds []interface{}
}

func (s *OrderStore) ID(id dot.ID) *OrderStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *OrderStore) ShopID(id dot.ID) *OrderStore {
	s.preds = append(s.preds, s.ft.ByShopID(id))
	return s
}

func (s *OrderStore) OrderIDs(ids dot.ID) *OrderStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "order_id", ids))
	return s
}

var _ = selTotal(&Total{})

type Total struct {
	TotalAmount  int     `sel:"SUM(total_amount)"`
	TotalOrder   int     `sel:"COUNT(id)"`
	AverageOrder float64 `sel:"AVG(total_amount)"`
}

func (s *OrderStore) GetOrderSummary(dateFrom time.Time, dateTo time.Time) (*Total, error) {
	var total Total
	q := s.query().Where(s.preds)
	q = q.SQL("FROM \"order\" ")
	q = q.Where("created_at BETWEEN ? AND ?", dateFrom, dateTo).Where("status = 1")
	_, err := q.Get(&total)
	return &total, err
}

type TotalPerDate struct {
	Day         time.Time
	TotalAmount int
	Count       int
}

func (s *OrderStore) GetAmoumtPerDay(dateFrom time.Time, dateTo time.Time) ([]*TotalPerDate, error) {
	var totalPerDate []*TotalPerDate
	q := s.query().Where(s.preds)
	q = q.SQL("SELECT DATE(created_at), SUM(total_amount), COUNT(id) FROM  \"order\" ")
	q = q.Where("created_at BETWEEN ? AND ?", dateFrom, dateTo)
	q = q.Where("status = ?", 1)
	q = q.GroupBy("date(created_at)")
	rows, err := q.Query()

	defer rows.Close()

	oldTime := time.Date(dateFrom.Year(), dateFrom.Month(), dateFrom.Day(), 0, 0, 0, 0, time.UTC)
	for rows.Next() {
		var day time.Time
		var total int
		var count int
		err = rows.Scan(&day, &total, &count)
		if err != nil {
			return nil, err
		}
		for oldTime.Before(day) && oldTime.Before(dateTo) {
			totalPerDate = append(totalPerDate, &TotalPerDate{
				Day:         oldTime,
				TotalAmount: 0,
				Count:       0,
			})
			oldTime = oldTime.Add(24 * time.Hour)
		}
		totalPerDate = append(totalPerDate, &TotalPerDate{
			Day:         day,
			TotalAmount: total,
			Count:       count,
		})
		oldTime = day
	}
	return totalPerDate, err
}

type TopSellItem struct {
	ProductId dot.ID
	Name      string
	Count     int
	ImageUrls []string
}

func (s *OrderStore) GetTopSellItem(shopID dot.ID, dateFrom time.Time, dateTo time.Time) ([]*TopSellItem, error) {
	var topItem []*TopSellItem

	q := s.query().SQL("SELECT shop_product.image_urls , shop_product.name, order_line.product_id, SUM(quantity) as sum FROM order_line, \"order\", shop_product ")
	q = q.Where("\"order\".created_at BETWEEN ? AND ?", dateFrom, dateTo)
	q = q.Where("order_line.order_id = \"order\".id and shop_product.product_id = order_line.product_id")
	q = q.Where("\"order\".status = ?", 1).Where("\"order\".shop_id = ?", shopID)
	q = q.GroupBy("order_line.product_id, shop_product.name, shop_product.image_urls").Limit(10).OrderBy("sum desc")
	rows, err := q.Query()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var productID dot.ID
		var name string
		var quantity int
		var imageUrls []string
		var imageUrlsValue = core.Array{V: &imageUrls, Opts: core.Opts{UseArrayInsteadOfJSON: true}}
		err = rows.Scan(imageUrlsValue, &name, &productID, &quantity)
		if err != nil {
			return nil, err
		}
		topItem = append(topItem, &TopSellItem{
			ProductId: productID,
			Name:      name,
			Count:     quantity,
			ImageUrls: imageUrls,
		})
	}
	return topItem, nil
}
