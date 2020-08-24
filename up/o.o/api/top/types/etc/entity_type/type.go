package entity_type

// +enum
// +enum:zero=null
type EntityType int

type NullEntityType struct {
	Enum  EntityType
	Valid bool
}

const (
	// +enum=order
	Order EntityType = 11

	// +enum=fulfillment
	Fulfillment EntityType = 13

	// +enum=product
	Product EntityType = 19

	// +enum=variant
	Variant EntityType = 23

	// +enum=customer
	Customer EntityType = 27

	// +enum=inventory_level
	InventoryLevel EntityType = 31

	// +enum=customer_address
	CustomerAddress EntityType = 35

	// +enum=customer_group
	CustomerGroup EntityType = 39

	// +enum=customer_group_relationship
	CustomerGroupRelationship EntityType = 43

	// +enum=product_collection
	ProductCollection EntityType = 47

	// +enum=product_collection_relationship
	ProductCollectionRelationship EntityType = 49

	// +enum=shipnow_fulfillment
	ShipnowFulfillment EntityType = 53
)

func Contain(list []EntityType, _type EntityType) bool {
	for _, item := range list {
		if item == _type {
			return true
		}
	}
	return false
}
