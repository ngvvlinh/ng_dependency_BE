package ahamove

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"o.o/api/main/identity"
	"o.o/api/main/shipnow/carrier"
	shipnowcarrier "o.o/backend/com/main/shipnow-carrier"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/integration/shipnow/ahamove/client"
	"o.o/capi/dot"
)

var _ shipnowcarrier.ShipnowCarrierAccount = &CarrierAccount{}

type URLConfig struct {
	ThirdPartyHost       string
	PathUserVerification string
}

type CarrierAccount struct {
	client        *client.Client
	urlConfig     URLConfig
	IdentityQuery identity.QueryBus
}

func NewCarrierAccount(ahamoveClient *client.Client, urlConfig URLConfig, identityBus identity.QueryBus) *CarrierAccount {
	ca := &CarrierAccount{
		client:        ahamoveClient,
		urlConfig:     urlConfig,
		IdentityQuery: identityBus,
	}
	return ca
}

func (c *CarrierAccount) RegisterExternalAccount(ctx context.Context, args *shipnowcarrier.RegisterExternalAccountArgs) (*carrier.RegisterExternalAccountResult, error) {
	request := &client.RegisterAccountRequest{
		Mobile:  args.Phone,
		Name:    args.Name,
		Address: args.Address,
	}

	response, err := c.client.RegisterAccount(ctx, request)
	if err != nil {
		return nil, err
	}
	res := &carrier.RegisterExternalAccountResult{
		Token: response.Token,
	}
	return res, nil
}

func (c *CarrierAccount) GetExternalAccount(ctx context.Context, args *shipnowcarrier.GetExternalAccountArgs) (*carrier.ExternalAccount, error) {
	token, err := getToken(ctx, c.IdentityQuery, args.OwnerID)
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Token không được để trống. Vui lòng tạo tài khoản Ahamove")
	}

	request := &client.GetAccountRequest{
		Token: token,
	}
	account, err := c.client.GetAccount(ctx, request)
	if err != nil {
		return nil, err
	}

	createAt := time.Unix(int64(account.CreateTime), 0)
	res := &carrier.ExternalAccount{
		ID:        account.ID,
		Name:      account.Name,
		Email:     account.Email,
		Verified:  account.Verified,
		CreatedAt: createAt,
	}
	return res, nil
}

func (c *CarrierAccount) VerifyExternalAccount(ctx context.Context, args *shipnowcarrier.VerifyExternalAccountArgs) (*carrier.VerifyExternalAccountResult, error) {
	token, err := getToken(ctx, c.IdentityQuery, args.OwnerID)
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Token không được để trống. Vui lòng tạo tài khoản Ahamove")
	}

	description, err := getDescriptionForVerification(ctx, c.IdentityQuery, c.urlConfig, args.OwnerID)
	if err != nil {
		return nil, err
	}

	request := &client.VerifyAccountRequest{
		Token:       token,
		Description: description,
	}
	response, err := c.client.VerifyAccount(ctx, request)
	if err != nil {
		return nil, err
	}

	externalTicket := response.Ticket
	res := &carrier.VerifyExternalAccountResult{
		TicketID:    strconv.Itoa(externalTicket.ID),
		Subject:     externalTicket.Subject,
		Description: externalTicket.Description,
		CreatedAt:   externalTicket.CreatedAt,
	}
	return res, nil
}

func prepareAhamovePhotoUrl(
	urlConfig URLConfig,
	ahamoveAccount *identity.ExternalAccountAhamove,
	uri string, typeImg string,
) string {
	ext := filepath.Ext(uri)
	filename := strings.TrimSuffix(filepath.Base(uri), ext)
	newName := fmt.Sprintf("user_%v_%v_%v", typeImg, ahamoveAccount.ExternalID, ahamoveAccount.ExternalCreatedAt.Unix())

	// example result:
	// https://example.com/ahamove/user_verification/BdVzaWz6ssamNKrRV7W8/user_id_front_84909090999_1444118656.jpg
	return fmt.Sprintf("%v%v/%v/%v%v",
		urlConfig.ThirdPartyHost,
		urlConfig.PathUserVerification,
		filename, newName, ext)
}

// description format: <user._id>, <user.name>, <photo_urls>
// photo_url format: <topship_domain>/upload/ahamove/user_verification/user_id_front<user.id>_<user.create_time>.jpg
func getDescriptionForVerification(ctx context.Context, identityQuery identity.QueryBus, urlConfig URLConfig, userID dot.ID) (des string, _err error) {
	queryUser := &identity.GetUserByIDQuery{
		UserID: userID,
	}
	if err := identityQuery.Dispatch(ctx, queryUser); err != nil {
		return "", err
	}
	user := queryUser.Result

	query := &identity.GetExternalAccountAhamoveQuery{
		Phone:   user.Phone,
		OwnerID: user.ID,
	}
	if err := identityQuery.Dispatch(ctx, query); err != nil {
		return "", err
	}
	account := query.Result

	var photoImgs []string
	front := prepareAhamovePhotoUrl(urlConfig, account, account.IDCardFrontImg, "id_front")
	back := prepareAhamovePhotoUrl(urlConfig, account, account.IDCardBackImg, "id_back")
	portrait := prepareAhamovePhotoUrl(urlConfig, account, account.PortraitImg, "portrait")

	photoImgs = append(photoImgs, front, back, portrait)
	photoImgs = append(photoImgs, account.CompanyImgs...)
	photoImgs = append(photoImgs, account.BusinessLicenseImgs...)
	if account.FanpageURL != "" {
		photoImgs = append(photoImgs, account.FanpageURL)
	}
	if account.WebsiteURL != "" {
		photoImgs = append(photoImgs, account.WebsiteURL)
	}

	des = fmt.Sprintf("%v, %v, %v", account.ExternalID, account.Name, strings.Join(photoImgs, ", "))
	return des, nil
}
