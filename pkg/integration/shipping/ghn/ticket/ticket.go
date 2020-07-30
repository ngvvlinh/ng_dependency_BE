package ticket

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/schema"
	"gopkg.in/resty.v1"

	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/httpreq"
	"o.o/backend/pkg/common/cmenv"
	"o.o/common/l"
)

var ll = l.New()
var encoder = schema.NewEncoder()

type Client struct {
	userID  int
	baseUrl string
	token   string

	rclient *httpreq.Resty
}

func New(env string) *Client {
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	rcfg := httpreq.RestyConfig{Client: client}
	// TODO(Nam) config
	c := &Client{
		token:   "f2de18e8-bb22-11e9-8485-7afeb53b4264",
		userID:  2043022936363,
		rclient: httpreq.NewResty(rcfg),
	}
	switch env {
	case cmenv.PartnerEnvTest, cmenv.PartnerEnvDev:
		c.baseUrl = "https://online-gateway.ghn.vn/shiip/public-ap"
	case cmenv.PartnerEnvProd:
		c.baseUrl = "https://console.ghn.vn/api/v1/apiv3/"
	default:
		ll.Fatal("ghn: Invalid env")
	}
	return c
}

func (c *Client) CreateTicket(ctx context.Context, req *CreateTicketRequest) (*CreateTicketResponse, error) {
	req.CEmail = "etop@etop.vn" //TODO(Nam) config
	var resp CreateTicketResponse
	fullUrl := c.baseUrl + "/ticket/create"
	if err := c.sendPostFormRequest(ctx, "post", fullUrl, req, &resp, "Không thể tạo ticket"); err != nil {
		return nil, err
	}
	return &resp, nil
}

// reply in ghn = comment etop
func (c *Client) CreateReply(ctx context.Context, req *CreateTicketReplyRequest) (*CreateTicketReplyResponse, error) {
	req.UserID = "2043053874563" //TODO(Nam) Config
	var resp CreateTicketReplyResponse
	fullUrl := c.baseUrl + "/ticket/reply"
	if err := c.sendPostFormRequest(ctx, "put", fullUrl, req, &resp, "Không thể tạo ticket"); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) sendPostFormRequest(ctx context.Context, method string, fullURL string, req interface{}, resp interface{}, msg string) error {
	values := url.Values{}
	if req != nil {
		if err := encoder.Encode(req, values); err != nil {
			return err
		}
	}
	var formData = make(map[string]string)
	for key := range values {
		formData[key] = values.Get(key)
	}
	var res *resty.Response
	var err error
	switch method {
	case "post":
		res, err = c.rclient.R().SetFormData(formData).Post(fullURL)
	case "put":
		res, err = c.rclient.R().SetFormData(formData).Put(fullURL)
	default:
		return cm.Errorf(cm.ExternalServiceError, err, "Method not found")
	}
	if err != nil {
		return cm.Errorf(cm.ExternalServiceError, err, "Lỗi kết nối với GHN")
	}
	return httpreq.HandleResponse(ctx, res, resp, msg)
}
