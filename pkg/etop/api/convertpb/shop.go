package convertpb

import (
	"context"
	"net/url"

	"o.o/api/main/department"
	"o.o/api/main/location"
	"o.o/api/main/transaction"
	"o.o/api/shopping/addressing"
	"o.o/api/shopping/customering"
	"o.o/api/summary"
	apishop "o.o/api/top/int/shop"
	shoptypes "o.o/api/top/int/shop/types"
	"o.o/api/top/int/types"
	identitymodel "o.o/backend/com/main/identity/model"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etop/model"
)

func PbSummaryTables(tables []*model.SummaryTable) []*apishop.SummaryTable {
	res := make([]*apishop.SummaryTable, len(tables))
	for i, table := range tables {
		res[i] = &apishop.SummaryTable{
			Label:   table.Label,
			Tags:    table.Tags,
			Columns: PbSummaryColRow(table.Cols),
			Rows:    PbSummaryColRow(table.Rows),
			Data:    PbSummaryData(table.Data),
		}
	}
	return res
}

func PbSummaryColRow(items []model.SummaryColRow) []apishop.SummaryColRow {
	res := make([]apishop.SummaryColRow, len(items))
	for i, item := range items {
		res[i] = apishop.SummaryColRow{
			Label:  item.Label,
			Spec:   item.Spec,
			Unit:   item.Unit,
			Indent: item.Indent,
		}
	}
	return res
}

func PbSummaryData(data []model.SummaryItem) []apishop.SummaryItem {
	res := make([]apishop.SummaryItem, len(data))
	for i, item := range data {
		res[i] = apishop.SummaryItem{
			Spec:  item.Spec,
			Value: item.Value,
			Unit:  item.Unit,
		}
	}
	return res
}

// ObjectFrom up
func PbSummaryTablesNew(tables []*summary.SummaryTable) []*apishop.SummaryTable {
	res := make([]*apishop.SummaryTable, len(tables))
	for i, table := range tables {
		res[i] = &apishop.SummaryTable{
			Label:   table.Label,
			Tags:    table.Tags,
			Columns: PbSummaryColRowNew(table.Cols),
			Rows:    PbSummaryColRowNew(table.Rows),
			Data:    PbSummaryDataNew(table.Data),
		}
	}
	return res
}

func PbSummaryColRowNew(items []summary.SummaryColRow) []apishop.SummaryColRow {
	res := make([]apishop.SummaryColRow, len(items))
	for i, item := range items {
		res[i] = apishop.SummaryColRow{
			Label:  item.Label,
			Spec:   item.Spec,
			Unit:   item.Unit,
			Indent: item.Indent,
		}
	}
	return res
}

func PbSummaryDataNew(data []summary.SummaryItem) []apishop.SummaryItem {
	res := make([]apishop.SummaryItem, len(data))
	for i, item := range data {
		res[i] = apishop.SummaryItem{
			Spec:      item.Spec,
			Value:     item.Value,
			Unit:      item.Unit,
			ImageUrls: item.ImageUrls,
			Label:     item.Label,
		}
	}
	return res
}

func PbExportAttempts(ms []*model.ExportAttempt) []*apishop.ExportItem {
	res := make([]*apishop.ExportItem, len(ms))
	for i, m := range ms {
		res[i] = PbExportAttempt(m)
	}
	return res
}

func PbExportAttempt(m *model.ExportAttempt) *apishop.ExportItem {
	return &apishop.ExportItem{
		Id:           m.ID,
		Filename:     m.FileName,
		ExportType:   m.ExportType,
		DownloadUrl:  m.DownloadURL,
		AccountId:    m.AccountID,
		UserId:       m.UserID,
		CreatedAt:    cmapi.PbTime(m.CreatedAt),
		DeletedAt:    cmapi.PbTime(m.DeletedAt),
		RequestQuery: m.RequestQuery,
		MimeType:     m.MimeType,
		Status:       0,
		ExportErrors: cmapi.PbErrorsFromModel(m.Errors),
		Error:        cmapi.PbError(m.GetAbortedError()),
	}
}

func PbAuthorizedPartner(item *identitymodel.Partner, s *identitymodel.Shop) *apishop.AuthorizedPartnerResponse {
	redirectUrl := ""
	if item.AvailableFromEtopConfig != nil {
		redirectUrl = item.AvailableFromEtopConfig.RedirectUrl
	}
	rUrl := GenerateRedirectAuthorizedPartnerURL(redirectUrl, s)
	return &apishop.AuthorizedPartnerResponse{
		Partner:     PbPublicAccountInfo(item),
		RedirectUrl: rUrl,
	}
}

func PbAuthorizedPartners(items []*identitymodel.Partner, s *identitymodel.Shop) []*apishop.AuthorizedPartnerResponse {
	res := make([]*apishop.AuthorizedPartnerResponse, len(items))
	for i, item := range items {
		res[i] = PbAuthorizedPartner(item, s)
	}
	return res
}

