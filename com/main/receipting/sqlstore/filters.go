package sqlstore

import (
	"etop.vn/backend/pkg/common/sq"
	"etop.vn/backend/pkg/etop/sqlstore"
)

func (ft *ReceiptFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

func (ft *ReceiptLineFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

var SortReceipt = map[string]string{
	"id":         "",
	"created_at": "",
	"updated_at": "",
	"title":      "",
}

var FilterReceipt = sqlstore.FilterWhitelist{}

var SortReceiptLine = map[string]string{
	"id":         "",
	"created_at": "",
	"updated_at": "",
	"title":      "",
}

var FilterReceiptLine = sqlstore.FilterWhitelist{}
