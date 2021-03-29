package authetop

import (
	"o.o/backend/pkg/etop/authorize/auth"
	"o.o/backend/pkg/etop/authorize/authcommon"
)

const Policy auth.Policy = authcommon.CommonPolicy + `
	# etelecom
	p, shop/extension:create, admin, owner, staff_management
	p, shop/extension:delete, admin, owner, staff_management
	p, shop/extension:view, admin, owner, analyst, salesman, accountant, purchasing_management, inventory_management, staff_management, telecom_customerservice
	p, shop/hotline:view, admin, owner, staff_management
	p, shop/calllog:view, admin, owner, staff_management, telecom_customerservice
	p, shop/calllog:create, admin, owner, staff_management, telecom_customerservice
	p, admin/hotline:create, admin
	p, admin/hotline:update, admin
	# dashboard
	p, shop/dashboard:view, admin, owner, salesman, accountant, purchasing_management, inventory_management, staff_management
	p, shop/etelecom_user_setting:view, admin, owner, salesman, accountant, purchasing_management, inventory_management, staff_management
	p, admin/etelecom_user_setting:view, admin, ad_accountant, ad_salelead, ad_sale, ad_customerservice, ad_customerservice_lead
	p, admin/etelecom_user_setting:update, admin, ad_accountant, ad_salelead, ad_sale, ad_customerservice, ad_customerservice_lead
`
