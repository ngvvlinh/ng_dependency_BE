package convertpb

import (
	"context"

	"o.o/api/main/authorization"
	"o.o/api/main/credit"
	"o.o/api/main/identity"
	identitytypes "o.o/api/main/identity/types"
	"o.o/api/main/invitation"
	"o.o/api/main/location"
	ordertypes "o.o/api/main/ordering/types"
	etop "o.o/api/top/int/etop"
	"o.o/api/top/types/etc/account_type"
	addresstype "o.o/api/top/types/etc/address_type"
	"o.o/api/top/types/etc/status3"
	notimodel "o.o/backend/com/eventhandler/notifier/model"
	addressmodel "o.o/backend/com/main/address/model"
	creditmodel "o.o/backend/com/main/credit/model"
	identitymodel "o.o/backend/com/main/identity/model"
	identitysharemodel "o.o/backend/com/main/identity/sharemodel"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/integration/bank"
	"o.o/capi/dot"
)

func CompanyInfoToModel(e *etop.CompanyInfo) *identitysharemodel.CompanyInfo {
	if e == nil {
		return nil
	}
	return &identitysharemodel.CompanyInfo{
		Name:                e.Name,
		TaxCode:             e.TaxCode,
		Address:             e.Address,
		Website:             e.Website,
		LegalRepresentative: ContactPersonToModel(e.LegalRepresentative),
	}
}

func BankAccountToModel(b *etop.BankAccount) *identitysharemodel.BankAccount {
	if b == nil {
		return nil
	}
	return &identitysharemodel.BankAccount{
		Name:          b.Name,
		Province:      b.Province,
		Branch:        b.Branch,
		AccountNumber: b.AccountNumber,
		AccountName:   b.AccountName,
	}
}

func BankAccountToCoreBankAccount(b *etop.BankAccount) *identitytypes.BankAccount {
	if b == nil {
		return nil
	}
	return &identitytypes.BankAccount{
		Name:          b.Name,
		Province:      b.Province,
		Branch:        b.Branch,
		AccountNumber: b.AccountNumber,
		AccountName:   b.AccountName,
	}
}

func AddressToModel(a *etop.Address) (*addressmodel.Address, error) {
	if a == nil {
		return nil, nil
	}
	res := &addressmodel.Address{
		ID:           a.Id,
		Province:     a.Province,
		ProvinceCode: a.ProvinceCode,
		District:     a.District,
		DistrictCode: a.DistrictCode,
		Ward:         a.Ward,
		WardCode:     a.WardCode,
		Address1:     a.Address1,
		Address2:     a.Address2,
		FullName:     a.FullName,
		FirstName:    a.FirstName,
		LastName:     a.LastName,
		Phone:        a.Phone,
		Position:     a.Position,
		Email:        a.Email,
		Notes:        PbAddressNoteToModel(a.Notes),
		Type:         a.Type.String(),
	}
	locationQuery := &location.FindOrGetLocationQuery{
		ProvinceCode: a.ProvinceCode,
		DistrictCode: a.DistrictCode,
		WardCode:     a.WardCode,
		Province:     a.Province,
		District:     a.District,
		Ward:         a.Ward,
	}
	if err := locationBus.Dispatch(context.TODO(), locationQuery); err != nil {
		return nil, err
	}
	loc := locationQuery.Result
	if loc.Province == nil || loc.District == nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Cần cung cấp thông tin tỉnh/thành phố và quận/huyện")
	}

	res.Province = loc.Province.Name
	res.ProvinceCode = loc.Province.Code
	res.District = loc.District.Name
	res.DistrictCode = loc.District.Code
	if loc.Ward != nil {
		res.Ward = loc.Ward.Name
		res.WardCode = loc.Ward.Code
	}
	if a.Coordinates != nil {
		res.Coordinates = &addressmodel.Coordinates{
			Latitude:  a.Coordinates.Latitude,
			Longitude: a.Coordinates.Longitude,
		}
	}
	return res, nil
}

