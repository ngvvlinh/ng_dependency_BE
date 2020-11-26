package sadmin

import (
	"context"
	"fmt"

	"o.o/api/top/int/sadmin"
	"o.o/api/top/types/common"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/common/l"
)

const (
	PrefixWebhookCallbackURL = "WebhookCallbackURL"
	Version                  = "v1.0"
	OneDay                   = 24 * 60 * 60
)

var ll = l.New()

type WebhookCallbackService struct {
	rd redis.Store
}

func NewWebhookCallbackService(_rd redis.Store) *WebhookCallbackService {
	return &WebhookCallbackService{
		rd: _rd,
	}
}

type WebhookService struct {
	session.Session
	WebhookCallbackService *WebhookCallbackService
}

func (s *WebhookService) Clone() sadmin.WebhookService {
	res := *s
	return &res
}

type Callback struct {
	URL     string            `json:"url"`
	Options map[string]string `json:"options"`
}

func (s *WebhookService) RegisterWebhook(
	ctx context.Context, request *sadmin.SAdminRegisterWebhookRequest,
) (*common.Empty, error) {
	if request.CallbackURL == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "callback must not be empty")
	}
	if err := s.WebhookCallbackService.AddWebhookCallback(request.Type.String(), Callback{
		URL:     request.CallbackURL,
		Options: request.Options,
	}); err != nil {
		return nil, err
	}
	ll.SendMessagef("[sadmin] Webhook registered type: %s, callback_url: %s", request.Type.String(), request.CallbackURL)
	return &common.Empty{}, nil
}

func (s *WebhookService) UnregisterWebhook(
	ctx context.Context,
	req *sadmin.SAdminUnregisterWebhookRequest,
) (_ *common.Empty, _err error) {
	if req.CallbackURL == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "callback must not be empty")
	}
	if req.RemoveAll {
		_err = s.WebhookCallbackService.RemoveAllWebhookCallbackURLs(req.Type.String())
	} else {
		_err = s.WebhookCallbackService.RemoveWebhookCallbackURL(req.Type.String(), req.CallbackURL)
	}
	if _err != nil {
		return nil, _err
	}
	ll.SendMessagef("[sadmin] Webhook unregistered type: %s, callback_url: %s", req.Type.String(), req.CallbackURL)
	return &common.Empty{}, nil
}

func WebhookCallbackKey(webhookType string) string {
	return fmt.Sprintf("%v-%v-%v", PrefixWebhookCallbackURL, Version, webhookType)
}

func (s *WebhookCallbackService) AddWebhookCallback(webhookType string, newCallback Callback) error {
	callbackURLs, err := s.GetWebhookCallbacks(webhookType)
	if err != nil {
		return err
	}

	var isExist bool
	for _, callbackURL := range callbackURLs {
		if newCallback.URL == callbackURL.URL {
			isExist = true
			break
		}
	}

	if !isExist {
		callbackURLs = append(callbackURLs, newCallback)
	}

	return s.SetWebhookCallbacks(webhookType, callbackURLs)
}

func (s *WebhookCallbackService) RemoveAllWebhookCallbackURLs(webhookType string) error {
	key := WebhookCallbackKey(webhookType)
	switch err := s.rd.Del(key); err {
	case nil, redis.ErrNil:
		return nil
	default:
		return err
	}
}

func (s *WebhookCallbackService) RemoveWebhookCallbackURL(webhookType, oldCallbackURL string) error {
	callbacks, err := s.GetWebhookCallbacks(webhookType)
	if err != nil {
		return err
	}

	var newCallbacks []Callback
	for _, callback := range callbacks {
		if oldCallbackURL != callback.URL {
			newCallbacks = append(newCallbacks, callback)
		}
	}

	return s.SetWebhookCallbacks(webhookType, newCallbacks)
}

func (s *WebhookCallbackService) GetWebhookCallbacks(webhookType string) (callbacks []Callback, _ error) {
	key := WebhookCallbackKey(webhookType)
	err := s.rd.Get(key, &callbacks)
	switch err {
	case redis.ErrNil:
		return nil, nil
	case nil:
		return callbacks, nil
	default:
		return nil, err
	}
}

func (s *WebhookCallbackService) SetWebhookCallbacks(webhookType string, callbacks []Callback) error {
	key := WebhookCallbackKey(webhookType)
	if err := s.rd.SetWithTTL(key, callbacks, OneDay); err != nil {
		return err
	}
	return nil
}
