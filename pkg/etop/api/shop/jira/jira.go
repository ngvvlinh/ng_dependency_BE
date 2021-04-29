package jira

import (
	"context"

	api "o.o/api/top/int/shop"
	cm "o.o/api/top/types/common"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/backend/pkg/integration/jira/driver"
)

type JiraService struct {
	session.Session
	Driver *driver.JiraDriver
}

func (s *JiraService) Clone() api.JiraService { res := *s; return &res }

func (s *JiraService) GetJiraCustomFields(ctx context.Context, r *cm.Empty) (*api.GetCustomFieldsResponse, error) {
	customFields := s.Driver.ListCustomFields()
	result := Convert_jira_CustomFields_to_api_CustomFields(customFields)
	return &api.GetCustomFieldsResponse{
		CustomFields: result,
	}, nil
}

func (s *JiraService) CreateJiraIssue(ctx context.Context, r *api.CreateJiraIssueRequest) (*cm.Empty, error) {
	customFields := make([]*driver.CustomField, 0, len(r.CustomFields))
	for _, field := range r.CustomFields {
		customFields = append(customFields, &driver.CustomField{
			Key:   field.Key,
			Value: field.Value,
		})
	}
	_, err := s.Driver.CreateIssue(ctx, &driver.CreateIssueRequest{
		Summary:      r.Summary,
		Description:  r.Description,
		CustomFields: customFields,
	})
	if err != nil {
		return nil, err
	}
	return &cm.Empty{}, nil
}