func PbUser(m *identitymodel.User) *etop.User {
	if m == nil {
		panic("Nil user")
	}
	return &etop.User{
		Id:        m.ID,
		FullName:  m.FullName,
		ShortName: m.ShortName,
		Phone:     m.Phone,
		Email:     m.Email,
		CreatedAt: cmapi.PbTime(m.CreatedAt),
		UpdatedAt: cmapi.PbTime(m.UpdatedAt),

		EmailVerifiedAt: cmapi.PbTime(m.EmailVerifiedAt),
		PhoneVerifiedAt: cmapi.PbTime(m.PhoneVerifiedAt),

		EmailVerificationSentAt: cmapi.PbTime(m.EmailVerificationSentAt),
		PhoneVerificationSentAt: cmapi.PbTime(m.PhoneVerificationSentAt),
		Source:                  m.Source,
	}
}

func PbUsers(ms []*identity.User) []*etop.User {
	res := make([]*etop.User, len(ms))
	for i, m := range ms {
		res[i] = Convert_core_User_To_api_User(m)
	}
	return res
}

func PbAccountType(t account_type.AccountType) account_type.AccountType {
	return t
}

func PbLoginAccount(m *identitymodel.AccountUserExtended) *etop.LoginAccount {
	account := m.Account
	return &etop.LoginAccount{
		Id:          account.ID,
		Name:        account.Name,
		Type:        PbAccountType(account.Type),
		AccessToken: "", // Will be filled later
		ExpiresIn:   0,  // Will be filled later
		ImageUrl:    account.ImageURL,
		UrlSlug:     account.URLSlug,
		UserAccount: PbUserAccount(m),
	}
}

func PbUserAccount(m *identitymodel.AccountUserExtended) *etop.UserAccountInfo {
	account := m.Account
	accUser := m.AccountUser
	user := m.User

	fullName, shortName := m.GetUserName()
	return &etop.UserAccountInfo{
		UserId:               user.ID,
		UserFullName:         fullName,
		UserShortName:        shortName,
		AccountId:            account.ID,
		AccountName:          account.Name,
		AccountType:          PbAccountType(account.Type),
		Position:             accUser.Position,
		Permission:           PbPermission(accUser),
		Status:               accUser.Status,
		ResponseStatus:       accUser.ResponseStatus,
		InvitationSentBy:     accUser.InvitationSentBy,
		InvitationSentAt:     cmapi.PbTime(accUser.InvitationSentAt),
		InvitationAcceptedAt: cmapi.PbTime(accUser.InvitationAcceptedAt),
		DisabledAt:           cmapi.PbTime(accUser.DisabledAt),
	}
}

func PbUserAccounts(items []*identitymodel.AccountUserExtended) []*etop.UserAccountInfo {
	result := make([]*etop.UserAccountInfo, len(items))
	for i, item := range items {
		result[i] = PbUserAccount(item)
	}
	return result
}

func PbUserAccountIncomplete(accUser *identitymodel.AccountUser, account *identitymodel.Account) *etop.UserAccountInfo {
	return &etop.UserAccountInfo{
		UserId:               accUser.UserID,
		UserFullName:         "",
		UserShortName:        "",
		AccountId:            accUser.AccountID,
		AccountName:          account.Name,
		AccountType:          PbAccountType(account.Type),
		Position:             accUser.Position,
		Permission:           PbPermission(accUser),
		Status:               accUser.Status,
		InvitationSentBy:     accUser.InvitationSentBy,
		InvitationSentAt:     cmapi.PbTime(accUser.InvitationSentAt),
		InvitationAcceptedAt: cmapi.PbTime(accUser.InvitationAcceptedAt),
		DisabledAt:           cmapi.PbTime(accUser.DisabledAt),
	}
}

func PbPermission(m *identitymodel.AccountUser) *etop.Permission {
	return &etop.Permission{
		Roles:       m.Roles,
		Permissions: m.Permissions,
	}
}

func PbPartner(m *identitymodel.Partner) *etop.Partner {
	return &etop.Partner{
		Id:             m.ID,
		Name:           m.Name,
		PublicName:     m.PublicName,
		Status:         m.Status,
		IsTest:         m.IsTest != 0,
		ContactPersons: PbContactPersons(m.ContactPersons),
		Phone:          m.Phone,
		WebsiteUrl:     m.WebsiteURL,
		ImageUrl:       m.ImageURL,
		Email:          m.Email,
		OwnerId:        m.OwnerID,
		User:           nil, // TODO
	}
}

func PbPublicPartners(items []*identitymodel.Partner) []*etop.PublicAccountInfo {
	res := make([]*etop.PublicAccountInfo, len(items))
	for i, item := range items {
		res[i] = PbPublicAccountInfo(item)
	}
	return res
}

