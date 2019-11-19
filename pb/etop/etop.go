package etop

import (
	"context"

	"etop.vn/api/main/identity"
	"etop.vn/api/main/invitation"
	"etop.vn/api/main/location"
	ordertypes "etop.vn/api/main/ordering/types"
	notimodel "etop.vn/backend/com/handler/notifier/model"
	servicelocation "etop.vn/backend/com/main/location"
	"etop.vn/backend/pb/common"
	"etop.vn/backend/pb/etop/etc/address_type"
	"etop.vn/backend/pb/etop/etc/status3"
	pbs3 "etop.vn/backend/pb/etop/etc/status3"
	"etop.vn/backend/pb/etop/etc/try_on"
	"etop.vn/backend/pb/etop/etc/user_source"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/integration/bank"
)

var locationBus = servicelocation.New().MessageBus()

func (m *CreateUserRequest) Censor() {
	if m.Password != "" {
		m.Password = "..."
	}
	if m.RegisterToken != "" {
		m.RegisterToken = "..."
	}
}

func (m *LoginRequest) Censor() {
	if m.Password != "" {
		m.Password = "..."
	}
}

func (m *ChangePasswordRequest) Censor() {
	if m.CurrentPassword != "" {
		m.CurrentPassword = "..."
	}
	if m.NewPassword != "" {
		m.NewPassword = "..."
	}
	if m.ConfirmPassword != "" {
		m.ConfirmPassword = "..."
	}
}

func (m *ChangePasswordUsingTokenRequest) Censor() {
	if m.ResetPasswordToken != "" {
		m.ResetPasswordToken = "..."
	}
	if m.NewPassword != "" {
		m.NewPassword = "..."
	}
	if m.ConfirmPassword != "" {
		m.ConfirmPassword = "..."
	}
}

func (e *CompanyInfo) ToModel() *model.CompanyInfo {
	if e == nil {
		return nil
	}
	return &model.CompanyInfo{
		Name:                e.Name,
		TaxCode:             e.TaxCode,
		Address:             e.Address,
		Website:             e.Website,
		LegalRepresentative: e.LegalRepresentative.ToModel(),
	}
}

func (b *BankAccount) ToModel() *model.BankAccount {
	if b == nil {
		return nil
	}
	return &model.BankAccount{
		Name:          b.Name,
		Province:      b.Province,
		Branch:        b.Branch,
		AccountNumber: b.AccountNumber,
		AccountName:   b.AccountName,
	}
}

func (b *BankAccount) ToCoreBankAccount() *identity.BankAccount {
	if b == nil {
		return nil
	}
	return &identity.BankAccount{
		Name:          b.Name,
		Province:      b.Province,
		Branch:        b.Branch,
		AccountNumber: b.AccountNumber,
		AccountName:   b.AccountName,
	}
}

func (a *Address) ToModel() (*model.Address, error) {
	if a == nil {
		return nil, nil
	}
	res := &model.Address{
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
		Type:         a.Type.ToModel(),
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
		res.Coordinates = &model.Coordinates{
			Latitude:  a.Coordinates.Latitude,
			Longitude: a.Coordinates.Longitude,
		}
	}
	return res, nil
}

func PbUser(m *model.User) *User {
	if m == nil {
		panic("Nil user")
	}
	return &User{
		Id:        m.ID,
		FullName:  m.FullName,
		ShortName: m.ShortName,
		Phone:     m.Phone,
		Email:     m.Email,
		CreatedAt: common.PbTime(m.CreatedAt),
		UpdatedAt: common.PbTime(m.UpdatedAt),

		EmailVerifiedAt: common.PbTime(m.EmailVerifiedAt),
		PhoneVerifiedAt: common.PbTime(m.PhoneVerifiedAt),

		EmailVerificationSentAt: common.PbTime(m.EmailVerificationSentAt),
		PhoneVerificationSentAt: common.PbTime(m.PhoneVerificationSentAt),
		Source:                  user_source.PbUserSource(m.Source),
	}
}

func PbAccountType(t model.AccountType) AccountType {
	return AccountType(AccountType_value[string(t)])
}

