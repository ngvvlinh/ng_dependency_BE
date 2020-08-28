package ahamove

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"o.o/api/main/accountshipnow"
	"o.o/api/main/identity"
	"o.o/api/main/shipnow/carrier"
	shipnowcarriertypes "o.o/backend/com/main/shipnow/carrier/types"
	"o.o/backend/pkg/integration/shipnow/ahamove/client"
	"o.o/capi/dot"
)

func (c *Carrier) RegisterExternalAccount(ctx context.Context, args *shipnowcarriertypes.RegisterExternalAccountArgs) (*carrier.RegisterExternalAccountResult, error) {
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

func (c *Carrier) GetExternalAccount(ctx context.Context, args *shipnowcarriertypes.GetExternalAccountArgs) (*carrier.ExternalAccount, error) {
	request := &client.GetAccountRequest{}
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

func (c *Carrier) VerifyExternalAccount(ctx context.Context, args *shipnowcarriertypes.VerifyExternalAccountArgs) (*carrier.VerifyExternalAccountResult, error) {
	description, err := c.getDescriptionForVerification(ctx, c.urlConfig, args.OwnerID)
	if err != nil {
		return nil, err
	}

	request := &client.VerifyAccountRequest{
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
	ahamoveAccount *accountshipnow.ExternalAccountAhamove,
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
func (c *Carrier) getDescriptionForVerification(ctx context.Context, urlConfig URLConfig, userID dot.ID) (des string, _err error) {
	queryUser := &identity.GetUserByIDQuery{
		UserID: userID,
	}
	if err := c.identityQuery.Dispatch(ctx, queryUser); err != nil {
		return "", err
	}
	user := queryUser.Result

	query := &accountshipnow.GetExternalAccountAhamoveQuery{
		Phone:   user.Phone,
		OwnerID: user.ID,
	}
	if err := c.accountshipnowQuery.Dispatch(ctx, query); err != nil {
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