func PbPublicAccountInfo(m identitymodel.AccountInterface) *etop.PublicAccountInfo {
	switch m := m.(type) {
	case *identitymodel.Partner:
		return &etop.PublicAccountInfo{
			Id:       m.ID,
			Name:     m.PublicName, // public name here!
			Type:     PbAccountType(account_type.Partner),
			ImageUrl: m.ImageURL,
			Website:  m.WebsiteURL,
		}
	default:
		account := m.GetAccount()
		return &etop.PublicAccountInfo{
			Id:       account.ID,
			Name:     account.Name,
			Type:     PbAccountType(account.Type),
			ImageUrl: account.ImageURL,
			Website:  "",
		}
	}
}

func PbShop(m *identitymodel.Shop) *etop.Shop {
	return &etop.Shop{
		Id:          m.ID,
		Name:        m.Name,
		Status:      m.Status,
		Phone:       m.Phone,
		BankAccount: PbBankAccount(m.BankAccount),
		WebsiteUrl:  m.WebsiteURL.String,
		ImageUrl:    m.ImageURL,
		Email:       m.Email,
		OwnerId:     m.OwnerID,
		TryOn:       m.TryOn,
		Code:        m.Code,
	}
}

func Convert_core_Shop_To_api_Shop(in *identity.Shop) *etop.Shop {
	if in == nil {
		return nil
	}
	return &etop.Shop{
		Id:                    in.ID,
		Name:                  in.Name,
		Status:                in.Status,
		IsTest:                in.IsTest == 1,
		Phone:                 in.Phone,
		AutoCreateFfm:         in.AutoCreateFFM,
		WebsiteUrl:            in.WebsiteURL,
		ImageUrl:              in.ImageURL,
		Email:                 in.Email,
		ShipToAddressId:       in.ShipToAddressID,
		ShipFromAddressId:     in.ShipFromAddressID,
		OwnerId:               in.OwnerID,
		Code:                  in.Code,
		BankAccount:           Convert_core_BankAccount_To_api_BankAccount(in.BankAccount),
		TryOn:                 in.TryOn,
		MoneyTransactionRrule: in.MoneyTransactionRRule,
	}
}

func Convert_core_Shops_To_api_Shops(items []*identity.Shop) []*etop.Shop {
	result := make([]*etop.Shop, len(items))
	for i, item := range items {
		result[i] = Convert_core_Shop_To_api_Shop(item)
	}
	return result
}

func PbShopExtended(m *identitymodel.ShopExtended) *etop.Shop {
	return &etop.Shop{
		Id:                            m.ID,
		InventoryOverstock:            m.InventoryOverstock.Apply(true),
		Name:                          m.Name,
		Status:                        m.Status,
		Address:                       PbAddress(m.Address),
		Phone:                         m.Phone,
		BankAccount:                   PbBankAccount(m.BankAccount),
		WebsiteUrl:                    m.WebsiteURL.String,
		ImageUrl:                      m.ImageURL,
		Email:                         m.Email,
		ShipToAddressId:               m.ShipToAddressID,
		ShipFromAddressId:             m.ShipFromAddressID,
		AutoCreateFfm:                 m.AutoCreateFFM,
		TryOn:                         m.TryOn,
		GhnNoteCode:                   m.GhnNoteCode,
		OwnerId:                       m.OwnerID,
		User:                          PbUser(m.User),
		CompanyInfo:                   PbCompanyInfo(m.CompanyInfo),
		MoneyTransactionRrule:         m.MoneyTransactionRRule,
		SurveyInfo:                    PbSurveyInfos(m.SurveyInfo),
		ShippingServiceSelectStrategy: PbShippingServiceSelectStrategy(m.ShippingServiceSelectStrategy),
		Code:                          m.Code,

		// deprecated: 2018.07.24+14
		ProductSourceId: m.ID,
	}
}

func PbShopExtendeds(items []*identitymodel.ShopExtended) []*etop.Shop {
	result := make([]*etop.Shop, len(items))
	for i, item := range items {
		result[i] = PbShopExtended(item)
	}
	return result
}

