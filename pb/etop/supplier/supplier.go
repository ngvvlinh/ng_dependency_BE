package supplier

import (
	"etop.vn/backend/pb/common"
	"etop.vn/backend/pb/etop"
	"etop.vn/backend/pb/etop/etc/status3"
	"etop.vn/backend/pb/etop/etc/status5"
	"etop.vn/backend/pb/etop/order"
	"etop.vn/backend/pb/etop/order/source"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/etop/model"
	ordermodel "etop.vn/backend/pkg/services/ordering/model"
)

func PbBranches(bs []*model.KiotvietBranch) []*KiotvietBranch {
	res := make([]*KiotvietBranch, len(bs))
	for i, b := range bs {
		res[i] = PbBranch(b)
	}
	return res
}

func PbBranch(m *model.KiotvietBranch) *KiotvietBranch {
	return &KiotvietBranch{
		Id:      m.ID,
		Name:    m.Name,
		Code:    m.Code,
		Address: m.Address,
	}
}

func PbBrand(m *model.ProductBrand) *Brand {
	return &Brand{
		Id:          m.ID,
		SupplierId:  m.SupplierID,
		Name:        m.Name,
		Description: m.Description,
		Policy:      m.Policy,
		ImageUrls:   m.ImageURLs,
		UpdatedAt:   common.PbTime(m.UpdatedAt),
		CreatedAt:   common.PbTime(m.CreatedAt),
	}
}

func PbBrands(cs []*model.ProductBrand) []*Brand {
	res := make([]*Brand, len(cs))
	for i, c := range cs {
		res[i] = PbBrand(c)
	}
	return res
}

func PbBrandExt(m *model.ProductBrandExtended) *Brand {
	pb := PbBrand(m.ProductBrand)
	pb.SupplierName = m.Supplier.Name
	return pb
}

func PbBrandsExt(cs []*model.ProductBrandExtended) []*Brand {
	res := make([]*Brand, len(cs))
	for i, c := range cs {
		res[i] = PbBrandExt(c)
	}
	return res
}

func PbCategories(cs []*model.ProductSourceCategoryExtended) []*Category {
	res := make([]*Category, len(cs))
	for i, c := range cs {
		res[i] = PbCategory(c)
	}
	return res
}

func PbCategory(m *model.ProductSourceCategoryExtended) *Category {
	return &Category{
		Id:                m.ProductSourceCategory.ID,
		Name:              m.Name,
		ProductSourceId:   m.ProductSourceCategory.ProductSourceID,
		ProductSourceType: m.ProductSourceCategory.ProductSourceType,
		ParentId:          m.ParentID,
		ShopId:            m.ShopID,
		SupplierId:        m.SupplierID,
		XId:               m.ExternalID,
		XName:             m.ExternalName,
		XParentId:         m.ExternalParentID,
	}
}

func PbVariants(items []*model.VariantExtended) []*Variant {
	res := make([]*Variant, len(items))
	for i, item := range items {
		res[i] = PbVariant(item)
	}
	return res
}

