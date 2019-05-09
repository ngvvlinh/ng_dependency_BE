package kiotviet

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	resty "gopkg.in/resty.v1"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/l"
	"etop.vn/backend/pkg/etop/model"
)

// Constants
const (
	DefaultPageSize = 100
	BaseURLID       = "https://id.kiotviet.vn/"
	BaseURLPublic   = "https://public.kiotapi.com/"
)

var (
	ll  = l.New()
	cfg SyncConfig

	client = &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	rclient = resty.NewWithClient(client).SetDebug(true)
)

func init() {
	bus.AddHandler("kiotviet", RequestAccessToken)
	bus.AddHandler("kiotviet", GetOrRenewToken)
}

type SyncConfig struct {
	WebhookBaseURL string `yaml:"webhook_base_url" envconfig:"webhook_base_url"`

	skipWebhook bool
}

func Init(config SyncConfig) {
	cfg = config
	if cfg.WebhookBaseURL == "" {
		if cm.IsDev() {
			cfg.skipWebhook = true
			ll.Warn("kiotviet: Skip webhook")
			return
		}
		ll.Fatal("Missing webhook")
	}
	if !strings.HasSuffix(cfg.WebhookBaseURL, "/") {
		cfg.WebhookBaseURL = cfg.WebhookBaseURL + "/"
	}
}

func NewConnection(ctx context.Context, retailerID, clientID, clientSecret string) (*Connection, []*Branch, error) {
	if retailerID == "" {
		return nil, nil, cm.Error(cm.InvalidArgument, "Thiếu subdomain", nil)
	}
	if !validateInput(clientID) || !validateInput(clientSecret) {
		return nil, nil, cm.Error(cm.InvalidArgument, "Invalid Kiotviet credential", nil)
	}

	tokenCmd := &RequestAccessTokenCommand{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}
	if err := bus.Dispatch(ctx, tokenCmd); err != nil {
		return nil, nil, err
	}

	// Verify retailerID
	conn := &Connection{
		RetailerID: retailerID,
		TokenStr:   tokenCmd.Result.TokenStr,
		ExpiresAt:  tokenCmd.Result.ExpiresAt,
	}

	var branchResp BranchResult
	_, err := getx(conn, BaseURLPublic+"branches", &branchResp)
	if err != nil {
		return nil, nil, cm.Error(cm.InvalidArgument, "Unable to verify RetailerID", err)
	}
	return conn, branchResp.Data, nil
}

func RequestAccessToken(ctx context.Context, cmd *RequestAccessTokenCommand) error {
	if cmd.ClientID == "" || cmd.ClientSecret == "" {
		return errors.New("kiotviet: Missing credential")
	}

	const reqURL = BaseURLID + "connect/token"
	q := make(url.Values)
	q.Set("scopes", "PublicApi.Access")
	q.Set("grant_type", "client_credentials")
	q.Set("client_id", cmd.ClientID)
	q.Set("client_secret", cmd.ClientSecret)

	req, err := http.NewRequest("POST", reqURL, strings.NewReader(q.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	now := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	var tokenResp AccessTokenResponse
	var errResp *ErrorResponse
	_, err = parseResponse(req, resp, &tokenResp, &errResp)
	if err != nil {
		if errResp != nil {
			return cm.Error(cm.ExternalServiceError, "Lỗi từ Kiotviet: "+errResp.Error(), errResp)
		}
		return err
	}

	// Try getting expiresAt from JWT first, then from response
	expiresAt := getJWTExpires(tokenResp.TokenStr)
	if expiresAt.IsZero() {
		expiresAt = now.Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	}

	cmd.Result.TokenStr = tokenResp.TokenStr
	cmd.Result.ExpiresAt = expiresAt
	return nil
}

func GetOrRenewToken(ctx context.Context, cmd *GetOrRenewTokenCommand) error {
	psQuery := &model.GetProductSourceExtendedQuery{
		GetProductSourceProps: model.GetProductSourceProps{
			ID: cmd.SupplierID,
		},
	}
	if err := bus.Dispatch(ctx, psQuery); err != nil {
		return cm.Error(cm.NotFound, "", err)
	}

	ps := psQuery.Result.ProductSource
	psi := psQuery.Result.ProductSourceInternal
	if psi == nil {
		return cm.Error(cm.FailedPrecondition, "ProductSource is not external", nil)
	}

	secret := psi.Secret
	if secret == nil || secret.RetailerID != ps.ExternalKey {
		return cm.Error(cm.Internal, "Mismatch Key", nil).
			Log("", l.Int64("ProductSource", ps.ID))
	}

	retailerID := secret.RetailerID
	conn := &Connection{
		RetailerID: retailerID,
		TokenStr:   psi.AccessToken,
		ExpiresAt:  psi.ExpiresAt,
	}
	cmd.Result.Connection = conn

	// Renew access token
	now := time.Now()
	if psi.AccessToken == "" || psi.ExpiresAt.Sub(now) < 30*time.Minute {
		tokenCmd := &RequestAccessTokenCommand{
			ClientID:     secret.ClientID,
			ClientSecret: secret.ClientSecret,
		}
		if err := bus.Dispatch(ctx, tokenCmd); err != nil {
			return err
		}
		conn.TokenStr = tokenCmd.Result.TokenStr
		conn.ExpiresAt = tokenCmd.Result.ExpiresAt

		updateTokenCmd := &model.UpdateKiotvietAccessTokenCommand{
			ProductSourceID: ps.ID,
			ClientToken:     conn.TokenStr,
			ExpiresAt:       conn.ExpiresAt,
		}
		if err := bus.Dispatch(ctx, updateTokenCmd); err != nil {
			return err
		}
	}
	return nil
}

func getOrRenewToken(ctx context.Context, supplierID int64, conn *Connection) (*Connection, error) {
	if conn == nil || conn.ExpiresAt.Sub(time.Now()) < 30*time.Minute {
		tokenCmd := &GetOrRenewTokenCommand{SupplierID: supplierID}
		if err := bus.Dispatch(ctx, tokenCmd); err != nil {
			return nil, err
		}

		newConn := &Connection{
			RetailerID: conn.RetailerID,
			TokenStr:   tokenCmd.Result.TokenStr,
			ExpiresAt:  tokenCmd.Result.ExpiresAt,
		}
		return newConn, nil
	}
	return conn, nil
}