func Convert_core_ShopExtended_To_api_ShopExtended(m *identity.ShopExtended) *etop.Shop {
	if m == nil {
		return nil
	}
	return &etop.Shop{
		Id:                            m.ID,
		InventoryOverstock:            m.InventoryOverstock.Apply(true),
		Name:                          m.Name,
		Status:                        m.Status,
		Address:                       Convert_core_Address_To_api_Address(m.Address),
		Phone:                         m.Phone,
		BankAccount:                   Convert_core_BankAccount_To_api_BankAccount(m.BankAccount),
		WebsiteUrl:                    m.WebsiteURL,
		ImageUrl:                      m.ImageURL,
		Email:                         m.Email,
		ShipToAddressId:               m.ShipToAddressID,
		ShipFromAddressId:             m.ShipFromAddressID,
		AutoCreateFfm:                 m.AutoCreateFFM,
		TryOn:                         m.TryOn,
		GhnNoteCode:                   m.GhnNoteCode,
		OwnerId:                       m.OwnerID,
		User:                          Convert_core_User_To_api_User(m.User),
		CompanyInfo:                   Convert_core_CompanyInfo_To_api_CompanyInfo(m.CompanyInfo),
		MoneyTransactionRrule:         m.MoneyTransactionRRule,
		SurveyInfo:                    Convert_core_SurveyInfos_To_api_SurveyInfors(m.SurveyInfo),
		ShippingServiceSelectStrategy: Convert_core_ShippingServiceSelectStrategy_To_api_ShippingServiceSelectStrategy(m.ShippingServiceSelectStrategy),
		Code:                          m.Code,

		// deprecated: 2018.07.24+14
		ProductSourceId: m.ID,
	}
}

func Convert_core_ShopExtendeds_To_api_ShopExtendeds(items []*identity.ShopExtended) []*etop.Shop {
	result := make([]*etop.Shop, len(items))
	for i, item := range items {
		result[i] = Convert_core_ShopExtended_To_api_ShopExtended(item)
	}
	return result
}

func Convert_core_User_To_api_User(in *identity.User) *etop.User {
	if in == nil {
		return nil
	}
	isBlocked := false
	if in.Status == status3.N {
		isBlocked = true
	}
	return &etop.User{
		Id:              in.ID,
		FullName:        in.FullName,
		ShortName:       in.ShortName,
		Phone:           in.Phone,
		Email:           in.Email,
		CreatedAt:       cmapi.PbTime(in.CreatedAt),
		UpdatedAt:       cmapi.PbTime(in.UpdatedAt),
		EmailVerifiedAt: cmapi.PbTime(in.EmailVerifiedAt),
		PhoneVerifiedAt: cmapi.PbTime(in.PhoneVerifiedAt),
		Source:          in.Source,
		IsBlocked:       isBlocked,
	}
}

func PbProvinces(items []*location.Province) []*etop.Province {
	res := make([]*etop.Province, len(items))
	for i, item := range items {
		res[i] = PbProvince(item)
	}
	return res
}

func PbProvince(m *location.Province) *etop.Province {
	return &etop.Province{
		Code:   m.Code,
		Name:   m.Name,
		Region: m.Region.Name(),
	}
}

func PbDistricts(items []*location.District) []*etop.District {
	res := make([]*etop.District, len(items))
	for i, item := range items {
		res[i] = PbDistrict(item)
	}
	return res
}

func PbDistrict(item *location.District) *etop.District {
	return &etop.District{
		Code:         item.Code,
		ProvinceCode: item.ProvinceCode,
		Name:         item.Name,
	}
}

func PbWards(items []*location.Ward) []*etop.Ward {
	res := make([]*etop.Ward, len(items))
	for i, item := range items {
		res[i] = PbWard(item)
	}
	return res
}

func PbWard(item *location.Ward) *etop.Ward {
	return &etop.Ward{
		Code:         item.Code,
		DistrictCode: item.DistrictCode,
		Name:         item.Name,
	}
}

func PbBanks(items []*bank.Bank) []*etop.Bank {
	res := make([]*etop.Bank, len(items))
	for i, item := range items {
		res[i] = PbBank(item)
	}
	return res
}

func PbBank(item *bank.Bank) *etop.Bank {
	return &etop.Bank{
		Code: item.MaNganHang,
		Name: item.TenNH,
		Type: item.Loai,
	}
}

func PbBankProvinces(items []*bank.Province) []*etop.BankProvince {
	res := make([]*etop.BankProvince, len(items))
	for i, item := range items {
		res[i] = PbBankProvince(item)
	}
	return res
}

