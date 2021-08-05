package authetop

import (
	"o.o/backend/pkg/etop/authorize/auth"
	"o.o/backend/pkg/etop/authorize/authcommon"
)

const Policy auth.Policy = authcommon.CommonPolicy + `
	# product
	p, shop/product:create, admin, owner, purchasing_management, m_admin
	p, shop/product/basic_info:update, admin, owner, purchasing_management, m_admin
	p, shop/product:delete, admin, owner, purchasing_management, m_admin
	p, shop/product:import, admin, owner, purchasing_management, m_admin
	p, shop/product/basic_info:view, admin, owner, salesman, purchasing_management, inventory_management, m_admin
	p, shop/product/cost_price:view, admin, owner, purchasing_management, m_admin
	p, shop/product/cost_price:update, admin, owner, purchasing_management, m_admin
	p, shop/product/retail_price:update, admin, owner, m_admin
	# etelecom
	p, shop/extension:create, admin, owner, staff_management, m_admin
	p, shop/extension:delete, admin, owner, staff_management, m_admin
	p, shop/extension:update, admin, owner, staff_management, m_admin
	p, shop/extension:view, admin, owner, analyst, salesman, accountant, purchasing_management, inventory_management, staff_management, telecom_customerservice, telecom_customerservice_management, m_admin
	p, shop/hotline:view, admin, owner, staff_management, telecom_customerservice, telecom_customerservice_management, m_admin
	p, shop/calllog:view, admin, owner, staff_management, telecom_customerservice, telecom_customerservice_management, m_admin
	p, shop/calllog:create, admin, owner, staff_management, telecom_customerservice, telecom_customerservice_management, m_admin
	p, shop/tenant:create, admin, owner, staff_management, telecom_customerservice, telecom_customerservice_management, m_admin
	p, shop/tenant:view, admin, owner, staff_management, telecom_customerservice, telecom_customerservice_management, m_admin
	p, shop/call_session:delete, admin, owner, analyst, salesman, accountant, purchasing_management, inventory_management, staff_management, telecom_customerservice, telecom_customerservice_management, m_admin
	p, admin/hotline:create, admin, ad_voip
	p, admin/hotline:update, admin, ad_voip
	p, admin/hotline:view, admin, ad_voip
	p, admin/hotline:delete, admin, ad_voip
	p, admin/tenant:create, admin, ad_voip
	p, admin/tenant:update, admin, ad_voip
	p, admin/tenant:view, admin, ad_voip
	# dashboard
	p, shop/dashboard:view, admin, owner, salesman, accountant, purchasing_management, inventory_management, staff_management, m_admin
	p, shop/etelecom_user_setting:view, admin, owner, salesman, accountant, purchasing_management, inventory_management, staff_management, m_admin
	p, admin/etelecom_user_setting:view, admin, ad_accountant, ad_salelead, ad_sale, ad_customerservice, ad_customerservice_lead
	p, admin/etelecom_user_setting:update, admin, ad_accountant, ad_salelead, ad_sale, ad_customerservice, ad_customerservice_lead
`
