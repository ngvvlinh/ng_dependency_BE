package root

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"

	"o.o/api/main/authorization"
	"o.o/api/main/identity"
	api "o.o/api/top/int/etop"
	"o.o/api/top/types/etc/status3"
	identitymodel "o.o/backend/com/main/identity/model"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/code/gencode"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/common/validate"
	"o.o/capi/dot"
)

var keyWebphoneRequestLogin = "webphone-request-login"

// create an user and a default shop for vnpost
const (
	vnpostShopID = 1134164111674536521
	vnpostKey    = "vnpost"
)

func (s *UserService) WebphoneRequestLogin(ctx context.Context, r *api.WebphoneRequestLoginRequest) (*api.WebphoneRequestLoginResponse, error) {
	phone, ok := validate.NormalizePhone(r.Phone)
	if !ok {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Số điện thoại không hợp lệ")
	}

	// send secret_key
	key := fmt.Sprintf("%v-%v", keyWebphoneRequestLogin, phone.String())
	code, err := s.RedisStore.GetString(key)
	if err != nil || code == "" {
		code = gencode.GenerateCode(gencode.Alphabet54, 24)
		_ = s.RedisStore.SetStringWithTTL(key, code, 5*60)
	}
	return &api.WebphoneRequestLoginResponse{
		SecretKey: code,
	}, nil
}

func (s *UserService) WebphoneLogin(ctx context.Context, r *api.WebphoneLoginRequest) (*api.LoginResponse, error) {
	phone, err := s.ValidateVNPostWebphoneLoginRequest(r)
	if err != nil {
		return nil, err
	}
	cmd := &identity.RegisterSimplifyCommand{
		Phone: phone,
	}
	if err := s.IdentityAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	query := &identitymodelx.GetUserByEmailOrPhoneQuery{
		Phone: phone,
	}
	if err := s.UserStoreIface.GetUserByEmailOrPhone(ctx, query); err != nil {
		return nil, err
	}
	user := query.Result
	isExisted, err := s.checkIfAccountUserWithVNPostShopExist(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	if !isExisted {
		createAccountUserCmd := &identitymodelx.CreateAccountUserCommand{
			AccountUser: &identitymodel.AccountUser{
				AccountID: vnpostShopID,
				UserID:    user.ID,
				Status:    status3.P,
				Permission: identitymodel.Permission{
					Roles: []string{authorization.RoleSalesMan.String()},
				},
			},
		}
		if err := s.AccountUserStore.CreateAccountUser(ctx, createAccountUserCmd); err != nil {
			return nil, err
		}
	}

	resp, err := s.CreateLoginResponse(ctx, nil, "", user.ID, user, vnpostShopID, 0, true, 0)
	if err != nil {
		return nil, err
	}
	// Only return token of default vnpost shop (vnpostShopID)
	resp.AvailableAccounts = nil
	return resp, nil
}

func (s *UserService) ValidateVNPostWebphoneLoginRequest(r *api.WebphoneLoginRequest) (_phone string, err error) {
	phone, ok := validate.NormalizePhone(r.Phone)
	if !ok {
		return "", cm.Errorf(cm.InvalidArgument, nil, "Số điện thoại không hợp lệ")
	}
	key := fmt.Sprintf("%v-%v", keyWebphoneRequestLogin, phone.String())
	secretKey, err := s.RedisStore.GetString(key)
	if err != nil && err != redis.ErrNil {
		return "", cm.Errorf(cm.Internal, err, "")
	}

	// data = phone + public_key + "vnpost"
	data := phone.String() + string(s.WebphonePublicKey) + vnpostKey
	checkSum, err := CheckSum(data, secretKey)
	if err != nil {
		return "", err
	}
	if checkSum == r.Code {
		return phone.String(), nil
	}
	return "", cm.Errorf(cm.InvalidArgument, nil, "Mã code không hợp lệ")
}

func CheckSum(data, secretKey string) (string, error) {
	hash := hmac.New(sha1.New, []byte(secretKey))
	_, err := hash.Write([]byte(data))
	if err != nil {
		return "", cm.Errorf(cm.Internal, err, "")
	}
	sum := hash.Sum(nil)
	return base64.StdEncoding.EncodeToString(sum), nil
}

func (s *UserService) checkIfAccountUserWithVNPostShopExist(ctx context.Context, userID dot.ID) (bool, error) {
	query := &identitymodelx.GetAccountUserQuery{
		UserID:    userID,
		AccountID: vnpostShopID,
	}
	err := s.AccountUserStore.GetAccountUser(ctx, query)
	switch cm.ErrorCode(err) {
	case cm.NotFound:
		return false, nil
	case cm.NoError:
		return true, nil
	default:
		return false, err
	}

}