func PbBankProvince(item *bank.Province) *etop.BankProvince {
	return &etop.BankProvince{
		Code:     item.MaTinh,
		Name:     item.TenTinhThanh,
		BankCode: item.MaNganHang,
	}
}

func PbBankBranches(items []*bank.Branch) []*etop.BankBranch {
	res := make([]*etop.BankBranch, len(items))
	for i, item := range items {
		res[i] = PbBankBranch(item)
	}
	return res
}

func PbBankBranch(item *bank.Branch) *etop.BankBranch {
	return &etop.BankBranch{
		Code:         item.MaChiNhanh,
		Name:         item.TenChiNhanh,
		BankCode:     item.MaNganHang,
		ProvinceCode: item.MaTinh,
	}
}

func PbAddresses(items []*addressmodel.Address) []*etop.Address {
	result := make([]*etop.Address, len(items))
	for i, item := range items {
		result[i] = PbAddress(item)
	}
	return result
}

func PbAddress(a *addressmodel.Address) *etop.Address {
	if a == nil {
		return nil
	}
	addressType, _ := addresstype.ParseAddressType(a.Type)
	res := &etop.Address{
		Id:           a.ID,
		Province:     a.Province,
		ProvinceCode: a.ProvinceCode,
		District:     a.District,
		DistrictCode: a.DistrictCode,
		Ward:         a.Ward,
		WardCode:     a.WardCode,
		Address1:     a.Address1,
		Address2:     a.Address2,
		Zip:          a.Zip,
		Country:      a.Country,
		FullName:     a.FullName,
		FirstName:    a.FirstName,
		LastName:     a.LastName,
		Phone:        a.Phone,
		Email:        a.Email,
		Position:     a.Position,
		Type:         addressType,
		Notes:        PbAddressNote(a.Notes),
	}
	if a.Coordinates != nil {
		res.Coordinates = &etop.Coordinates{
			Latitude:  a.Coordinates.Latitude,
			Longitude: a.Coordinates.Longitude,
		}
	}
	return res
}

func PbBankAccount(b *identitysharemodel.BankAccount) *etop.BankAccount {
	if b == nil {
		return nil
	}
	return &etop.BankAccount{
		Name:          b.Name,
		Province:      b.Province,
		Branch:        b.Branch,
		AccountName:   b.AccountName,
		AccountNumber: b.AccountNumber,
	}
}

func Convert_core_BankAccount_To_api_BankAccount(in *identitytypes.BankAccount) *etop.BankAccount {
	if in == nil {
		return nil
	}
	return &etop.BankAccount{
		Name:          in.Name,
		Province:      in.Province,
		Branch:        in.Branch,
		AccountNumber: in.AccountNumber,
		AccountName:   in.AccountName,
	}
}

func ContactPersonToModel(m *etop.ContactPerson) *identitysharemodel.ContactPerson {
	return &identitysharemodel.ContactPerson{
		Name:     m.Name,
		Position: m.Position,
		Phone:    m.Phone,
		Email:    m.Email,
	}
}

func ContactPersonsToModel(items []*etop.ContactPerson) []*identitysharemodel.ContactPerson {
	result := make([]*identitysharemodel.ContactPerson, 0, len(items))
	for _, item := range items {
		result = append(result, ContactPersonToModel(item))
	}
	return result
}

func PbContactPerson(c *identitysharemodel.ContactPerson) *etop.ContactPerson {
	if c == nil {
		return nil
	}
	return &etop.ContactPerson{
		Name:     c.Name,
		Position: c.Position,
		Email:    c.Email,
		Phone:    c.Phone,
	}
}

func Convert_core_ContactPerson_To_api_ContactPerson(in *identitytypes.ContactPerson) *etop.ContactPerson {
	if in == nil {
		return nil
	}
	return &etop.ContactPerson{
		Name:     in.Name,
		Position: in.Position,
		Email:    in.Email,
		Phone:    in.Phone,
	}
}

func PbContactPersons(items []*identitysharemodel.ContactPerson) []*etop.ContactPerson {
	if items == nil {
		return nil
	}
	result := make([]*etop.ContactPerson, 0, len(items))
	for _, item := range items {
		result = append(result, PbContactPerson(item))
	}
	return result
}

func Convert_core_ContactPersons_To_api_ContactPersons(items []*identitytypes.ContactPerson) []*etop.ContactPerson {
	result := make([]*etop.ContactPerson, len(items))
	for i, item := range items {
		result[i] = Convert_core_ContactPerson_To_api_ContactPerson(item)
	}
	return result
}

