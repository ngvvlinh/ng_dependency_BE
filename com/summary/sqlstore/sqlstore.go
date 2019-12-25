package sqlstore

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/backend/pkg/common/sql/sq"
	"etop.vn/backend/pkg/common/sql/sq/core"
	. "etop.vn/backend/pkg/etop/logic/summary"
	"etop.vn/capi/dot"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

type SummaryStoreFactory func(context.Context) *SummaryStore

func NewSummaryStoreFactory(db *cmsql.Database) SummaryStoreFactory {
	return func(ctx context.Context) *SummaryStore {
		return &SummaryStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

var _ = selTotal(&Total{})

type Total struct {
	TotalAmount  int     `sel:"SUM(total_amount)"`
	TotalOrder   int     `sel:"COUNT(id)"`
	AverageOrder float64 `sel:"AVG(total_amount)"`
}

func (s *SummaryStore) GetOrderSummary(shopID dot.ID, dateFrom time.Time, dateTo time.Time) (*Total, error) {
	var total Total
	q := s.query().Where("shop_id = ?", shopID)
	q = q.SQL("FROM \"order\" ")
	q = q.Where("created_at BETWEEN ? AND ?", dateFrom, dateTo).Where("status != -1")
	_, err := q.Get(&total)
	return &total, err
}

var _ = selTotalPerDate(&TotalPerDate{})

type TotalPerDate struct {
	Day         time.Time
	TotalAmount int
	Count       int
	StartTime   time.Time
	EndTime     time.Time
}

func listTotalByDate(dateFrom time.Time, dateTo time.Time) []*TotalPerDate {
	var result []*TotalPerDate
	var timeStart = dateFrom
	for timeStart.Before(dateTo) {
		result = append(result, &TotalPerDate{
			Day:         timeStart,
			TotalAmount: 0,
			Count:       0,
			StartTime:   timeStart,
			EndTime:     timeStart.Add(24 * time.Hour),
		})
		timeStart = timeStart.Add(24 * time.Hour)
	}
	return result
}

func (s *SummaryStore) buildSqlTotalPerDate(args []*TotalPerDate) *SummaryQueryBuilder {
	builder := NewSummaryQueryBuilder("order")
	index := 0
	for key, row := range args {
		index++
		// must always use [i] because we want to take the address
		builder.AddCell(&Subject{
			Expr:  " Sum(total_amount) ",
			Label: "Day :" + strconv.Itoa(key),
			Pred: &Predicate{
				Label: "Day :" + strconv.Itoa(key),
				Spec:  "Day :" + strconv.Itoa(key),
				Expr:  sq.NewExpr("created_at >= ? AND created_at < ?", row.StartTime, row.EndTime),
			},
		}, (*core.Int)(&row.TotalAmount))
		builder.AddCell(&Subject{
			Expr:  " count(id)",
			Label: "Day :" + strconv.Itoa(key),
			Pred: &Predicate{
				Label: "Day :" + strconv.Itoa(key),
				Spec:  "Day :" + strconv.Itoa(key),
				Expr:  sq.NewExpr("created_at >= ? AND created_at < ?", row.StartTime, row.EndTime),
			},
		}, (*core.Int)(&row.Count))
	}
	return builder
}

func (s *SummaryStore) GetAmoumtPerDay(shopID dot.ID, dateFrom time.Time, dateTo time.Time) ([]*TotalPerDate, error) {
	totalPerDates := listTotalByDate(dateFrom, dateTo)
	builder := s.buildSqlTotalPerDate(totalPerDates)
	q := s.query().Where("shop_id = ?", shopID).Where("status!=?", -1)
	err := q.SQL(builder).Scan(builder.ScanArgs...)
	if err != nil {
		return nil, err
	}
	return totalPerDates, nil
}

var _ = selTopSellItem(&TopSellItem{})

type TopSellItem struct {
	ProductCode string   `sel:"sp.code"`
	ProductId   dot.ID   `sel:"ol.product_id"`
	Name        string   `sel:"sp.name"`
	Count       int      `sel:"SUM(quantity) as sum"`
	ImageUrls   []string `sel:"sp.image_urls"`
}

func (s *SummaryStore) GetTopSellItem(shopID dot.ID, dateFrom time.Time, dateTo time.Time) (TopSellItems, error) {
	var topItem TopSellItems

	q := s.query().SQL("FROM order_line ol, \"order\" o, shop_product sp")
	q = q.Where("o.created_at BETWEEN ? AND ?", dateFrom, dateTo)
	q = q.Where("ol.order_id = o.id and sp.product_id = ol.product_id")
	q = q.Where("o.status != -1").Where("o.shop_id = ?", shopID)
	q = q.GroupBy("sp.code, ol.product_id, sp.name, sp.image_urls").Limit(10).OrderBy("sum desc")
	err := q.Find(&topItem)
	return topItem, err
}

var _ = selStaffOrder(&StaffOrder{})

type StaffOrder struct {
	UserName        string `sel:"u.full_name"`
	UserID          dot.ID `sel:"u.id"`
	TotalCount      int32  `sel:"count(o.id) as total_amount"`
	TotalAmount     int32  `sel:"sum(o.total_amount) as order_count"`
	OrderSourceID   dot.ID `sel:"o.order_source_id"`
	OrderSourceType string `sel:"o.order_source_type"`
}

func (s *SummaryStore) GetListStaffOrder(shopID dot.ID, dateFrom time.Time, dateTo time.Time) (StaffOrders, error) {
	var result StaffOrders

	sqlStr := `from "order" o, "order_line" ol, "user" u`
	q := s.query().SQL(sqlStr)
	q = q.Where("o.created_at BETWEEN ? AND ?", dateFrom, dateTo)
	q = q.Where(`o.status != ?`, -1)
	q = q.Where(`o.shop_id = ?`, shopID)
	q = q.Where(`o.id = ol.order_id`)
	q = q.Where(`o.created_by = u.id`)
	q = q.GroupBy(`u.id, u.full_name, o.order_source_id, o.order_source_type`).OrderBy("total_amount desc")
	err := q.Find(&result)
	return result, err
}

var _ = selFfmByArea(&FfmByArea{})

type FfmByArea struct {
	Count        int    `sel:"count(id)"`
	ProvinceCode string `sel:"address_to_province_code"`
	DistrictCode string `sel:"address_to_district_code"`
}

var sqlCountFfmByArea = strings.ReplaceAll(strings.ReplaceAll(`
		from
			%v
		where
			 shop_id = ?
			and shipping_state != 'cancelled'
			and created_at between ? and ?
		group by
			address_to_province_code,
			address_to_district_code
		order by
			address_to_district_code

		`, "\n", " "), "\t", " ")

func (s *SummaryStore) GetCountFfmByArea(shopID dot.ID, dateFrom time.Time, dateTo time.Time) ([]*FfmByArea, error) {
	var ffm FfmByAreas
	qFfm := s.query().SQL(fmt.Sprintf(sqlCountFfmByArea, "fulfillment"), shopID, dateFrom, dateTo)
	err := qFfm.Find(&ffm)
	if err != nil {
		return nil, err
	}

	var shipnowFfm FfmByAreas
	qShipnowFfm := s.query().SQL(fmt.Sprintf(sqlCountFfmByArea, "shipnow_fulfillment"), shopID, dateFrom, dateTo)
	err = qShipnowFfm.Find(&shipnowFfm)
	if err != nil {
		return nil, err
	}
	var keyShipnowFfm = 0
	for keyFfm := range ffm {
		if keyShipnowFfm > len(shipnowFfm)-1 {
			break
		}
		if shipnowFfm[keyShipnowFfm].DistrictCode == ffm[keyFfm].DistrictCode {
			ffm[keyFfm].Count = ffm[keyFfm].Count + shipnowFfm[keyShipnowFfm].Count
			keyFfm++
		}
		for shipnowFfm[keyShipnowFfm].DistrictCode > ffm[keyFfm].DistrictCode {
			ffm = append(ffm, shipnowFfm[keyFfm])
			keyFfm++
		}
	}
	sort.Slice(ffm, func(i, j int) bool {
		return ffm[i].Count > ffm[j].Count
	})
	return ffm, nil
}
