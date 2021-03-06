package identity

import (
	"fmt"
	"strings"
	"time"

	"o.o/api/main/address"
	identitytypes "o.o/api/main/identity/types"
	"o.o/api/meta"
	"o.o/api/top/types/etc/account_type"
	"o.o/api/top/types/etc/ghn_note_code"
	"o.o/api/top/types/etc/shop_user_role"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/try_on"
	"o.o/api/top/types/etc/user_otp_action"
	"o.o/api/top/types/etc/user_role"
	"o.o/api/top/types/etc/user_source"
	"o.o/capi/dot"
	cmutil "o.o/capi/util"
	"o.o/common/xerrors"
)

// +gen:event:topic=event/identity

type (
	SubjectType string
)

const (
	SubjectTypeAccount SubjectType = "account"
	SubjectTypeUser    SubjectType = "user"
)

const LatestUserOTPKey = "LatestUserOTP"

func GetLatestUserOTPRedisKey(userID dot.ID, phone string) string {
	if userID != 0 {
		return fmt.Sprintf("%v:uID:%v", LatestUserOTPKey, userID)
	}
	return fmt.Sprintf("%v:uPhone:%v", LatestUserOTPKey, phone)
}

type Permission struct {
	Roles       []string
	Permissions []string
}

func (p Permission) GetUserRoles() []user_role.UserRole {
	var res = make([]user_role.UserRole, 0, len(p.Roles))
	for _, role := range p.Roles {
		r, ok := user_role.ParseUserRole(role)
		if !ok {
			continue
		}
		res = append(res, r)
	}
	return res
}

func (p Permission) GetShopUserRoles() []shop_user_role.UserRole {
	var res = make([]shop_user_role.UserRole, 0, len(p.Roles))
	for _, role := range p.Roles {
		r, ok := shop_user_role.ParseUserRole(role)
		if !ok {
			continue
		}
		res = append(res, r)
	}
	return res
}

type Shop struct {
	ID      dot.ID
	Name    string
	OwnerID dot.ID
	IsTest  int

	AddressID          dot.ID
	ShipToAddressID    dot.ID
	ShipFromAddressID  dot.ID
	Phone              string
	WebsiteURL         string
	ImageURL           string
	Email              string
	Code               string
	AutoCreateFFM      bool
	InventoryOverstock dot.NullBool
	GhnNoteCode        ghn_note_code.GHNNoteCode
	RecognizedHosts    []string
	// MoneyTransactionRRule format:
	// FREQ=DAILY
	// FREQ=WEEKLY;BYDAY=TU,TH,SA
	// FREQ=MONTHLY;BYDAY=-1FR
	MoneyTransactionRRule         string `sq:"'money_transaction_rrule'"`
	SurveyInfo                    []*SurveyInfo
	ShippingServiceSelectStrategy []*ShippingServiceSelectStrategyItem

	Status                  status3.Status
	CreatedAt               time.Time
	UpdatedAt               time.Time
	DeletedAt               time.Time
	BankAccount             *identitytypes.BankAccount
	TryOn                   try_on.TryOnCode
	CompanyInfo             *identitytypes.CompanyInfo
	WLPartnerID             dot.ID
	IsPriorMoneyTransaction dot.NullBool
}

func (s *Shop) CheckInfo() error {
	if s.MoneyTransactionRRule != "" {
		acceptStrings := []string{"FREQ", "BYMONTHDAY", "BYDAY", "BYSETPOS"}
		ss := strings.Split(s.MoneyTransactionRRule, ";")
		for _, s := range ss {
			prefix := strings.Split(s, "=")[0]
			if !cmutil.ListStringsContain(acceptStrings, prefix) {
				return xerrors.Errorf(xerrors.InvalidArgument, nil, "?????nh d???ng MoneyTransactionRRule kh??ng h???p l???. ?????nh d???ng ????ng: FREQ=WEEKLY;BYDAY=TU,TH,SA. Vui l??ng truy c???p https://icalendar.org/rrule-tool.html ????? bi???t th??m chi ti???t.")
			}
		}
	}
	if s.ShippingServiceSelectStrategy != nil {
		for _, item := range s.ShippingServiceSelectStrategy {
			key := strings.ToLower(item.Key)
			if key == "provider" || key == "shipping_provider" {
				return xerrors.Errorf(xerrors.InvalidArgument, nil, "Vui l??ng s??? d???ng `carrier` thay v?? %v", item.Key)
			}
		}
	}
	return nil
}