func PbCompanyInfo(info *identitysharemodel.CompanyInfo) *etop.CompanyInfo {
	if info == nil {
		return nil
	}
	return &etop.CompanyInfo{
		Name:                info.Name,
		TaxCode:             info.TaxCode,
		Address:             info.Address,
		LegalRepresentative: PbContactPerson(info.LegalRepresentative),
	}
}

func Convert_core_CompanyInfo_To_api_CompanyInfo(in *identitytypes.CompanyInfo) *etop.CompanyInfo {
	if in == nil {
		return nil
	}
	return &etop.CompanyInfo{
		Name:                in.Name,
		TaxCode:             in.TaxCode,
		Address:             in.Address,
		LegalRepresentative: Convert_core_ContactPerson_To_api_ContactPerson(in.LegalRepresentative),
	}
}

func PbAddressNote(item *addressmodel.AddressNote) *etop.AddressNote {
	if item == nil {
		return nil
	}
	return &etop.AddressNote{
		OpenTime:   item.OpenTime,
		LunchBreak: item.LunchBreak,
		Note:       item.Note,
		Other:      item.Other,
	}
}

func PbAddressNoteToModel(item *etop.AddressNote) *addressmodel.AddressNote {
	if item == nil {
		return nil
	}
	return &addressmodel.AddressNote{
		OpenTime:   item.OpenTime,
		LunchBreak: item.LunchBreak,
		Note:       item.Note,
		Other:      item.Other,
	}
}

func PbCreateAddressToModel(accountID dot.ID, p *etop.CreateAddressRequest) (*addressmodel.Address, error) {
	address := &etop.Address{
		FullName:     p.FullName,
		FirstName:    p.FirstName,
		LastName:     p.LastName,
		Phone:        p.Phone,
		Position:     p.Position,
		Email:        p.Email,
		Country:      p.Country,
		Province:     p.Province,
		District:     p.District,
		Ward:         p.Ward,
		Zip:          p.Zip,
		DistrictCode: p.DistrictCode,
		ProvinceCode: p.ProvinceCode,
		WardCode:     p.WardCode,
		Address1:     p.Address1,
		Address2:     p.Address2,
		Type:         p.Type,
		Notes:        p.Notes,
		Coordinates:  p.Coordinates,
	}
	res, err := AddressToModel(address)
	if err != nil {
		return nil, err
	}
	res.AccountID = accountID
	return res, nil
}

func PbUpdateAddressToModel(accountID dot.ID, p *etop.UpdateAddressRequest) (*addressmodel.Address, error) {
	address := &etop.Address{
		FullName:     p.FullName,
		FirstName:    p.FirstName,
		LastName:     p.LastName,
		Phone:        p.Phone,
		Position:     p.Position,
		Email:        p.Email,
		Country:      p.Country,
		Province:     p.Province,
		District:     p.District,
		Ward:         p.Ward,
		Zip:          p.Zip,
		DistrictCode: p.DistrictCode,
		ProvinceCode: p.ProvinceCode,
		WardCode:     p.WardCode,
		Address1:     p.Address1,
		Address2:     p.Address2,
		Type:         p.Type,
		Notes:        p.Notes,
		Coordinates:  p.Coordinates,
	}
	res, err := AddressToModel(address)
	if err != nil {
		return nil, err
	}
	res.ID = p.Id
	res.AccountID = accountID
	return res, nil
}

func PbCreditExtended(item *creditmodel.CreditExtended) *etop.Credit {
	if item == nil {
		return nil
	}

	return &etop.Credit{
		Id:        item.ID,
		Amount:    item.Amount,
		ShopId:    item.ShopID,
		Type:      item.Type,
		Shop:      PbShop(item.Shop),
		CreatedAt: cmapi.PbTime(item.CreatedAt),
		UpdatedAt: cmapi.PbTime(item.UpdatedAt),
		PaidAt:    cmapi.PbTime(item.PaidAt),
		Status:    item.Status,
	}
}

func PbCreditExtendeds(items []*creditmodel.CreditExtended) []*etop.Credit {
	result := make([]*etop.Credit, len(items))
	for i, item := range items {
		result[i] = PbCreditExtended(item)
	}
	return result
}

