package client

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
	"net/http"
	"o.o/api/top/int/etop"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/httpreq"
	"o.o/backend/pkg/common/authorization/auth"
	cc "o.o/backend/pkg/common/config"
	"o.o/common/l"
	"time"
)

var ll = l.New()

type Config struct {
	Issuer       string `yaml:"issuer"`
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	BaseURL      string `yaml:"base_url" `
	RedirectUrls struct {
		App string `yaml:"app"`
		Web string `yaml:"web"`
	} `yaml:"redirect_urls"`
}

func (c *Config) MustLoadEnv(prefix ...string) {
	p := "ET_OIDC"
	if len(prefix) > 0 {
		p = prefix[0]
	}
	cc.EnvMap{
		p + "_ISSUER":        &c.Issuer,
		p + "_CLIENT_ID":     &c.ClientID,
		p + "_CLIENT_SECRET": &c.ClientSecret,
		p + "_BASE_URL":      &c.BaseURL,
		p + "_REDIRECT_URLS": &c.RedirectUrls,
	}.MustLoad()
}

type Client struct {
	cfg          Config
	oAuth2Config *oauth2.Config
	provider     *oidc.Provider
	rClient      *httpreq.Resty
}

func New(cfg Config) *Client {
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	rcfg := httpreq.RestyConfig{Client: client}
	provider, err := oidc.NewProvider(context.Background(), cfg.Issuer)
	if err != nil {
		ll.Fatal("Invalid provider")
	}

	return &Client{
		cfg:      cfg,
		provider: provider,
		oAuth2Config: &oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURL:  cfg.RedirectUrls.Web,
			Endpoint:     provider.Endpoint(),
			Scopes:       []string{oidc.ScopeOpenID, oidc.ScopeOfflineAccess, "profile", "phone", "email", "role"},
		},
		rClient: httpreq.NewResty(rcfg),
	}
}

func (c *Client) GetAuthURL(redirectType etop.RedirectType) string {
	if redirectType == etop.APP {
		c.oAuth2Config.RedirectURL = c.cfg.RedirectUrls.App
	}
	return c.oAuth2Config.AuthCodeURL(auth.RandomToken(16))
}

func (c *Client) VerifyToken(ctx context.Context, code string) (*VerifyTokenResponse, error) {
	verifier := c.provider.Verifier(&oidc.Config{
		ClientID: c.cfg.ClientID,
	})

	token, err := c.oAuth2Config.Exchange(ctx, code)
	if err != nil {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Failed to exchange token: %v.", err)
	}

	rawToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Invalid token.")
	}

	idToken, err := verifier.Verify(ctx, rawToken)
	if err != nil {
		return nil, cm.Errorf(cm.FailedPrecondition, err, "Failed to verify token: %v", err)
	}

	var resp VerifyTokenResponse
	if err := idToken.Claims(&resp.IDTokenClaims); err != nil {
		return nil, cm.Errorf(cm.ExternalServiceError, err, "Invalid claims: %v", err)
	}

	claims, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(claims, &resp.IDTokenClaims); err != nil {
		return nil, err
	}

	return &resp, nil
}
