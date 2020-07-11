package sadmin

import (
	"context"
	"fmt"
	"strings"

	"o.o/api/top/int/sadmin"
	"o.o/api/top/types/common"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/common/l"
)

const (
	PrefixWebhookCallbackURL = "WebhookCallbackURL"
	OneDay                   = 24 * 60 * 60
	Separator                = ","
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

func (s *WebhookService) RegisterWebhook(
	ctx context.Context, request *sadmin.SAdminRegisterWebhookRequest,
) (*common.Empty, error) {
	if request.CallbackURL == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "callback must not be empty")
	}
	if err := s.WebhookCallbackService.AddWebhookCallbackURL(request.Type.String(), request.CallbackURL); err != nil {
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
	return fmt.Sprintf("%v-%v", PrefixWebhookCallbackURL, webhookType)
}

func (s *WebhookCallbackService) AddWebhookCallbackURL(webhookType, newCallbackURL string) error {
	callbackURLs, err := s.GetWebhookCallbackURLs(webhookType)
	if err != nil {
		return err
	}

	var isExist bool
	for _, callbackURL := range callbackURLs {
		if newCallbackURL == callbackURL {
			isExist = true
			break
		}
	}

	if !isExist {
		callbackURLs = append(callbackURLs, newCallbackURL)
	}

	return s.SetWebhookCallbackURLs(webhookType, callbackURLs)
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
	callbackURLs, err := s.GetWebhookCallbackURLs(webhookType)
	if err != nil {
		return err
	}

	var newCallbackURLs []string
	for _, callbackURL := range callbackURLs {
		if oldCallbackURL != callbackURL {
			newCallbackURLs = append(newCallbackURLs, callbackURL)
		}
	}

	return s.SetWebhookCallbackURLs(webhookType, newCallbackURLs)
}

func (s *WebhookCallbackService) GetWebhookCallbackURLs(webhookType string) ([]string, error) {
	key := WebhookCallbackKey(webhookType)
	callbackURLsJoined, err := s.rd.GetString(key)
	switch err {
	case redis.ErrNil:
		return nil, nil
	case nil:
		if callbackURLsJoined == "" {
			return nil, nil
		}
		return strings.Split(callbackURLsJoined, Separator), nil
	default:
		return nil, err
	}
}

func (s *WebhookCallbackService) SetWebhookCallbackURLs(webhookType string, callbackURLs []string) error {
	key := WebhookCallbackKey(webhookType)
	if err := s.rd.SetStringWithTTL(key, strings.Join(callbackURLs, Separator), OneDay); err != nil {
		return err
	}
	return nil
}
