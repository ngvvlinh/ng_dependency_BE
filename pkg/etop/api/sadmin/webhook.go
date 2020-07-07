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
	ctx context.Context, request *sadmin.SAdminUnregisterWebhookRequest,
) (*common.Empty, error) {
	if request.CallbackURL == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "callback must not be empty")
	}

	if err := s.WebhookCallbackService.RemoveWebhookCallbackURL(request.Type.String(), request.CallbackURL); err != nil {
		return nil, err
	}
	ll.SendMessagef("[sadmin] Webhook registered type: %s, callback_url: %s", request.Type.String(), request.CallbackURL)
	return &common.Empty{}, nil
}

func (s *WebhookCallbackService) GetWebhookCallbackKey(webhookType string) string {
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
	key := s.GetWebhookCallbackKey(webhookType)
	callbackURLsJoined, err := s.rd.GetString(key)
	switch err {
	case redis.ErrNil:
		return []string{}, nil
	case nil:
		return strings.Split(callbackURLsJoined, Separator), nil
	default:
		return []string{}, err
	}
}

func (s *WebhookCallbackService) SetWebhookCallbackURLs(webhookType string, callbackURLs []string) error {
	key := s.GetWebhookCallbackKey(webhookType)
	if err := s.rd.SetStringWithTTL(key, strings.Join(callbackURLs, Separator), OneDay); err != nil {
		return err
	}
	return nil
}
