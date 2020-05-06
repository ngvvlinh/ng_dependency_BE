package sqlstore

import (
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
)

func (ft *FbCustomerConversationFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

var FilterFbCustomerConversation = sqlstore.FilterWhitelist{
	Dates:   []string{"created_at", "updated_at", "last_message_at"},
	Numbers: []string{"id", "fb_page_id"},
	Equals:  []string{"type", "external_id", "external_user_id"},
}

var SortFbCustomerConversation = map[string]string{
	"id":         "id",
	"created_at": "",
	"updated_at": "",
}

func (ft *FbExternalConversationFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

var FilterFbExternalConversation = sqlstore.FilterWhitelist{
	Dates:   []string{"created_at", "updated_at"},
	Numbers: []string{"id", "fb_page_id"},
	Equals:  []string{"type", "external_id"},
}

var SortFbExternalConversation = map[string]string{
	"id":         "id",
	"created_at": "",
	"updated_at": "",
}

func (ft *FbExternalMessageFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

var FilterFbExternalMessage = sqlstore.FilterWhitelist{
	Dates:   []string{"created_at", "updated_at", "external_created_time"},
	Numbers: []string{"id", "fb_page_id"},
	Equals:  []string{"type", "external_id"},
}

var SortFbExternalMessage = map[string]string{
	"id":                    "id",
	"created_at":            "",
	"updated_at":            "",
	"external_created_time": "",
}