func PbVariant(m *model.VariantExtended) *Variant {
	// units := make([]*Unit, len(m.ExternalUnits))
	// for i, u := range m.ExternalUnits {
	// 	units[i] = &Unit{
	// 		Code:     u.Code,
	// 		Name:     u.Name,
	// 		FullName: u.FullName,
	// 		Unit:     u.Unit,
	// 		UnitConv: float32(u.UnitConv),
	// 		Price:    int32(u.Price),
	// 	}
	// }

	return &Variant{
		Id:         m.ID,
		SupplierId: m.SupplierID,
		CategoryId: m.Product.EtopCategoryID,

		// SCategoryId:  m.ProductSourceCategoryExtendedID, // unused

		// SMeta: common.RawJSONObjectMsg(m.SupplierMeta),
		SName:        m.EdName,
		SShortDesc:   m.EdShortDesc,
		SDescription: m.EdDescription,
		SDescHtml:    m.EdDescHTML,
		ImageUrls:    m.ImageURLs,

		XId:          m.VariantExternal.ExternalID,
		XBaseId:      m.VariantExternal.ExternalProductID,
		XCategoryId:  m.VariantExternal.ExternalCategoryID,
		XCode:        m.VariantExternal.ExternalCode,
		XName:        m.VariantExternal.ExternalName,
		XDescription: m.VariantExternal.ExternalDescription,
		XImageUrls:   m.VariantExternal.ExternalImageURLs,
		XUnit:        m.VariantExternal.ExternalUnit,
		XUnitConv:    float32(m.VariantExternal.ExternalUnitConv),
		XPrice:       int32(m.VariantExternal.ExternalPrice),
		// XUnits:      units,
		XAttributes: PbAttributes(m.VariantExternal.ExternalAttributes),

		// XCreatedAt: etop.PbTime(m.ExternalCreatedAt),
		XUpdatedAt: common.PbTime(m.VariantExternal.ExternalUpdatedAt),
		// XSyncAt:    common.PbTime(m.LastSyncAt),
		UpdatedAt: common.PbTime(m.UpdatedAt),
		CreatedAt: common.PbTime(m.CreatedAt),
		SStatus:   status3.Pb(m.EdStatus),
		XStatus:   status3.Pb(m.VariantExternal.ExternalStatus),
		EStatus:   status3.Pb(m.EtopStatus),
		Status:    status3.Pb(m.Status),

		QuantityAvailable: int32(m.QuantityAvailable),
		QuantityOnHand:    int32(m.QuantityOnHand),
		QuantityReserved:  int32(m.QuantityReserved),

		SWholesalePrice: int32(m.EdWholesalePrice),
		SListPrice:      int32(m.EdListPrice),
		SRetailPriceMin: int32(m.EdRetailPriceMin),
		SRetailPriceMax: int32(m.EdRetailPriceMax),

		WholesalePrice_0: int32(m.WholesalePrice0),
		WholesalePrice:   int32(m.WholesalePrice),
		ListPrice:        int32(m.ListPrice),
		RetailPriceMin:   int32(m.RetailPriceMin),
		RetailPriceMax:   int32(m.RetailPriceMax),
		SBrandId:         int64(m.ProductBrandID),
	}
}

func PbProducts(items []*model.ProductFtVariant) []*Product {
	res := make([]*Product, len(items))
	for i, item := range items {
		res[i] = PbProduct(item)
	}
	return res
}

func PbProduct(m *model.ProductFtVariant) *Product {
	return &Product{
		Id:         m.Product.ID,
		SupplierId: m.Product.SupplierID,
		CategoryId: m.EtopCategoryID,

		SName:        cm.Coalesce(m.EdName, m.Product.Name),
		SDescription: cm.Coalesce(m.EdDescription, m.Product.Description),
		SShortDesc:   cm.Coalesce(m.EdShortDesc, m.Product.ShortDesc),
		SDescHtml:    cm.Coalesce(m.EdDescHTML, m.Product.DescHTML),
		STags:        m.EdTags,
		ImageUrls:    m.Product.ImageURLs,

		// XId:          m.Product.ExternalID,
		XCategoryId:  m.ExternalCategoryID,
		XCode:        m.ExternalCode,
		XName:        m.ExternalName,
		XDescription: m.ExternalDescription,
		XImageUrls:   m.ExternalImageURLs,
		XUnit:        m.ExternalUnit,
		// XUnits:      units,
		// XAttributes: PbAttributes(m.ExternalAttributes),

		XCreatedAt: common.PbTime(m.ExternalCreatedAt),
		XUpdatedAt: common.PbTime(m.ExternalUpdatedAt),
		// XSyncAt:    common.PbTime(m.LastSyncAt),
		UpdatedAt: common.PbTime(m.Product.UpdatedAt),
		CreatedAt: common.PbTime(m.Product.CreatedAt),
		// SStatus:    status3.Pb(m.SupplierStatus),
		XStatus: status3.Pb(m.ExternalStatus),
		// EStatus:    status3.Pb(m.EtopStatus),
		Status: status3.Pb(m.Product.Status),

		QuantityAvailable: int32(m.QuantityAvailable),
		QuantityOnHand:    int32(m.QuantityOnHand),
		QuantityReserved:  int32(m.QuantityReserved),

		SBrandId: int64(m.ProductBrandID),
		Brand:    PbBrand(m.ProductBrand),
		Variants: PbVariantFromExternals(m.Variants, m.Product),
	}
}

func PbVariantFromExternals(items []*model.VariantExternalExtended, p *model.Product) []*Variant {
	res := make([]*Variant, len(items))
	for i, item := range items {
		res[i] = PbVariantFromExternal(item, p)
	}
	return res
}

