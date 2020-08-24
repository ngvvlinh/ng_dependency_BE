package model

import (
	"encoding/json"
	"strings"
	"time"

	"o.o/api/top/types/etc/entity_type"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/cmenv"
	"o.o/backend/pkg/common/validate"
	"o.o/capi/dot"
)

// +sqlgen
type Callback struct {
	ID        dot.ID
	WebhookID dot.ID
	AccountID dot.ID
	CreatedAt time.Time `sq:"create"`
	Changes   json.RawMessage
	Result    json.RawMessage // WebhookStatesError
}

// +sqlgen
type Webhook struct {
	ID        dot.ID
	AccountID dot.ID
	Entities  []entity_type.EntityType
	Fields    []string
	URL       string
	Metadata  string
	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
	DeletedAt time.Time
}

func (m *Webhook) BeforeInsert() error {
	if m == nil {
		return cm.Errorf(cm.InvalidArgument, nil, "empty data")
	}
	if m.AccountID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "missing Name")
	}
	if !validate.URL(m.URL) {
		return cm.Errorf(cm.InvalidArgument, nil, "Địa chỉ url không hợp lệ")
	}
	if cmenv.IsProd() && !strings.HasPrefix(m.URL, "https://") {
		return cm.Errorf(cm.InvalidArgument, nil, "Địa chỉ url phải là https://")
	}
	if len(m.Entities) == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "missing entity")
	}
	if len(m.Fields) > 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Thông tin fields chưa được hỗ trợ, vui lòng để trống")
	}

	mp := make(map[entity_type.EntityType]bool)
	for _, item := range m.Entities {
		if !validate.LowercaseID(item.String()) {
			return cm.Errorf(cm.InvalidArgument, nil, `invalid entity: "%v"`, item)
		}
		switch item {
		case entity_type.Order, entity_type.Fulfillment,
			entity_type.Product, entity_type.Variant, entity_type.Customer,
			entity_type.InventoryLevel, entity_type.CustomerAddress,
			entity_type.CustomerGroup, entity_type.CustomerGroupRelationship,
			entity_type.ProductCollection, entity_type.ProductCollectionRelationship,
			entity_type.ShipnowFulfillment:
			if mp[item] {
				return cm.Errorf(cm.InvalidArgument, nil, `duplicated entity: "%v"`, item)
			}
			mp[item] = true
		default:
			return cm.Errorf(cm.InvalidArgument, nil, `unknown entity: "%v"`, item)
		}
	}

	m.ID = cm.NewID()
	return nil
}