func PbLoginAccount(m *model.AccountUserExtended) *LoginAccount {
	account := m.Account
	return &LoginAccount{
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

func PbUserAccount(m *model.AccountUserExtended) *UserAccountInfo {
	account := m.Account
	accUser := m.AccountUser
	user := m.User

	fullName, shortName := m.GetUserName()
	return &UserAccountInfo{
		UserId:               user.ID,
		UserFullName:         fullName,
		UserShortName:        shortName,
		AccountId:            account.ID,
		AccountName:          account.Name,
		AccountType:          PbAccountType(account.Type),
		Position:             accUser.Position,
		Permission:           PbPermission(accUser),
		Status:               status3.Pb(accUser.Status),
		ResponseStatus:       status3.Pb(accUser.ResponseStatus),
		InvitationSentBy:     accUser.InvitationSentBy,
		InvitationSentAt:     common.PbTime(accUser.InvitationSentAt),
		InvitationAcceptedAt: common.PbTime(accUser.InvitationAcceptedAt),
		DisabledAt:           common.PbTime(accUser.DisabledAt),
	}
}

func PbUserAccounts(items []*model.AccountUserExtended) []*UserAccountInfo {
	result := make([]*UserAccountInfo, len(items))
	for i, item := range items {
		result[i] = PbUserAccount(item)
	}
	return result
}

func PbUserAccountIncomplete(accUser *model.AccountUser, account *model.Account) *UserAccountInfo {
	return &UserAccountInfo{
		UserId:               accUser.UserID,
		UserFullName:         "",
		UserShortName:        "",
		AccountId:            accUser.AccountID,
		AccountName:          account.Name,
		AccountType:          PbAccountType(account.Type),
		Position:             accUser.Position,
		Permission:           PbPermission(accUser),
		Status:               status3.Pb(accUser.Status),
		InvitationSentBy:     accUser.InvitationSentBy,
		InvitationSentAt:     common.PbTime(accUser.InvitationSentAt),
		InvitationAcceptedAt: common.PbTime(accUser.InvitationAcceptedAt),
		DisabledAt:           common.PbTime(accUser.DisabledAt),
	}
}

func PbPermission(m *model.AccountUser) *Permission {
	return &Permission{
		Roles:       m.Roles,
		Permissions: m.Permissions,
	}
}

func PbPartner(m *model.Partner) *Partner {
	return &Partner{
		Id:             m.ID,
		Name:           m.Name,
		PublicName:     m.PublicName,
		Status:         status3.Pb(m.Status),
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

func PbPublicPartners(items []*model.Partner) []*PublicAccountInfo {
	res := make([]*PublicAccountInfo, len(items))
	for i, item := range items {
		res[i] = PbPublicAccountInfo(item)
	}
	return res
}

func PbPublicAccountInfo(m model.AccountInterface) *PublicAccountInfo {
	switch m := m.(type) {
	case *model.Partner:
		return &PublicAccountInfo{
			Id:       m.ID,
			Name:     m.PublicName, // public name here!
			Type:     PbAccountType(model.TypePartner),
			ImageUrl: m.ImageURL,
			Website:  m.WebsiteURL,
		}
	default:
		account := m.GetAccount()
		return &PublicAccountInfo{
			Id:       account.ID,
			Name:     account.Name,
			Type:     PbAccountType(account.Type),
			ImageUrl: account.ImageURL,
			Website:  "",
		}
	}
}

func PbShop(m *model.Shop) *Shop {
	return &Shop{
		Id:          m.ID,
		Name:        m.Name,
		Status:      status3.Pb(m.Status),
		Phone:       m.Phone,
		BankAccount: PbBankAccount(m.BankAccount),
		WebsiteUrl:  m.WebsiteURL,
		ImageUrl:    m.ImageURL,
		Email:       m.Email,
		OwnerId:     m.OwnerID,
		TryOn:       try_on.PbTryOn(m.TryOn),
	}
}

func PbShopExtended(m *model.ShopExtended) *Shop {
	return &Shop{
		Id:                            m.ID,
		InventoryOverstock:            cm.BoolDefault(m.InventoryOverstock, true),
		Name:                          m.Name,
		Status:                        status3.Pb(m.Status),
		Address:                       PbAddress(m.Address),
		Phone:                         m.Phone,
		BankAccount:                   PbBankAccount(m.BankAccount),
		WebsiteUrl:                    m.WebsiteURL,
		ImageUrl:                      m.ImageURL,
		Email:                         m.Email,
		ShipToAddressId:               m.ShipToAddressID,
		ShipFromAddressId:             m.ShipFromAddressID,
		AutoCreateFfm:                 m.AutoCreateFFM,
		TryOn:                         try_on.PbTryOn(m.TryOn),
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

func PbShopExtendeds(items []*model.ShopExtended) []*Shop {
	result := make([]*Shop, len(items))
	for i, item := range items {
		result[i] = PbShopExtended(item)
	}
	return result
}

func PbProvinces(items []*location.Province) []*Province {
	res := make([]*Province, len(items))
	for i, item := range items {
		res[i] = PbProvince(item)
	}
	return res
}

func PbProvince(m *location.Province) *Province {
	return &Province{
		Code:       m.Code,
		Name:       m.Name,
		Region:     m.Region.Name(),
		RegionCode: int64(m.Region),
	}
}

func PbDistricts(items []*location.District) []*District {
	res := make([]*District, len(items))
	for i, item := range items {
		res[i] = PbDistrict(item)
	}
	return res
}

func PbDistrict(item *location.District) *District {
	return &District{
		Code:         item.Code,
		ProvinceCode: item.ProvinceCode,
		Name:         item.Name,
	}
}

func PbWards(items []*location.Ward) []*Ward {
	res := make([]*Ward, len(items))
	for i, item := range items {
		res[i] = PbWard(item)
	}
	return res
}

func PbWard(item *location.Ward) *Ward {
	return &Ward{
		Code:         item.Code,
		DistrictCode: item.DistrictCode,
		Name:         item.Name,
	}
}

func PbBanks(items []*bank.Bank) []*Bank {
	res := make([]*Bank, len(items))
	for i, item := range items {
		res[i] = PbBank(item)
	}
	return res
}

func PbBank(item *bank.Bank) *Bank {
	return &Bank{
		Code: item.MaNganHang,
		Name: item.TenNH,
		Type: item.Loai,
	}
}

func PbBankProvinces(items []*bank.Province) []*BankProvince {
	res := make([]*BankProvince, len(items))
	for i, item := range items {
		res[i] = PbBankProvince(item)
	}
	return res
}

func PbBankProvince(item *bank.Province) *BankProvince {
	return &BankProvince{
		Code:     item.MaTinh,
		Name:     item.TenTinhThanh,
		BankCode: item.MaNganHang,
	}
}

func PbBankBranches(items []*bank.Branch) []*BankBranch {
	res := make([]*BankBranch, len(items))
	for i, item := range items {
		res[i] = PbBankBranch(item)
	}
	return res
}

func PbBankBranch(item *bank.Branch) *BankBranch {
	return &BankBranch{
		Code:         item.MaChiNhanh,
		Name:         item.TenChiNhanh,
		BankCode:     item.MaNganHang,
		ProvinceCode: item.MaTinh,
	}
}

func PbAddresses(items []*model.Address) []*Address {
	result := make([]*Address, len(items))
	for i, item := range items {
		result[i] = PbAddress(item)
	}
	return result
}

func PbAddress(a *model.Address) *Address {
	if a == nil {
		return nil
	}
	res := &Address{
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
		Type:         address_type.PbType(a.Type),
		Notes:        PbAddressNote(a.Notes),
	}
	if a.Coordinates != nil {
		res.Coordinates = &Coordinates{
			Latitude:  a.Coordinates.Latitude,
			Longitude: a.Coordinates.Longitude,
		}
	}
	return res
}

func PbBankAccount(b *model.BankAccount) *BankAccount {
	if b == nil {
		return nil
	}
	return &BankAccount{
		Name:          b.Name,
		Province:      b.Province,
		Branch:        b.Branch,
		AccountName:   b.AccountName,
		AccountNumber: b.AccountNumber,
	}
}

func Convert_core_BankAccount_To_api_BankAccount(in *identity.BankAccount) *BankAccount {
	if in == nil {
		return nil
	}
	return &BankAccount{
		Name:          in.Name,
		Province:      in.Province,
		Branch:        in.Branch,
		AccountNumber: in.AccountNumber,
		AccountName:   in.AccountName,
	}
}

func (m *ContactPerson) ToModel() *model.ContactPerson {
	return &model.ContactPerson{
		Name:     m.Name,
		Position: m.Position,
		Phone:    m.Phone,
		Email:    m.Email,
	}
}

func ContactPersonsToModel(items []*ContactPerson) []*model.ContactPerson {
	result := make([]*model.ContactPerson, 0, len(items))
	for _, item := range items {
		result = append(result, item.ToModel())
	}
	return result
}

func PbContactPerson(c *model.ContactPerson) *ContactPerson {
	if c == nil {
		return nil
	}
	return &ContactPerson{
		Name:     c.Name,
		Position: c.Position,
		Email:    c.Email,
		Phone:    c.Phone,
	}
}

func PbContactPersons(items []*model.ContactPerson) []*ContactPerson {
	if items == nil {
		return nil
	}
	result := make([]*ContactPerson, 0, len(items))
	for _, item := range items {
		result = append(result, PbContactPerson(item))
	}
	return result
}

func PbCompanyInfo(info *model.CompanyInfo) *CompanyInfo {
	if info == nil {
		return nil
	}
	return &CompanyInfo{
		Name:                info.Name,
		TaxCode:             info.TaxCode,
		Address:             info.Address,
		LegalRepresentative: PbContactPerson(info.LegalRepresentative),
	}
}

func PbAddressNote(item *model.AddressNote) *AddressNote {
	if item == nil {
		return nil
	}
	return &AddressNote{
		OpenTime:   item.OpenTime,
		LunchBreak: item.LunchBreak,
		Note:       item.Note,
		Other:      item.Other,
	}
}

func PbAddressNoteToModel(item *AddressNote) *model.AddressNote {
	if item == nil {
		return nil
	}
	return &model.AddressNote{
		OpenTime:   item.OpenTime,
		LunchBreak: item.LunchBreak,
		Note:       item.Note,
		Other:      item.Other,
	}
}

func PbCreateAddressToModel(accountID int64, p *CreateAddressRequest) (*model.Address, error) {
	address := &Address{
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
	res, err := address.ToModel()
	if err != nil {
		return nil, err
	}
	res.AccountID = accountID
	return res, nil
}

func PbUpdateAddressToModel(accountID int64, p *UpdateAddressRequest) (*model.Address, error) {
	address := &Address{
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
	res, err := address.ToModel()
	if err != nil {
		return nil, err
	}
	res.ID = p.Id
	res.AccountID = accountID
	return res, nil
}

func PbCreditExtended(item *model.CreditExtended) *Credit {
	if item == nil {
		return nil
	}

	return &Credit{
		Id:        item.ID,
		Amount:    int64(item.Amount),
		ShopId:    item.ShopID,
		Type:      item.Type,
		Shop:      PbShop(item.Shop),
		CreatedAt: common.PbTime(item.CreatedAt),
		UpdatedAt: common.PbTime(item.UpdatedAt),
		PaidAt:    common.PbTime(item.PaidAt),
		Status:    status3.Pb(item.Status),
	}
}

func PbCreditExtendeds(items []*model.CreditExtended) []*Credit {
	result := make([]*Credit, len(items))
	for i, item := range items {
		result[i] = PbCreditExtended(item)
	}
	return result
}

func ShippingServiceSelectStrategyToModel(s []*ShippingServiceSelectStrategyItem) []*model.ShippingServiceSelectStrategyItem {
	if s == nil {
		return nil
	}
	var result = make([]*model.ShippingServiceSelectStrategyItem, len(s))
	for i, item := range s {
		result[i] = &model.ShippingServiceSelectStrategyItem{
			Key:   item.Key,
			Value: item.Value,
		}
	}
	return result
}

func (m *SurveyInfo) ToModel() *model.SurveyInfo {
	return &model.SurveyInfo{
		Key:      m.Key,
		Question: m.Question,
		Answer:   m.Answer,
	}
}

func SurveyInfosToModel(items []*SurveyInfo) []*model.SurveyInfo {
	result := make([]*model.SurveyInfo, 0, len(items))
	for _, item := range items {
		result = append(result, item.ToModel())
	}
	return result
}

func PbSurveyInfo(info *model.SurveyInfo) *SurveyInfo {
	if info == nil {
		return nil
	}
	return &SurveyInfo{
		Key:      info.Key,
		Question: info.Question,
		Answer:   info.Answer,
	}
}

func PbSurveyInfos(items []*model.SurveyInfo) []*SurveyInfo {
	result := make([]*SurveyInfo, len(items))
	for i, item := range items {
		result[i] = PbSurveyInfo(item)
	}
	return result
}

func PbShippingServiceSelectStrategy(items []*model.ShippingServiceSelectStrategyItem) []*ShippingServiceSelectStrategyItem {
	if items == nil {
		return nil
	}
	var result = make([]*ShippingServiceSelectStrategyItem, len(items))
	for i, item := range items {
		result[i] = &ShippingServiceSelectStrategyItem{
			Key:   item.Key,
			Value: item.Value,
		}
	}
	return result
}

func PbDevice(m *notimodel.Device) *Device {
	return &Device{
		Id:                m.ID,
		AccountId:         m.AccountID,
		DeviceId:          m.DeviceID,
		DeviceName:        m.DeviceName,
		ExternalDeviceId:  m.ExternalDeviceID,
		ExternalServiceId: int32(m.ExternalServiceID),
		CreatedAt:         common.PbTime(m.CreatedAt),
		UpdatedAt:         common.PbTime(m.UpdatedAt),
	}
}

func PbNotification(m *notimodel.Notification) *Notification {
	return &Notification{
		Id:               m.ID,
		AccountId:        m.AccountID,
		Title:            m.Title,
		Message:          m.Message,
		IsRead:           m.IsRead,
		Entity:           string(m.Entity),
		EntityId:         m.EntityID,
		SendNotification: m.SendNotification,
		SyncStatus:       pbs3.Pb(m.SyncStatus),
		SeenAt:           common.PbTime(m.SeenAt),
		CreatedAt:        common.PbTime(m.CreatedAt),
		UpdatedAt:        common.PbTime(m.UpdatedAt),
	}
}

func PbNotifications(items []*notimodel.Notification) []*Notification {
	result := make([]*Notification, len(items))
	for i, item := range items {
		result[i] = PbNotification(item)
	}
	return result
}

func PbCoordinates(in *ordertypes.Coordinates) *Coordinates {
	if in == nil {
		return nil
	}
	return &Coordinates{
		Latitude:  in.Latitude,
		Longitude: in.Longitude,
	}
}

func PbCoordinatesToModel(in *Coordinates) *ordertypes.Coordinates {
	if in == nil {
		return nil
	}
	return &ordertypes.Coordinates{
		Latitude:  in.Latitude,
		Longitude: in.Longitude,
	}
}

func PbInvitation(m *invitation.Invitation) *Invitation {
	if m == nil {
		return nil
	}
	var roles []string
	for _, role := range m.Roles {
		roles = append(roles, string(role))
	}
	return &Invitation{
		Id:         m.ID,
		ShopId:     m.AccountID,
		Email:      m.Email,
		Roles:      roles,
		Token:      m.Token,
		Status:     pbs3.Pb(model.Status3(m.Status)),
		InvitedBy:  m.InvitedBy,
		AcceptedAt: common.PbTime(m.AcceptedAt),
		DeclinedAt: common.PbTime(m.RejectedAt),
		ExpiredAt:  common.PbTime(m.ExpiresAt),
		CreatedAt:  common.PbTime(m.CreatedAt),
		UpdatedAt:  common.PbTime(m.UpdatedAt),
	}
}

func PbInvitations(ms []*invitation.Invitation) []*Invitation {
	res := make([]*Invitation, len(ms))
	for i, m := range ms {
		res[i] = PbInvitation(m)
	}
	return res
}
