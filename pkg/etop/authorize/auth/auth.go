package auth

import (
	"strings"

	"github.com/casbin/casbin"

	"etop.vn/backend/pkg/common/authorization/casbin/adapter"
)

const (
	model = `
	[request_definition]
	r = action, role
	
	[policy_definition]
	p = action, role
	
	[role_definition]
	g = _, _
	
	[policy_effect]
	e = some(where (p.eft == allow))
	
	[matchers]
	m = g(r.role, p.role) && r.action == p.action`

	policy = `
	# carrier
    p, shop/carrier:create, admin, owner, salesman
	p, shop/carrier:update, admin, owner, salesman
	p, shop/carrier:delete, admin, owner, salesman
	p, shop/carrier:view, admin, owner, salesman, accountant
	# category
	p, shop/category:create, admin, owner, salesman
	p, shop/category:update, admin, owner, salesman, purchasing_management
    p, shop/category:delete, admin, owner, salesman
    p, shop/category:view, admin, owner, salesman, accountant, purchasing_management, inventory_management
	# collection
	p, shop/collection:create, admin, owner, salesman, purchasing_management
	p, shop/collection:update, admin, owner, salesman, purchasing_management
    p, shop/collection:view, admin, owner, salesman, accountant, purchasing_management, inventory_management
	# customer
	p, shop/customer:create, admin, owner, salesman
	p, shop/customer:update, admin, owner, salesman
	p, shop/customer:delete, admin, owner, salesman
	p, shop/customer:view, admin, owner, salesman, accountant
	p, shop/customer:manage, admin, owner, salesman
	p, shop/customer_group:manage, admin, owner, salesman
    # dashboard
    p, shop/dashboard:view, admin, owner, salesman, accountant, purchasing_management, inventory_management, staff_management
	# external account
    p, shop/external_account:manage, admin, owner
    # fulfillment
	p, shop/fulfillment:create, admin, owner, salesman
	p, shop/fulfillment:cancel, admin, owner, salesman
	p, shop/fulfillment:view, admin, owner, salesman, accountant
	p, shop/fulfillment:export, admin, owner, salesman
	# inventory
	p, shop/inventory:view, admin, owner, inventory_management, purchasing_management
	p, shop/inventory:update, admin, owner, inventory_management, purchasing_management
    p, shop/inventory:create, admin, owner, salesman, inventory_management, purchasing_management
    p, shop/inventory:confirm, admin, owner, inventory_management, salesman
    p, shop/inventory:cancel, admin, owner, inventory_management
	# ledger
	p, shop/ledger:create, admin, owner, accountant
	p, shop/ledger:view, admin, owner, accountant, salesman, purchasing_management
	p, shop/ledger:update, admin, owner, accountant
	p, shop/ledger:delete, admin, owner, accountant
	# money_transaction
	p, shop/money_transaction:view, admin, owner, accountant
	p, shop/money_transaction:export, admin, owner
	# order
	p, shop/order:create, admin, owner, salesman
	p, shop/order:confirm, admin, owner, salesman
	p, shop/order:update, admin, owner, salesman
	p, shop/order:cancel, admin, owner, salesman
    # order: hotfix, all users can access order page
	p, shop/order:view, admin, owner, salesman, accountant, purchasing_management, inventory_management, staff_management
	# purchase_order
	p, shop/purchase_order:view, admin, owner, accountant, purchasing_management
	p, shop/purchase_order:create, admin, owner, purchasing_management
    p, shop/purchase_order:update, admin, owner, purchasing_management
	p, shop/purchase_order:confirm, admin, owner, purchasing_management
    p, shop/purchase_order:cancel, admin, owner, purchasing_management
	# product
	p, shop/product:create, admin, owner, purchasing_management
	p, shop/product/basic_info:update, admin, owner, purchasing_management
	p, shop/product:delete, admin, owner, purchasing_management
	p, shop/product:import, admin, owner, purchasing_management
	p, shop/product/basic_info:view, admin, owner, salesman, purchasing_management, inventory_management
	p, shop/product/cost_price:view, admin, owner, purchasing_management
	p, shop/product/cost_price:update, admin, owner, purchasing_management
	p, shop/product/retail_price:update, admin, owner
	# receipt
	p, shop/receipt:create, admin, owner, salesman, accountant, purchasing_management
	p, shop/receipt:update, admin, owner, accountant
	p, shop/receipt:confirm, admin, owner, salesman, accountant, purchasing_management
	p, shop/receipt:cancel, admin, owner, accountant
	p, shop/receipt:view, admin, owner, accountant
	# settings
	p, shop/settings/shop_info:update, admin, owner
	p, shop/settings/shop_info:view, admin, owner, salesman, accountant, purchasing_management, inventory_management, staff_management
	p, shop/settings/company_info:view, admin, owner
	p, shop/settings/company_info:update, admin, owner
	p, shop/settings/bank_info:view, admin, owner
	p, shop/settings/bank_info:update, admin, owner
	p, shop/settings/payment_schedule:view, admin, owner
	p, shop/settings/payment_schedule:update, admin, owner
	p, shop/settings/topship_balance:view, admin, owner
	p, shop/settings/wallet_balance:view, admin, owner
	p, shop/settings/shipping_setting:view, admin, owner
	p, shop/settings/shipping_setting:update, admin, owner
	# shipnow FFM
	p, shop/shipnow:view, admin, owner, salesman, accountant
	p, shop/shipnow:create, admin, owner, salesman
	p, shop/shipnow:update, admin, owner, salesman
    p, shop/shipnow:cancel, admin, owner, salesman
	p, shop/shipnow:confirm, admin, owner, salesman
	# shop
	p, shop/balance:view, admin, owner, accountant
	p, shop/account:delete, owner
	# staff
	p, relationship/invitation:create, admin, owner, staff_management
	p, relationship/invitation:view, admin, owner, staff_management
	p, relationship/invitation:delete, admin, owner, staff_management
	p, relationship/permission:update, admin, owner, staff_management
	p, relationship/relationship:update, admin, owner, staff_management
	p, relationship/relationship:view, admin, owner, staff_management
	p, relationship/user:remove, admin, owner, staff_management
	# stocktake
	p, shop/stocktake:create, admin, owner, purchasing_management, inventory_management
	p, shop/stocktake:update, admin, owner, inventory_management
	p, shop/stocktake:confirm, admin, owner, purchasing_management, inventory_management
	p, shop/stocktake:cancel, admin, owner, inventory_management
	p, shop/stocktake:view, admin, owner, accountant, inventory_management
	# supplier
	p, shop/supplier:create, admin, owner, purchasing_management
	p, shop/supplier:update, admin, owner, purchasing_management
	p, shop/supplier:delete, admin, owner, purchasing_management
	p, shop/supplier:view, admin, owner, accountant, purchasing_management
	# refund
	p, shop/refund:create, admin, owner, salesman
	p, shop/refund:update, admin, owner, salesman
	p, shop/refund:delete, admin, owner, salesman
	p, shop/refund:view, admin, owner, salesman
	# trading
	p, trading/order:view, admin, owner
	p, trading/order:create, admin, owner
    p, shop/trading/product:view, admin, owner
    p, shop/trading/order:create, admin, owner
    p, shop/trading/order:view, admin, owner	
	# payment
	p, shop/payment:create, admin, owner
	p, shop/payment:view, admin, owner`
)

