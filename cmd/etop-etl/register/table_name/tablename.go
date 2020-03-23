package table_name

import "fmt"

// +enum
// +enum:zero=null
type TableName int

type NullTableName struct {
	Enum  TableName
	Valid bool
}

const (
	// +enum=user
	User TableName = 0

	// +enum=account
	Account TableName = 1

	// +enum=shop_customer
	ShopCustomer TableName = 2

	// +enum=order
	Order TableName = 3

	// +enum=shop
	Shop TableName = 4

	// +enum=fulfillment
	Fulfillment TableName = 5

	// +enum=shop_brand
	ShopBrand TableName = 6

	// +enum=shop_product
	ShopProduct TableName = 7
)

func ConvertStringsToTableNames(args []string) []TableName {
	tableNames := make([]TableName, 0, len(args))
	for _, arg := range args {
		tableName, ok := ParseTableName(arg)
		if !ok {
			panic(fmt.Sprintf("table name %q is invalid", arg))
		}
		tableNames = append(tableNames, tableName)
	}
	return tableNames
}
