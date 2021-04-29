package jira

import (
	api "o.o/api/top/int/shop"
	"o.o/backend/pkg/integration/jira/driver"
)

func Convert_jira_CustomField_to_api_CustomField(args *driver.CustomField) *api.JiraCustomField {
	res := &api.JiraCustomField{
		Name: args.Name,
		Key:  args.Key,
	}
	return res
}

func Convert_jira_CustomFields_to_api_CustomFields(args []*driver.CustomField) []*api.JiraCustomField {
	res := make([]*api.JiraCustomField, len(args))
	for i, item := range args {
		res[i] = Convert_jira_CustomField_to_api_CustomField(item)
	}
	return res
}
