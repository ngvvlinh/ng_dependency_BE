package accountuser

import (
	"bytes"
	"context"
	"io/ioutil"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"

	"o.o/api/main/authorization"
	"o.o/api/main/identity"
	"o.o/api/top/types/etc/shop_user_role"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/user_source"
	com "o.o/backend/com/main"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/httpx"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/validate"
	"o.o/capi/dot"
)

type AccountUserService struct {
	dbMain        *cmsql.Database
	identityAggr  identity.CommandBus
	identityQuery identity.QueryBus
}

var mapRole = map[string]authorization.Role{
	"inventory_management":    authorization.RoleInventoryManagement,
	"salesman":                authorization.RoleSalesMan,
	"analyst":                 authorization.RoleAnalyst,
	"accountant":              authorization.RoleAccountant,
	"purchasing_management":   authorization.RolePurchasingManagement,
	"staff_management":        authorization.RoleStaffManagement,
	"telecom_customerservice": authorization.RoleTelecomCustomerService,
}

type Line struct {
	ID         string
	Name       string
	Phone      string
	Email      string
	Roles      []string
	ShopCode   string
	OwnerPhone string
	AccountID  dot.ID
	UserID     dot.ID
}

func New(
	dbMain com.MainDB,
	identityAggr identity.CommandBus,
	identityQuery identity.QueryBus,
) *AccountUserService {
	return &AccountUserService{
		dbMain:        dbMain,
		identityAggr:  identityAggr,
		identityQuery: identityQuery,
	}
}

func (s *AccountUserService) HandleImportAccountUser(c *httpx.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		return cm.Errorf(cm.InvalidArgument, nil, "Invalid request")
	}

	files := form.File["files"]
	switch len(files) {
	case 0:
		return cm.Errorf(cm.InvalidArgument, nil, "No file")
	case 1:
	// continue
	default:
		return cm.Errorf(cm.InvalidArgument, nil, "Too many files")
	}

	file, err := files[0].Open()
	if err != nil {
		return cm.Errorf(cm.InvalidArgument, err, "Can not read file")
	}
	defer file.Close()

	rawData, err := ioutil.ReadAll(file)
	if err != nil {
		return cm.Errorf(cm.InvalidArgument, err, "Không thể đọc được file.")
	}
	excelFile, err := excelize.OpenReader(bytes.NewReader(rawData))
	if err != nil {
		return cm.Errorf(cm.InvalidArgument, err, "Không thể đọc được file")
	}

	sheetName := excelFile.GetSheetName(0)
	rows, err := excelFile.GetRows(sheetName)
	if err != nil {
		return cm.Errorf(cm.InvalidArgument, err, "Không thể đọc được file")
	}
	if len(rows) <= 1 {
		return cm.Errorf(cm.InvalidArgument, nil, "File không có nội dung")
	}

	ctx := bus.Ctx()
	success := 0
	total := len(rows) - 1
	for i, row := range rows {
		if i == 0 {
			continue
		}

		line, _err := s.parseRow(ctx, row)
		if _err != nil {
			return cm.Errorf(cm.InvalidArgument, _err, "Can not parse row").WithMetap("row", row)
		}

		accountUser, _err := s.getAccountUser(ctx, line.OwnerPhone, line.ShopCode)
		if _err != nil {
			return cm.Errorf(cm.InvalidArgument, _err, "Can not get account").WithMetap("row", row)
		}
		line.AccountID = accountUser.AccountID

		var user *identity.User
		user, _err = s.registerUser(ctx, line)
		if _err != nil {
			return cm.Errorf(cm.InvalidArgument, _err, "Can not register user").WithMetap("row", row)
		}

		line.UserID = user.ID
		err = s.createAccountUser(ctx, line)
		if _err != nil {
			return cm.Errorf(cm.InvalidArgument, err, "Can not create account user").WithMetap("row", row)
		}

		success += 1
	}

	c.SetResult(map[string]interface{}{
		"code":    "ok",
		"total":   total,
		"success": success,
	})
	return nil
}