func PbVariantFromExternal(m *model.VariantExternalExtended, p *model.Product) *Variant {
	return &Variant{
		Id:         m.Variant.ID,
		SupplierId: m.SupplierID,
		CategoryId: p.EtopCategoryID,

		// SCategoryId: p.ProductSourceCategoryID, // unused

		// SMeta:        common.RawJSONObjectMsg(m.SupplierMeta),
		SName:        m.Variant.GetName(),
		SDescription: cm.Coalesce(m.EdDescription, m.Variant.Description),
		SShortDesc:   cm.Coalesce(m.EdShortDesc, m.Variant.ShortDesc),
		SDescHtml:    cm.Coalesce(m.EdDescHTML, m.Variant.DescHTML),
		ImageUrls:    m.ImageURLs,

		XId: m.ExternalID,
		// XBaseId:      m.ExternalBaseID,
		XCategoryId:  m.ExternalCategoryID,
		XCode:        m.ExternalCode,
		XName:        m.ExternalName,
		XDescription: m.ExternalDescription,
		XImageUrls:   m.ExternalImageURLs,
		XUnit:        m.ExternalUnit,
		XUnitConv:    float32(m.ExternalUnitConv),
		XPrice:       int32(m.ExternalPrice),
		// XUnits:      units,
		XAttributes: PbAttributes(m.ExternalAttributes),

		// XCreatedAt: common.PbTime(m.ExternalCreatedAt),
		XUpdatedAt: common.PbTime(m.ExternalUpdatedAt),
		// XSyncAt:    common.PbTime(m.LastSyncAt),
		UpdatedAt: common.PbTime(m.UpdatedAt),
		CreatedAt: common.PbTime(m.CreatedAt),
		SStatus:   status3.Pb(m.EdStatus),
		XStatus:   status3.Pb(m.ExternalStatus),
		EStatus:   status3.Pb(m.EtopStatus),
		Status:    status3.Pb(m.Status),

		QuantityAvailable: int32(m.QuantityAvailable),
		QuantityOnHand:    int32(m.QuantityOnHand),
		QuantityReserved:  int32(m.QuantityReserved),

		SWholesalePrice: int32(m.EdWholesalePrice),
		SListPrice:      int32(m.EdListPrice),
		SRetailPriceMin: int32(m.EdRetailPriceMin),
		SRetailPriceMax: int32(m.EdRetailPriceMax),

		WholesalePrice_0: int32(m.WholesalePrice0),
		WholesalePrice:   int32(m.WholesalePrice),
		ListPrice:        int32(m.ListPrice),
		RetailPriceMin:   int32(m.RetailPriceMin),
		RetailPriceMax:   int32(m.RetailPriceMax),
		SBrandId:         int64(m.ProductBrandID),
		Attributes:       PbAttributes(m.Attributes),
	}
}

func PbAttributes(as []model.ProductAttribute) []*Attribute {
	attrs := make([]*Attribute, len(as))
	for i, a := range as {
		attrs[i] = &Attribute{
			Name:  a.Name,
			Value: a.Value,
		}
	}
	return attrs
}

func (a *Attribute) ToModel() model.ProductAttribute {
	if a == nil {
		return model.ProductAttribute{}
	}
	return model.ProductAttribute{
		Name:  a.Name,
		Value: a.Value,
	}
}

func AttributesTomodel(items []*Attribute) []model.ProductAttribute {
	result := make([]model.ProductAttribute, 0, len(items))
	for _, item := range items {
		if item.Name == "" {
			continue
		}
		result = append(result, item.ToModel())
	}
	return result
}

func PbUpdateVariantsToModel(items []*UpdateVariantRequest) []*model.Variant {
	res := make([]*model.Variant, len(items))
	for i, item := range items {
		res[i] = PbUpdateVariantToModel(item)
	}
	return res
}

func PbUpdateVariantToModel(p *UpdateVariantRequest) *model.Variant {
	res := &model.Variant{
		ID: p.Id,

		EdName:        p.SName,
		EdShortDesc:   p.SShortDesc,
		EdDescription: p.SDescription,
		EdDescHTML:    p.SDescHtml,
		// SupplierNote:        p.SNote,
		SupplierMeta: p.SMeta.GetData(),
		// SupplierWholesalePrice: int(p.SWholesalePrice),
		// SupplierListPrice:      int(p.SListPrice),
		// SupplierRetailPriceMin: int(p.SRetailPriceMin),
		// SupplierRetailPriceMax: int(p.SRetailPriceMax),
		ProductBrandID: int64(p.SBrandId),
	}
	return res
}

