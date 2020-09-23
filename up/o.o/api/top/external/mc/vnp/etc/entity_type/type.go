package entity_type

import (
	"o.o/api/top/types/etc/entity_type"
)

// +enum
// +enum:zero=null
type EntityType int

type NullEntityType struct {
	Enum  EntityType
	Valid bool
}

const (
	// +enum=shipnow_fulfillment
	ShipnowFulfillment EntityType = 53
)

func Convert_type_VnpEntities_To_type_Entities(vnpEntityTypes []EntityType) ([]entity_type.EntityType, bool) {
	var res []entity_type.EntityType
	for _, _type := range vnpEntityTypes {
		entityType, ok := entity_type.ParseEntityType(_type.String())
		if !ok {
			return nil, ok
		}
		res = append(res, entityType)
	}
	return res, true
}

func Convert_type_Entities_To_type_VnpEntities(entityTypes []entity_type.EntityType) ([]EntityType, bool) {
	var res []EntityType
	for _, _type := range entityTypes {
		vnpEntityType, ok := ParseEntityType(_type.String())
		if !ok {
			return nil, ok
		}
		res = append(res, vnpEntityType)
	}
	return res, true
}