var (
	mode string
)

type Authorization struct {
	*casbin.Enforcer
}

type AuthorizationService interface {
	Check(roles []string, action string) bool
}

func New() *Authorization {
	sa := adapter.NewAdapter(policy)
	enforcer := casbin.NewEnforcer(casbin.NewModel(model), sa)
	return &Authorization{enforcer}
}

func SetMode(arg string) {
	mode = arg
}

func (a *Authorization) Check(roles []string, actionsArgs string, isTest int) bool {
	switch mode {
	case "":
		return true
	case "test":
		if isTest <= 0 {
			return true
		}
	case "all":
	// no-op
	default:
		//no-op
	}

	actions := strings.Split(actionsArgs, "|")
	for _, role := range roles {
		for _, action := range actions {
			if a.Enforcer.Enforce(action, role) {
				return true
			}
		}
	}
	return false
}

func ListActionsByRoles(roles []string) (actions []string) {
	strs := strings.Split(policy, "\n")

	mapRoleAndActions := make(map[string][]string)
	for _, str := range strs {
		// prefix '#' for comment
		str = strings.TrimSpace(str)
		if str == "" || strings.HasPrefix(str, "#") {
			continue
		}

		elements := strings.Split(str, ",")
		_ = elements[0] // prefix 'p' | 'g'
		action := strings.TrimSpace(elements[1])
		roles := elements[2:]
		for _, role := range roles {
			role = strings.TrimSpace(role)
			mapRoleAndActions[role] = append(mapRoleAndActions[role], action)
		}
	}

	mapAction := make(map[string]bool)
	for _, role := range roles {
		for _, action := range mapRoleAndActions[role] {
			mapAction[action] = true
		}
	}

	for action := range mapAction {
		actions = append(actions, action)
	}
	return
}