func PbPriceRules(m *model.SupplierPriceRules, ok bool) *PriceRulesResponse {
	rules := make([]*PriceRule, len(m.Rules))
	for i, item := range m.Rules {
		rules[i] = PbPriceRule(item)
	}
	return &PriceRulesResponse{
		General:  PbPriceRule(m.General),
		Rules:    rules,
		IsConfig: ok,
	}
}

func PbPriceRule(m *model.SupplierPriceRule) *PriceRule {
	return &PriceRule{
		SCategoryId:     m.SupplierCategoryID,
		XCategoryId:     m.ExternalCategoryID,
		STag:            m.Tag,
		WholesalePriceA: m.WholesalePriceA,
		WholesalePriceB: m.WholesalePriceB,
		ListPriceA:      m.ListPriceA,
		ListPriceB:      m.ListPriceB,
		RetailPriceMinA: m.RetailPriceMinA,
		RetailPriceMinB: m.RetailPriceMinB,
		RetailPriceMaxA: m.RetailPriceMaxA,
		RetailPriceMaxB: m.RetailPriceMaxB,
	}
}

func PbPriceRuleToModel(p *PriceRule) *model.SupplierPriceRule {
	return &model.SupplierPriceRule{
		SupplierCategoryID: p.SCategoryId,
		ExternalCategoryID: p.XCategoryId,
		Tag:                p.STag,
		ListPriceA:         p.ListPriceA,
		ListPriceB:         p.ListPriceB,
		WholesalePriceA:    p.WholesalePriceA,
		WholesalePriceB:    p.WholesalePriceB,
		RetailPriceMinA:    p.RetailPriceMinA,
		RetailPriceMinB:    p.RetailPriceMinB,
		RetailPriceMaxA:    p.RetailPriceMaxA,
		RetailPriceMaxB:    p.RetailPriceMaxB,
	}
}

func PbOrders(items []*ordermodel.Order) []*SupplierOrder {
	res := make([]*SupplierOrder, len(items))
	for i, item := range items {
		res[i] = PbOrder(item, nil)
	}
	return res
}

func PbOrder(m *ordermodel.Order, shop *model.ShopExtended) *SupplierOrder {
	totalItems, basketValue := 0, 0
	for _, line := range m.Lines {
		totalItems += line.Quantity
		basketValue += line.ListPrice * line.Quantity
	}
	// TotalAmount is the same as BasketValue because there is no discount
	// between suppliers and merchants.
	totalAmount := basketValue

	return &SupplierOrder{
		Id:     m.ID,
		ShopId: m.ShopID,
		Code:   m.Code,
		Source: source.PbSource(m.OrderSourceType),
		Shop:   PbShopInfo(shop),

		PaymentMethod:   m.PaymentMethod,
		Customer:        order.PbOrderCustomer(m.Customer),
		ShippingAddress: order.PbOrderAddress(m.ShippingAddress),

		CreatedAt:    common.PbTime(m.CreatedAt),
		UpdatedAt:    common.PbTime(m.UpdatedAt),
		ProcessedAt:  common.PbTime(m.ProcessedAt),
		ClosedAt:     common.PbTime(m.ClosedAt),
		CancelledAt:  common.PbTime(m.CancelledAt),
		CancelReason: m.CancelReason,

		Lines:       PbOrderLines(m.Lines),
		TotalItems:  int32(totalItems),
		BasketValue: int32(basketValue),
		TotalAmount: int32(totalAmount),

		ShConfirm:                 status3.Pb(m.ShopConfirm),
		Confirm:                   status3.Pb(m.ConfirmStatus),
		Status:                    status5.Pb(m.Status),
		FulfillmentShippingStatus: status5.Pb(m.FulfillmentShippingStatus),
		CustomerPaymentStatus:     status3.Pb(m.CustomerPaymentStatus),
		// Fulfillments:              order.PbFulfillments(m.Fulfillments, 0),
		ExternalData: order.PbOrderExternal(m.ExternalData),
	}
}

