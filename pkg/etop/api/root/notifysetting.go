package root

import (
	"context"

	"o.o/api/main/notify"
	api "o.o/api/top/int/etop"
	pbcm "o.o/api/top/types/common"
	"o.o/backend/com/eventhandler/notifier"
	cm "o.o/backend/pkg/common"
)

func (s *UserService) GetNotifySetting(ctx context.Context, _ *pbcm.Empty) (*api.GetNotifySettingResponse, error) {
	userID := s.SS.User().ID
	cmd := &notify.GetOrCreateUserNotifySettingCommand{
		UserID:        userID,
		DisableTopics: []string{},
	}
	err := s.NotifyAggr.Dispatch(ctx, cmd)
	if err != nil {
		return nil, err
	}
	setting := cmd.Result
	return &api.GetNotifySettingResponse{
		UserID: setting.UserID,
		Topics: buildTopics(setting.DisableTopics),
	}, nil
}

func buildTopics(disableTopics []string) []*api.NotifyTopic {
	var res []*api.NotifyTopic
	for _, topic := range notifier.NotifyTopics {
		res = append(res, &api.NotifyTopic{
			Topic:  topic,
			Enable: !isContainTopic(topic, disableTopics),
		})
	}
	return res
}

func isContainTopic(topic string, disableTopics []string) bool {
	for _, _topic := range disableTopics {
		if _topic == topic {
			return true
		}
	}
	return false
}

func (s *UserService) EnableNotifyTopic(
	ctx context.Context,
	request *api.UpdateNotifyTopicRequest,
) (*api.GetNotifySettingResponse, error) {
	if !isContainTopic(request.Topic, notifier.NotifyTopics) {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "invalid topic")
	}

	userID := s.SS.User().ID
	cmd := &notify.GetOrCreateUserNotifySettingCommand{
		UserID:        userID,
		DisableTopics: []string{},
	}
	if err := s.NotifyAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	enableTopicCmd := &notify.EnableTopicCommand{
		UserID: userID,
		Topic:  request.Topic,
	}
	if err := s.NotifyAggr.Dispatch(ctx, enableTopicCmd); err != nil {
		return nil, err
	}

	setting := enableTopicCmd.Result
	return &api.GetNotifySettingResponse{
		UserID: setting.UserID,
		Topics: buildTopics(setting.DisableTopics),
	}, nil
}

func (s *UserService) DisableNotifyTopic(
	ctx context.Context,
	request *api.UpdateNotifyTopicRequest,
) (*api.GetNotifySettingResponse, error) {
	if !isContainTopic(request.Topic, notifier.NotifyTopics) {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "invalid topic")
	}

	userID := s.SS.User().ID
	cmd := &notify.GetOrCreateUserNotifySettingCommand{
		UserID:        userID,
		DisableTopics: []string{},
	}
	if err := s.NotifyAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	disableNotifyTopic := &notify.DisableTopicCommand{
		UserID: userID,
		Topic:  request.Topic,
	}
	if err := s.NotifyAggr.Dispatch(ctx, disableNotifyTopic); err != nil {
		return nil, err
	}

	setting := disableNotifyTopic.Result
	return &api.GetNotifySettingResponse{
		UserID: setting.UserID,
		Topics: buildTopics(setting.DisableTopics),
	}, nil
}
