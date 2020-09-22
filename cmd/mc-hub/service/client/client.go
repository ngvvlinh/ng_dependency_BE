package client

import (
	"context"
	"net/http"
	"strings"
	"time"

	"o.o/api/top/external/mc/vnp"
	externaltypes "o.o/api/top/external/types"
	pbcm "o.o/api/top/types/common"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/httpreq"
	"o.o/backend/pkg/common/headers"
	"o.o/common/l"
)

var ll = l.New()

type EndpointConfig struct {
	BaseURL string `yaml:"base_url"`
}

// Validate implements a simple validation for config
func (c EndpointConfig) Validate() error {
	if c.BaseURL == "" ||
		!strings.HasPrefix(c.BaseURL, "http") {
		return cm.Error(cm.Internal, "invalid base_url", nil)
	}
	return nil
}

var _ vnp.ShipnowService = &ShipnowClient{}

type ShipnowClient struct {
	EndpointConfig

	Client *httpreq.Resty
}

func New(endpoint EndpointConfig) (*ShipnowClient, error) {
	if err := endpoint.Validate(); err != nil {
		return nil, err
	}
	client := httpreq.NewResty(
		httpreq.RestyConfig{
			Client: &http.Client{
				Timeout: 10 * time.Second,
			},
		},
	)
	s := &ShipnowClient{
		EndpointConfig: endpoint,
		Client:         client,
	}
	s.BaseURL += "/v1"
	return s, nil
}

func (s *ShipnowClient) Clone() vnp.ShipnowService {
	res := *s
	return &res
}

func (s *ShipnowClient) Ping(ctx context.Context, req *pbcm.Empty) (resp *pbcm.Empty, _ error) {
	err := s.sendRequest(ctx, vnp.Path_Shipnow_Ping, req, &resp)
	return resp, err
}

func (s *ShipnowClient) GetShipnowServices(ctx context.Context, req *externaltypes.GetShipnowServicesRequest) (resp *externaltypes.GetShipnowServicesResponse, _ error) {
	err := s.sendRequest(ctx, vnp.Path_Shipnow_GetShipnowServices, req, &resp)
	return resp, err
}

func (s *ShipnowClient) CreateShipnowFulfillment(ctx context.Context, req *externaltypes.CreateShipnowFulfillmentRequest) (resp *externaltypes.ShipnowFulfillment, _ error) {
	err := s.sendRequest(ctx, vnp.Path_Shipnow_CreateShipnowFulfillment, req, &resp)
	return resp, err
}

func (s *ShipnowClient) CancelShipnowFulfillment(ctx context.Context, req *externaltypes.CancelShipnowFulfillmentRequest) (resp *pbcm.UpdatedResponse, _ error) {
	err := s.sendRequest(ctx, vnp.Path_Shipnow_CancelShipnowFulfillment, req, &resp)
	return resp, err
}

func (s *ShipnowClient) GetShipnowFulfillment(ctx context.Context, req *externaltypes.FulfillmentIDRequest) (resp *externaltypes.ShipnowFulfillment, _ error) {
	err := s.sendRequest(ctx, vnp.Path_Shipnow_GetShipnowFulfillment, req, &resp)
	return resp, err
}

func (s *ShipnowClient) sendRequest(ctx context.Context, path string, req, resp interface{}) error {
	r := s.Client.R().SetContext(ctx)
	r.SetHeader("Content-Type", "application/json")
	r.SetAuthToken(headers.GetAuthorizationToken(ctx)) // TODO(vu): use session
	r.ExpectContentType("application/json")
	r.SetBody(req).SetResult(resp).SetError(&pbcm.Error{})
	httpResp, err := r.Post(s.BaseURL + path)
	if err != nil {
		return cm.Errorf(cm.Internal, err, "can not call shipnow service")
	}
	// TODO(vu): move into a lib
	if httpResp.IsError() {
		ll.S.Infof("%s", httpResp.Body())
		errMsg := httpResp.Error().(*pbcm.Error)
		switch errMsg.Code {
		case cm.NotFound.String():
			return cm.Errorf(cm.NotFound, nil, "%v", errMsg.Msg)
		case cm.InvalidArgument.String():
			return cm.Errorf(cm.InvalidArgument, nil, "%v", errMsg.Msg)
		default:
			return cm.Errorf(cm.Internal, nil, "%v", errMsg.Msg)
		}
	}
	return nil
}
