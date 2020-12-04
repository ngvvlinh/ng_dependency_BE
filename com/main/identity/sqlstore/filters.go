package sqlstore

import (
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
)

func (ft *UserFilters) NotBelongWLPartner() sq.WriterTo {
	return ft.Filter("$.wl_partner_id IS NULL")
}

func (ft *ShopFilters) NotBelongWLPartner() sq.WriterTo {
	return ft.Filter("$.wl_partner_id IS NULL")
}

func (ft ShopFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

func (ft *AccountUserFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

func (ft *PartnerRelationFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

var filterShopExtendedWhitelist = sqlstore.FilterWhitelist{
	Arrays:   nil,
	Bools:    nil,
	Dates:    []string{"created_at"},
	Equals:   []string{"owner_id", "code", "phone", "email", "money_transaction_rrule"},
	Nullable: []string{"bank_account"},
	Numbers:  nil,
	Status:   nil,
	Unaccent: nil,
	PrefixOrRename: map[string]string{
		"name":                    "s",
		"created_at":              "s",
		"code":                    "s",
		"phone":                   "s",
		"email":                   "s",
		"money_transaction_rrule": "s",
		"bank_account":            "s",
	},
}

var SortUser = map[string]string{
	"id":         "",
	"created_at": "",
	"updated_at": "",
}