func (s *AccountUserService) parseRow(ctx context.Context, row []string) (*Line, error) {
	phone, ok := validate.NormalizePhone(row[2])
	if !ok {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Phone does not valid").WithMetap("row", row)
	}

	var email validate.NormalizedEmail
	if row[3] != "" {
		email, ok = validate.NormalizeEmail(row[3])
		if !ok {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Email does not valid").WithMetap("row", row)
		}
	}

	ownerPhone, ok := validate.NormalizePhone(row[5])
	if !ok {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Owner phone does not valid").WithMetap("row", row)
	}

	var shopCode string
	if row[6] != "" {
		shopCode = strings.TrimSpace(row[6])
	}

	strRoles := strings.Split(row[4], ",")
	var roles = []string{}
	for _, role := range strRoles {
		role = strings.ReplaceAll(role, " ", "")
		if _, ok = mapRole[role]; !ok {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Roles does not valid").WithMetap("row", row)
		}
		roles = append(roles, role)
	}

	return &Line{
		ID:         row[0],
		Name:       row[1],
		Phone:      phone.String(),
		Email:      email.String(),
		Roles:      roles,
		ShopCode:   shopCode,
		OwnerPhone: ownerPhone.String(),
	}, nil
}

func (s *AccountUserService) getAccountUser(ctx context.Context, phone, shopCode string) (*identity.AccountUser, error) {
	getUserQuery := &identity.GetUserByPhoneOrEmailQuery{
		Phone: phone,
	}
	if err := s.identityQuery.Dispatch(ctx, getUserQuery); err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Can not get owner by phone")
	}
	owner := getUserQuery.Result

	var accountUser *identity.AccountUser
	getAccountUsersQuery := &identity.ListAccountUsersQuery{
		Roles: []shop_user_role.UserRole{
			shop_user_role.Owner,
		},
		UserIDs: []dot.ID{owner.ID},
	}

	err := s.identityQuery.Dispatch(ctx, getAccountUsersQuery)
	switch cm.ErrorCode(err) {
	case cm.NoError:
		accountUsers := getAccountUsersQuery.Result.AccountUsers
		switch len(accountUsers) {
		case 0:
			return nil, cm.Errorf(cm.InvalidArgument, nil, "User doen not have any shop")
		case 1:
			accountUser = accountUsers[0]
		default:
			if shopCode == "" {
				return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing shop code")
			}
			getShopQuery := &identity.GetShopByCodeQuery{
				Code: shopCode,
			}
			if err = s.identityQuery.Dispatch(ctx, getShopQuery); err != nil {
				return nil, cm.Errorf(cm.InvalidArgument, nil, "Can not get shop by shop code")
			}
			shop := getShopQuery.Result
			for _, au := range accountUsers {
				if au.AccountID == shop.ID {
					accountUser = au
					break
				}
			}
			if accountUser == nil {
				return nil, cm.Errorf(cm.InvalidArgument, nil, "Shop code does not valid")
			}
		}
	default:
		return nil, err
	}

	return accountUser, nil

}

func (s *AccountUserService) registerUser(ctx context.Context, args *Line) (user *identity.User, err error) {
	getUserQuery := &identity.GetUserByPhoneOrEmailQuery{
		Phone: args.Phone,
	}
	err = s.identityQuery.Dispatch(ctx, getUserQuery)
	switch cm.ErrorCode(err) {
	case cm.NoError:
		user = getUserQuery.Result
	case cm.NotFound:
		cmd := &identity.CreateUserCommand{
			FullName:                args.Name,
			Phone:                   args.Phone,
			Password:                args.Phone,
			Email:                   args.Email,
			Status:                  status3.P,
			Source:                  user_source.Etop,
			PhoneVerifiedAt:         time.Now(),
			PhoneVerificationSentAt: time.Now(),
			AgreeTOS:                true,
		}
		if args.Email != "" {
			cmd.EmailVerificationSentAt = time.Now()
			cmd.EmailVerifiedAt = time.Now()
		}
		if err = s.identityAggr.Dispatch(ctx, cmd); err != nil {
			return nil, err
		}
		user = cmd.Result
	default:
		return nil, err
	}
	return
}

func (s *AccountUserService) createAccountUser(ctx context.Context, args *Line) error {
	getAccountUserQuery := &identity.GetAccountUserQuery{
		UserID:    args.UserID,
		AccountID: args.AccountID,
	}
	err := s.identityQuery.Dispatch(ctx, getAccountUserQuery)
	switch cm.ErrorCode(err) {
	case cm.NoError:
	case cm.NotFound:
		cmd := &identity.CreateAccountUserCommand{
			AccountID: args.AccountID,
			UserID:    args.UserID,
			Status:    status3.P,
			Permission: identity.Permission{
				Roles: args.Roles,
			},

			FullName: args.Name,
		}
		if err = s.identityAggr.Dispatch(ctx, cmd); err != nil {
			return err
		}
	default:
		return err
	}

	return nil
}