func GenerateRedirectAuthorizedPartnerURL(redirectUrl string, shop *identitymodel.Shop) string {
	u, err := url.Parse(redirectUrl)
	if err != nil {
		return ""
	}
	query := u.Query()
	query.Set("shop_id", shop.ID.String())
	query.Set("email", shop.Email)
	u.RawQuery = query.Encode()
	res, _ := url.QueryUnescape(u.String())
	return res
}

func PbCustomer(m *customering.ShopCustomer) *apishop.Customer {
	return &apishop.Customer{
		Id:        m.ID,
		ShopId:    m.ShopID,
		GroupIds:  m.GroupIDs,
		FullName:  m.FullName,
		Code:      m.Code,
		Note:      m.Note,
		Phone:     m.Phone,
		Email:     m.Email,
		Gender:    m.Gender,
		Type:      m.Type,
		Birthday:  m.Birthday,
		CreatedAt: cmapi.PbTime(m.CreatedAt),
		UpdatedAt: cmapi.PbTime(m.UpdatedAt),
		Deleted:   m.Deleted,
		Status:    m.Status,
	}
}

func PbCustopmerGroup(m *customering.ShopCustomerGroup) *apishop.CustomerGroup {
	return &apishop.CustomerGroup{
		Id:   m.ID,
		Name: m.Name,
	}
}

func PbCustomerGroups(ms []*customering.ShopCustomerGroup) []*apishop.CustomerGroup {
	res := make([]*apishop.CustomerGroup, len(ms))
	for i, m := range ms {
		res[i] = PbCustopmerGroup(m)
	}
	return res
}

func PbCustomers(ms []*customering.ShopCustomer) []*apishop.Customer {
	res := make([]*apishop.Customer, len(ms))
	for i, m := range ms {
		res[i] = PbCustomer(m)
	}
	return res
}

func PbShopAddress(ctx context.Context, in *addressing.ShopTraderAddress, locationBus location.QueryBus) (*apishop.CustomerAddress, error) {
	query := &location.GetLocationQuery{
		DistrictCode: in.DistrictCode,
		WardCode:     in.WardCode,
	}
	if err := locationBus.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	province, district, ward := query.Result.Province, query.Result.District, query.Result.Ward
	out := &apishop.CustomerAddress{
		Id:           in.ID,
		District:     "",
		DistrictCode: in.DistrictCode,
		Ward:         "",
		Company:      in.Company,
		WardCode:     in.WardCode,
		Address1:     in.Address1,
		Address2:     in.Address2,
		FullName:     in.FullName,
		Phone:        in.Phone,
		Email:        in.Email,
		Position:     in.Position,
		CustomerID:   in.TraderID,
		Coordinates:  PbCoordinates(in.Coordinates),
	}
	if ward != nil {
		out.Ward = ward.Name
	}
	if district != nil {
		out.District = district.Name
	}
	if province != nil {
		out.Province = province.Name
		out.ProvinceCode = province.Code
	}
	return out, nil
}

func PbShopAddresses(ctx context.Context, ms []*addressing.ShopTraderAddress, locationBus location.QueryBus) ([]*apishop.CustomerAddress, error) {
	var err error
	res := make([]*apishop.CustomerAddress, len(ms))
	for i, m := range ms {
		res[i], err = PbShopAddress(ctx, m, locationBus)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func Convert_core_Transaction_To_api_Transaction(in *transaction.Transaction) *types.Transaction {
	if in == nil {
		return nil
	}
	return &types.Transaction{
		Name:         in.Name,
		ID:           in.ID,
		Amount:       in.Amount,
		AccountID:    in.AccountID,
		Status:       in.Status,
		Type:         in.Type,
		Classify:     in.Classify,
		Note:         in.Note,
		ReferralType: in.ReferralType,
		ReferralIDs:  in.ReferralIDs,
		CreatedAt:    in.CreatedAt,
		UpdatedAt:    in.UpdatedAt,
	}
}

func Convert_core_Transactions_To_api_Transactions(in []*transaction.Transaction) []*types.Transaction {
	if in == nil {
		return nil
	}
	res := make([]*types.Transaction, len(in))
	for i, tran := range in {
		res[i] = Convert_core_Transaction_To_api_Transaction(tran)
	}
	return res
}

func Convert_core_Department_To_api_Department(in *department.Department) *shoptypes.Department {
	if in == nil {
		return nil
	}
	return &shoptypes.Department{
		ID:           in.ID,
		AccountID:    in.AccountID,
		Name:         in.Name,
		Description:  in.Description,
		TotalMembers: in.Count,
		CreatedAt:    in.CreatedAt,
		UpdatedAt:    in.UpdatedAt,
	}
}

func Convert_core_Departments_To_api_Departments(in []*department.Department) []*shoptypes.Department {
	if in == nil {
		return nil
	}
	res := make([]*shoptypes.Department, len(in))
	for i, department := range in {
		res[i] = Convert_core_Department_To_api_Department(department)
	}
	return res
}