func Convert_core_CreditExtended_to_api_Credit(item *credit.CreditExtended) *etop.Credit {
	if item == nil {
		return nil
	}

	return &etop.Credit{
		Id:        item.ID,
		Amount:    item.Amount,
		ShopId:    item.ShopID,
		Type:      item.Type,
		Shop:      Convert_core_Shop_To_api_Shop(item.Shop),
		CreatedAt: cmapi.PbTime(item.CreatedAt),
		UpdatedAt: cmapi.PbTime(item.UpdatedAt),
		PaidAt:    cmapi.PbTime(item.PaidAt),
		Status:    item.Status,
	}
}

func Convert_core_CreditExtendeds_to_api_Credits(items []*credit.CreditExtended) []*etop.Credit {
	result := make([]*etop.Credit, len(items))
	for i, item := range items {
		result[i] = Convert_core_CreditExtended_to_api_Credit(item)
	}
	return result
}

func ShippingServiceSelectStrategyToModel(s []*etop.ShippingServiceSelectStrategyItem) []*identitymodel.ShippingServiceSelectStrategyItem {
	if s == nil {
		return nil
	}
	var result = make([]*identitymodel.ShippingServiceSelectStrategyItem, len(s))
	for i, item := range s {
		result[i] = &identitymodel.ShippingServiceSelectStrategyItem{
			Key:   item.Key,
			Value: item.Value,
		}
	}
	return result
}

func SurveyInfoToModel(m *etop.SurveyInfo) *identitymodel.SurveyInfo {
	return &identitymodel.SurveyInfo{
		Key:      m.Key,
		Question: m.Question,
		Answer:   m.Answer,
	}
}

func SurveyInfosToModel(items []*etop.SurveyInfo) []*identitymodel.SurveyInfo {
	result := make([]*identitymodel.SurveyInfo, 0, len(items))
	for _, item := range items {
		result = append(result, SurveyInfoToModel(item))
	}
	return result
}

func PbSurveyInfo(info *identitymodel.SurveyInfo) *etop.SurveyInfo {
	if info == nil {
		return nil
	}
	return &etop.SurveyInfo{
		Key:      info.Key,
		Question: info.Question,
		Answer:   info.Answer,
	}
}

func PbSurveyInfos(items []*identitymodel.SurveyInfo) []*etop.SurveyInfo {
	result := make([]*etop.SurveyInfo, len(items))
	for i, item := range items {
		result[i] = PbSurveyInfo(item)
	}
	return result
}

func Convert_core_SurveyInfo_To_api_SurveyInfo(in *identity.SurveyInfo) *etop.SurveyInfo {
	if in == nil {
		return nil
	}
	return &etop.SurveyInfo{
		Key:      in.Key,
		Question: in.Question,
		Answer:   in.Answer,
	}
}

func Convert_core_SurveyInfos_To_api_SurveyInfors(items []*identity.SurveyInfo) []*etop.SurveyInfo {
	result := make([]*etop.SurveyInfo, len(items))
	for i, item := range items {
		result[i] = Convert_core_SurveyInfo_To_api_SurveyInfo(item)
	}
	return result
}

func PbShippingServiceSelectStrategy(items []*identitymodel.ShippingServiceSelectStrategyItem) []*etop.ShippingServiceSelectStrategyItem {
	if items == nil {
		return nil
	}
	var result = make([]*etop.ShippingServiceSelectStrategyItem, len(items))
	for i, item := range items {
		result[i] = &etop.ShippingServiceSelectStrategyItem{
			Key:   item.Key,
			Value: item.Value,
		}
	}
	return result
}

func Convert_core_ShippingServiceSelectStrategy_To_api_ShippingServiceSelectStrategy(ins []*identity.ShippingServiceSelectStrategyItem) []*etop.ShippingServiceSelectStrategyItem {
	if ins == nil {
		return nil
	}
	var result = make([]*etop.ShippingServiceSelectStrategyItem, len(ins))
	for i, item := range ins {
		result[i] = &etop.ShippingServiceSelectStrategyItem{
			Key:   item.Key,
			Value: item.Value,
		}
	}
	return result
}

func PbDevice(m *notimodel.Device) *etop.Device {
	return &etop.Device{
		Id:                m.ID,
		AccountId:         m.AccountID,
		DeviceId:          m.DeviceID,
		DeviceName:        m.DeviceName,
		ExternalDeviceId:  m.ExternalDeviceID,
		ExternalServiceId: m.ExternalServiceID,
		CreatedAt:         cmapi.PbTime(m.CreatedAt),
		UpdatedAt:         cmapi.PbTime(m.UpdatedAt),
	}
}