type ShippingServiceSelectStrategyItem struct {
	Key   string
	Value string
}

type SurveyInfo struct {
	Key      string `json:"key"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type Affiliate struct {
	ID          dot.ID
	OwnerID     dot.ID
	Name        string
	Phone       string
	Email       string
	IsTest      int
	Status      status3.Status
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
	BankAccount *identitytypes.BankAccount
}

type User struct {
	ID dot.ID

	FullName  string
	ShortName string
	Email     string
	Phone     string

	Status status3.Status // 1: actual user, 0: stub, -1: disabled

	EmailVerifiedAt         time.Time
	PhoneVerifiedAt         time.Time
	EmailVerificationSentAt time.Time
	PhoneVerificationSentAt time.Time

	IsTest      int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	RefUserID   dot.ID
	RefSaleID   dot.ID
	WLPartnerID dot.ID
	Source      user_source.UserSource

	BlockedAt   time.Time
	BlockedBy   dot.ID
	BlockReason string
	IsBlocked   bool

	AgreedTOSAt       time.Time
	AgreedEmailInfoAt time.Time
}

type UserFtRefSaff struct {
	*User
	RefAff  string
	RefSale string
}

type AccountCreatedEvent struct {
	meta.EventMeta
	ShopID dot.ID
	UserID dot.ID
}

type UserCreatedEvent struct {
	meta.EventMeta
	UserID dot.ID

	Invitation *UserInvitation
}

type AccountDeletingEvent struct {
	AccountID   dot.ID
	AccountType account_type.AccountType
}

type UserUpdatedEvent struct {
	AccountID dot.ID
	UserID    dot.ID
	FullName  string
	Phone     string
}

type UserInvitation struct {
	Token      string
	AutoAccept bool
}

type GetCustomersByShop struct {
	meta.EventMeta
	ShopID dot.ID
}

type Partner struct {
	ID dot.ID

	OwnerID dot.ID

	Name string

	PublicName string

	ImageURL string

	WebsiteURL string

	WhiteLabelKey string

	WhitelistIPs []string
}

type ShopExtended struct {
	*Shop
	Address *address.Address
	User    *User
}

type AccountUser struct {
	AccountID    dot.ID
	UserID       dot.ID
	DepartmentID dot.ID

	Status         status3.Status // 1: activated, -1: rejected/disabled, 0: pending
	ResponseStatus status3.Status // 1: accepted,  -1: rejected, 0: pending

	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
	DeletedAt time.Time

	Permission `sq:"inline"`

	FullName            string
	FullNameNorm        string
	Phone               string
	PhoneNorm           string
	ExtensionNumberNorm string
	ShortName           string
	Position            string

	InvitationSentAt     time.Time
	InvitationSentBy     dot.ID
	InvitationAcceptedAt time.Time
	InvitationRejectedAt time.Time

	DisabledAt    time.Time
	DisabledBy    time.Time
	DisableReason string

	Rid dot.ID
}

type UserRefSaff struct {
	UserID  dot.ID
	RefAff  dot.NullString
	RefSale dot.NullString
}

type UserInternal struct {
	ID      dot.ID
	Hashpwd string

	UpdatedAt time.Time `sq:"update"`
}

type PartnerRelation struct {
	AuthKey           string
	PartnerID         dot.ID
	SubjectID         dot.ID
	SubjectType       SubjectType
	ExternalSubjectID string

	Nonce     dot.ID
	Status    status3.Status
	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
	DeletedAt time.Time
}

type AccountUserExtended struct {
	UserID       dot.ID
	AccountID    dot.ID
	DepartmentID dot.ID
	Roles        []string
	Permissions  []string
	FullName     string
	ShortName    string
	Email        string
	Phone        string
	Position     string
	Deleted      bool
}

type AccountUserWithGroupByDepartment struct {
	DepartmentID dot.ID
	Count        int64
}

type LatestUserOTPData struct {
	OTP    string                        `json:"otp"`
	Action user_otp_action.UserOTPAction `json:"action"`
}
