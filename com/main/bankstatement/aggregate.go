package bankstatement

import (
	"context"
	"regexp"

	"o.o/api/main/bankstatement"
	"o.o/api/main/identity"
	com "o.o/backend/com/main"
	"o.o/backend/com/main/bankstatement/convert"
	"o.o/backend/com/main/bankstatement/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/validate"
	"o.o/capi"
	"o.o/common/l"
)

var ll = l.New()
var _ bankstatement.Aggregate = &BankStatementAggregate{}
var scheme = conversion.Build(convert.RegisterConversions)

type BankStatementAggregate struct {
	dbTx          cmsql.Transactioner
	store         sqlstore.BankStatementFactory
	eventBus      capi.EventBus
	identityQuery identity.QueryBus
}

func NewAggregateBankStatement(
	bus capi.EventBus,
	db com.MainDB,
	identityQ identity.QueryBus,
) *BankStatementAggregate {
	return &BankStatementAggregate{
		dbTx:          (*cmsql.Database)(db),
		identityQuery: identityQ,
		eventBus:      bus,
		store:         sqlstore.NewBankStatementStore(db),
	}
}

// Bank Statement description format: [shop_code] [phone]
// This regex devide to 4 part:
// [1]: any whitespace or tab in the beginning
// [2]: 4 characters include uppercase letter and number - shop code
// [3]: one or more whitescpace/tab in the middle
// [4]: Vietnamese phone number - user phone
var bankStatementDescriptionRegex = regexp.MustCompile(`^(\s*)([A-Z0-9]{4})(\s+)((((\+84|84|0){1})(3|5|7|8|9))+([0-9]{8}))`)

func BankStatementAggregateMessageBus(q *BankStatementAggregate) bankstatement.CommandBus {
	b := bus.New()
	return bankstatement.NewAggregateHandler(q).RegisterHandlers(b)
}

func (a BankStatementAggregate) CreateBankStatement(ctx context.Context, args *bankstatement.CreateBankStatementArgs) (*bankstatement.BankStatement, error) {
	if err := args.Validate(); err != nil {
		return nil, err
	}
	shopInfo, err := ParseBankStatementDescription(args.Description)
	if err != nil {
		return nil, err
	}

	_, err = a.store(ctx).ExternalTrxnID(args.ExternalTransactionID).Get()
	if err != nil && cm.ErrorCode(err) != cm.NotFound {
		return nil, err
	}
	if err == nil {
		// bank statement existed
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Bank statement was existed")
	}

	shopQuery := &identity.GetShopByCodeQuery{
		Code: shopInfo.ShopCode,
	}
	if err := a.identityQuery.Dispatch(ctx, shopQuery); err != nil {
		return nil, err
	}
	shop := shopQuery.Result

	userQuery := &identity.GetUserByPhoneQuery{
		Phone: shopInfo.UserPhone,
	}
	if err := a.identityQuery.Dispatch(ctx, userQuery); err != nil {
		return nil, err
	}
	if shop.OwnerID != userQuery.Result.ID {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Shop code does not belong to the user phone")
	}

	bankState := convert.Apply_bankstatement_CreateBankStatementArgs_bankstatement_BankStatement(args, nil)

	bankState.ID = cm.NewID()
	bankState.OtherInfo = args.OtherInfo
	bankState.AccountID = shop.ID

	if err := a.dbTx.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		if err := a.store(ctx).Create(bankState); err != nil {
			return err
		}

		event := &bankstatement.BankStatementCreatedEvent{
			ID: bankState.ID,
		}
		if err := a.eventBus.Publish(ctx, event); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return bankState, nil
}

type ShopInfo struct {
	ShopCode  string
	UserPhone string
}

func ParseBankStatementDescription(desc string) (*ShopInfo, error) {
	if desc == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing description")
	}
	res := bankStatementDescriptionRegex.FindStringSubmatch(desc)
	if len(res) < 5 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Description does not match format")
	}
	shopCode := res[2]
	phoneNumber := res[4]
	phone, ok := validate.NormalizePhone(phoneNumber)
	if !ok {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Phone number does not valid")
	}
	return &ShopInfo{
		ShopCode:  shopCode,
		UserPhone: phone.String(),
	}, nil
}