func PbNotification(m *notimodel.Notification) *etop.Notification {
	return &etop.Notification{
		Id:               m.ID,
		AccountId:        m.AccountID,
		Title:            m.Title,
		Message:          m.Message,
		IsRead:           m.IsRead,
		Entity:           m.Entity.String(),
		EntityId:         m.EntityID,
		SendNotification: m.SendNotification,
		SyncStatus:       m.SyncStatus,
		SeenAt:           cmapi.PbTime(m.SeenAt),
		CreatedAt:        cmapi.PbTime(m.CreatedAt),
		UpdatedAt:        cmapi.PbTime(m.UpdatedAt),
	}
}

func PbNotifications(items []*notimodel.Notification) []*etop.Notification {
	result := make([]*etop.Notification, len(items))
	for i, item := range items {
		result[i] = PbNotification(item)
	}
	return result
}

func PbCoordinates(in *ordertypes.Coordinates) *etop.Coordinates {
	if in == nil {
		return nil
	}
	return &etop.Coordinates{
		Latitude:  in.Latitude,
		Longitude: in.Longitude,
	}
}

func PbCoordinatesToModel(in *etop.Coordinates) *ordertypes.Coordinates {
	if in == nil {
		return nil
	}
	return &ordertypes.Coordinates{
		Latitude:  in.Latitude,
		Longitude: in.Longitude,
	}
}

func PbInvitation(m *invitation.Invitation) *etop.Invitation {
	if m == nil {
		return nil
	}
	var roles []string
	for _, role := range m.Roles {
		roles = append(roles, string(role))
	}
	return &etop.Invitation{
		Id:            m.ID,
		ShopId:        m.AccountID,
		Email:         m.Email,
		Phone:         m.Phone,
		FullName:      m.FullName,
		ShortName:     m.ShortName,
		Roles:         roles,
		Token:         m.Token,
		Status:        m.Status,
		InvitedBy:     m.InvitedBy,
		AcceptedAt:    cmapi.PbTime(m.AcceptedAt),
		DeclinedAt:    cmapi.PbTime(m.RejectedAt),
		ExpiresAt:     cmapi.PbTime(m.ExpiresAt),
		CreatedAt:     cmapi.PbTime(m.CreatedAt),
		UpdatedAt:     cmapi.PbTime(m.UpdatedAt),
		InvitationURL: m.InvitationURL,
	}
}

func PbInvitations(ms []*invitation.Invitation) []*etop.Invitation {
	res := make([]*etop.Invitation, len(ms))
	for i, m := range ms {
		res[i] = PbInvitation(m)
	}
	return res
}

func PbAuthorization(m *authorization.Authorization) *etop.Authorization {
	if m == nil {
		return nil
	}
	var roles, actions []string
	for _, role := range m.Roles {
		roles = append(roles, string(role))
	}
	for _, action := range m.Actions {
		actions = append(actions, string(action))
	}
	return &etop.Authorization{
		UserId: m.UserID,
		// TODO: fix
		Name:    m.FullName,
		Email:   m.Email,
		Roles:   roles,
		Actions: actions,
	}
}

func PbAuthorizations(ms []*authorization.Authorization) []*etop.Authorization {
	res := make([]*etop.Authorization, len(ms))
	for i, m := range ms {
		res[i] = PbAuthorization(m)
	}
	return res
}

func PbRelationship(m *authorization.Relationship) *etop.Relationship {
	if m == nil {
		return nil
	}
	var roles, actions []string
	for _, role := range m.Roles {
		roles = append(roles, string(role))
	}
	for _, action := range m.Actions {
		actions = append(actions, string(action))
	}
	return &etop.Relationship{
		UserID:      m.UserID,
		AccountID:   m.AccountID,
		FullName:    m.FullName,
		ShortName:   m.ShortName,
		Position:    m.Position,
		Roles:       roles,
		Permissions: actions,
		Deleted:     m.Deleted,
	}
}

func PbRelationships(ms []*authorization.Relationship) []*etop.Relationship {
	res := make([]*etop.Relationship, len(ms))
	for i, m := range ms {
		res[i] = PbRelationship(m)
	}
	return res
}
