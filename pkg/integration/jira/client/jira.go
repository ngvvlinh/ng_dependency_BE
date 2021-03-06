package client

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"gopkg.in/resty.v1"

	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/httpreq"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	cc "o.o/backend/pkg/common/config"
	"o.o/common/jsonx"
	"o.o/common/l"
	"o.o/common/xerrors"
)

type sendRequestArgs struct {
	url    string
	req    interface{}
	resp   interface{}
	method httpreq.RequestMethod
}

var ll = l.New()

type Config struct {
	UserEmail string `yaml:"user_email"    valid:"required"`
	APIKey    string `yaml:"api_key"     valid:"required"`
	BaseURL   string `yaml:"base_url" `
}

type Client struct {
	cfg     Config
	token   string
	rClient *httpreq.Resty
}

func (c *Config) MustLoadEnv(prefix ...string) {
	p := "ET_JIRA"
	if len(prefix) > 0 {
		p = prefix[0]
	}
	cc.EnvMap{
		p + "_EMAIL":    &c.UserEmail,
		p + "_BASE_URL": &c.BaseURL,
		p + "_API_KEY":  &c.APIKey,
	}.MustLoad()
}

func New(cfg Config) *Client {
	APITokenData := cfg.UserEmail + ":" + cfg.APIKey
	token := base64.StdEncoding.EncodeToString([]byte(APITokenData))

	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	rcfg := httpreq.RestyConfig{Client: client}

	return &Client{
		cfg:     cfg,
		token:   token,
		rClient: httpreq.NewResty(rcfg),
	}
}

func (c *Client) CreateIssue(ctx context.Context, requestData *CreateIssueRequest) (*CreateIssueResponse, error) {
	var resp CreateIssueResponse
	data, err := requestData.ToJiraBodyRequest()
	if err != nil {
		return nil, err
	}
	err = c.sendPostRequest(ctx, sendRequestArgs{
		url:    URL(c.cfg.BaseURL, "/rest/api/3/issue"),
		req:    data,
		resp:   &resp,
		method: httpreq.MethodPost,
	})
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) sendPostRequest(ctx context.Context, args sendRequestArgs) error {
	return c.sendRequest(ctx, args)
}

func (c *Client) sendRequest(ctx context.Context, args sendRequestArgs) error {
	var res *resty.Response
	var err error
	var errResp ErrorResponse

	request := c.rClient.R().
		SetHeader("Authorization", "Basic "+c.token).SetResult(&errResp).SetError(&errResp)

	switch args.method {
	case httpreq.MethodPost:
		res, err = request.SetBody(args.req).Post(args.url)
	default:
		panic(fmt.Sprintf("unsupported method %v", args.method))
	}
	if err != nil {
		return cm.Error(cm.ExternalServiceError, "L???i k???t n???i v???i Jira", err)
	}

	return handleResponse(ctx, res, args.resp, "Kh??ng th??? t???o issue", &errResp)

}

func handleResponse(ctx context.Context, res *httpreq.RestyResponse, result interface{}, errMsg string, errResp *ErrorResponse) error {
	status := res.StatusCode()
	body := res.Body()
	switch {
	case status >= 200 && status < 300:
		if result != nil {
			if httpreq.IsNullJsonRaw(body) {
				return cm.Errorf(cm.ExternalServiceError, nil, "L???i kh??ng x??c ?????nh t??? Jira: null response.")
			}
			if err := jsonx.Unmarshal(body, result); err != nil {
				return cm.Errorf(cm.ExternalServiceError, err, "L???i kh??ng x??c ?????nh t??? Jira: %v", err)
			}
		}
		return nil

	case status >= 400:
		var meta map[string]string
		var errJSON xerrors.ErrorJSON
		if !httpreq.IsNullJsonRaw(body) {
			if err := jsonx.Unmarshal(body, &meta); err != nil {
				var metaX map[string]interface{}
				_ = jsonx.Unmarshal(body, &metaX)
				meta = make(map[string]string)
				for k, v := range metaX {
					meta[k] = fmt.Sprint(v)
				}
			}
			errJSON.Msg = errMsg
			errJSON.Meta = errResp.Errors
		}

		return cm.Errorf(cm.ExternalServiceError, nil, "L???i t??? Jira: %v. N???u c???n th??m th??ng tin vui l??ng li??n h??? %v.", errJSON.Error(), wl.X(ctx).CSEmail).WithMetaM(meta)
	default:
		return cm.Errorf(cm.ExternalServiceError, nil, "L???i kh??ng x??c ?????nh t??? Jira: %v. Invalid status (%v).", errMsg, status)
	}
}
