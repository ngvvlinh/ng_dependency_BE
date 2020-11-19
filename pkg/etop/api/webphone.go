package api

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"

	"o.o/api/main/identity"
	api "o.o/api/top/int/etop"
	"o.o/api/top/types/etc/account_type"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	cm "o.o/backend/pkg/common"
	wldriver "o.o/backend/pkg/common/apifw/whitelabel/drivers"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/code/gencode"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/common/validate"
	"o.o/backend/pkg/etop/authorize/claims"
)

var keyWebphoneRequestLogin = "webphone-request-login"
var wlPartnerKeys = []string{wldriver.VNPostKey}

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
	phone, wlPartnerKey, err := s.ValidateWLPartnerWebphoneLoginRequest(r)
	if err != nil {
		return nil, err
	}
	// Wrap context to wl partner
	wlPartner := wl.GetWLPartnerByKey(wlPartnerKey)
	wlPartnerID := wlPartner.ID
	ctx = wl.WrapContextByPartnerID(ctx, wlPartnerID)
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

	// Gán wl_partner_id vào session
	// Các request sử dụng session này đều được tính cho wl_partner_id này
	claim := &claims.ClaimInfo{
		WLPartnerID: wlPartnerID,
	}
	resp, err := s.CreateLoginResponse(ctx, claim, "", user.ID, user, 0, account_type.Shop.Enum(), true, 0)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *UserService) ValidateWLPartnerWebphoneLoginRequest(r *api.WebphoneLoginRequest) (_phone, partnerKey string, err error) {
	phone, ok := validate.NormalizePhone(r.Phone)
	if !ok {
		return "", "", cm.Errorf(cm.InvalidArgument, nil, "Số điện thoại không hợp lệ")
	}
	key := fmt.Sprintf("%v-%v", keyWebphoneRequestLogin, phone.String())
	secretKey, err := s.RedisStore.GetString(key)
	if err != nil && err != redis.ErrNil {
		return "", "", cm.Errorf(cm.Internal, err, "")
	}

	// data = phone + public_key + wl_partner_key
	for _, partnerKey := range wlPartnerKeys {
		data := phone.String() + string(s.WebphonePublicKey) + partnerKey
		checkSum, err := CheckSum(data, secretKey)
		if err != nil {
			return "", "", err
		}
		if checkSum == r.Code {
			return phone.String(), partnerKey, nil
		}
	}
	return "", "", cm.Errorf(cm.InvalidArgument, nil, "Mã code không hợp lệ")
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
