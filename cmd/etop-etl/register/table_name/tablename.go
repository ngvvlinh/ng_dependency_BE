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

	// +enum=account_user
	AccountUser TableName = 8

	// +enum=address
	Address TableName = 9

	// +enum=inventory_variant
	InventoryVariant TableName = 10

	// +enum=inventory_voucher
	InventoryVoucher TableName = 11

	// +enum=invitation
	Invitation TableName = 12

	// +enum=money_transaction_shipping
	MoneyTransactionShipping TableName = 13

	// +enum=product_shop_collection
	ProductShopCollection TableName = 14

	// +enum=purchase_order
	PurchaseOrder TableName = 15

	// +enum=purchase_refund
	PurchaseRefund TableName = 16

	// +enum=receipt
	Receipt TableName = 17

	// +enum=refund
	Refund TableName = 18

	// +enum=shipnow_fulfillment
	ShipNowFufillment TableName = 19

	// +enum=shop_carrier
	ShopCarrier TableName = 20

	// +enum=shop_category
	ShopCategory TableName = 21

	// +enum=shop_collection
	ShopCollection TableName = 22

	// +enum=shop_customer_group
	ShopCustomerGroup TableName = 23

	// +enum=shop_customer_group_customer
	ShopCustomerGroupCustomer TableName = 24

	// +enum=shop_ledger
	ShopLedger TableName = 25

	// +enum=shop_product_collection
	ShopProductCollection TableName = 26

	// +enum=shop_stocktake
	ShopStocktake TableName = 27

	// +enum=shop_supplier
	ShopSupplier TableName = 28

	// +enum=shop_trader
	ShopTrader TableName = 29

	// +enum=shop_trader_address
	ShopTraderAddress TableName = 30

	// +enum=shop_variant
	ShopVariant TableName = 31

	// +enum=shop_variant_supplier
	ShopVariantSupplier TableName = 32
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