func PbOrderLines(items []*ordermodel.OrderLine) []*SupplierOrderLine {
	res := make([]*SupplierOrderLine, len(items))
	for i, item := range items {
		res[i] = PbOrderLine(item)
	}
	return res
}

func PbOrderLine(m *ordermodel.OrderLine) *SupplierOrderLine {
	return &SupplierOrderLine{
		OrderId:   m.OrderID,
		VariantId: m.VariantID,
		// VariantName: m.VariantName,

		UpdatedAt:    common.PbTime(m.UpdatedAt),
		ClosedAt:     common.PbTime(m.ClosedAt),
		ConfirmedAt:  common.PbTime(m.ConfirmedAt),
		CancelledAt:  common.PbTime(m.CancelledAt),
		CancelReason: m.CancelReason,
		SConfirm:     status3.Pb(m.SupplierConfirm),

		Quantity:   int32(m.Quantity),
		ListPrice:  int32(m.ListPrice),
		ImageUrl:   m.ImageURL,
		Attributes: PbAttributes(m.Attributes),
	}
}

func PbShopInfo(m *model.ShopExtended) *SupplierOrderShopInfo {
	if m == nil {
		return nil
	}
	return &SupplierOrderShopInfo{
		Id:         m.ID,
		Name:       m.Name,
		Phone:      m.Phone,
		Email:      m.Email,
		WebsiteUrl: m.WebsiteURL,
		ImageUrl:   m.ImageURL,
		Address: &order.OrderAddress{
			FullName:  m.Address.FullName,
			FirstName: m.Address.FirstName,
			LastName:  m.Address.LastName,
			Phone:     m.Address.Phone,
			Country:   "VN",
			City:      m.Address.City,
			Province:  m.Address.Province,
			District:  m.Address.District,
			Ward:      m.Address.Ward,
			Zip:       m.Address.Zip,
			Company:   "",
			Address1:  m.Address.Address1,
			Address2:  m.Address.Address1,
		},
	}
}

func PbContactPersonsToModel(items []*etop.ContactPerson) []*model.ContactPerson {
	res := make([]*model.ContactPerson, len(items))
	for i, item := range items {
		res[i] = PbContactPersonToModel(item)
	}
	return res
}

func PbContactPersonToModel(c *etop.ContactPerson) *model.ContactPerson {
	return &model.ContactPerson{
		Name:     c.Name,
		Phone:    c.Phone,
		Email:    c.Email,
		Position: c.Position,
	}
}

func PbSupplierExt(s *model.SupplierExtended) *etop.Supplier {
	return &etop.Supplier{
		Id:               s.Supplier.ID,
		Name:             s.Supplier.Name,
		Status:           status3.Pb(s.Supplier.Status),
		IsTest:           s.Supplier.IsTest > 0,
		CompanyInfo:      etop.PbCompanyInfo(s.Supplier.CompanyInfo),
		WarehouseAddress: etop.PbAddress(s.Address),
		BankAccount:      etop.PbBankAccount(s.Supplier.BankAccount),
		ContactPersons:   PbContactPersons(s.Supplier.ContactPersons),
	}
}

func PbContactPerson(r *model.ContactPerson) *etop.ContactPerson {
	if r == nil {
		return nil
	}
	return &etop.ContactPerson{
		Name:     r.Name,
		Position: r.Position,
		Phone:    r.Phone,
		Email:    r.Email,
	}
}

func PbContactPersons(items []*model.ContactPerson) []*etop.ContactPerson {
	if len(items) == 0 {
		return nil
	}
	res := make([]*etop.ContactPerson, len(items))
	for i, item := range items {
		res[i] = PbContactPerson(item)
	}
	return res
}

func PbUpdateProductToModel(p *UpdateProductRequest) *model.Product {
	res := &model.Product{
		ID:            p.Id,
		EdName:        p.SName,
		EdShortDesc:   p.SShortDesc,
		EdDescription: p.SDescription,
		EdDescHTML:    p.SDescHtml,
		// SupplierNote:        p.SNote,
		// SupplierMeta: p.SMeta.GetData(),
		// SupplierWholesalePrice: int(p.SWholesalePrice),
		// SupplierListPrice:      int(p.SListPrice),
		// SupplierRetailPriceMin: int(p.SRetailPriceMin),
		// SupplierRetailPriceMax: int(p.SRetailPriceMax),
		ProductBrandID: int64(p.SBrandId),
	}
	return res
}
