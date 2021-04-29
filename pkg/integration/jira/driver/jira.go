package driver

import (
	"context"

	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/integration/jira/client"
	"o.o/common/l"
)

// Lấy projectKey và issueTypeID từ api của Jira
// GET "https://e-b2b.atlassian.net/rest/api/3/issue/createmeta"
// projectKeys là key của project trên Jira
const jiraProjectKey = "ELM"
const jiraIssueTypeID = "10059" // Type: Lead

var ll = l.New()

type JiraDriver struct {
	jiraClient *client.Client
}

func New(cfg client.Config) *JiraDriver {
	c := JiraDriver{
		jiraClient: client.New(cfg),
	}
	return &c
}

type CustomField struct {
	Name    string
	Key     string
	KeyJira string
	Value   string
}

// Lấy custom fields từ api của Jira
// GET "https://e-b2b.atlassian.net/rest/api/3/issue/createmeta?projectKeys=ELM&expand=projects.issuetypes.fields"
// projectKeys là key của project trên Jira
var JiraProjectCustomFields = map[string]CustomField{
	"representator": {
		Name:    "Người đại diện",
		Key:     "representator",
		KeyJira: "customfield_10062",
	},
	"env": {
		Name:    "ENV",
		Key:     "env",
		KeyJira: "customfield_10080",
	},
	"platform": {
		Name:    "Platform",
		Key:     "platform",
		KeyJira: "customfield_10079",
	},
	"app_bundle": {
		Name:    "App bundle",
		Key:     "app_bundle",
		KeyJira: "customfield_10078",
	},
	"domain": {
		Name:    "Domain",
		Key:     "domain",
		KeyJira: "customfield_10077",
	},
	"product": {
		Name:    "Product",
		Key:     "product",
		KeyJira: "customfield_10076",
	},
	"sales_phone": {
		Name:    "Số điện thoại Sales phụ trách",
		Key:     "sales_phone",
		KeyJira: "customfield_10081",
	},
	"email": {
		Name:    "Email",
		Key:     "email",
		KeyJira: "customfield_10070",
	},
	"customer_phone": {
		Name:    "Số điện thoại khách hàng",
		Key:     "customer_phone",
		KeyJira: "customfield_10063",
	},
}

type CreateIssueRequest struct {
	Summary      string
	Description  string
	CustomFields []*CustomField
}

func (d *JiraDriver) ListCustomFields() (res []*CustomField) {
	for _, customField := range JiraProjectCustomFields {
		res = append(res, &CustomField{
			Key:  customField.Key,
			Name: customField.Name,
		})
	}
	return
}

func (c *JiraDriver) CreateIssue(ctx context.Context, req *CreateIssueRequest) (res *client.CreateIssueResponse, _err error) {
	mapCustomField := make(map[string]string)
	for _, field := range req.CustomFields {
		if _, ok := mapCustomField[field.Key]; ok {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Duplicate custom field key: %v", field.Key)
		}
		cField, ok := JiraProjectCustomFields[field.Key]
		if !ok {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Custom field key %v does not exist", field.Key)
		}
		mapCustomField[cField.KeyJira] = field.Value
	}

	DescriptionContents := []client.DescriptionContent{
		{
			Type: "paragraph",
			Content: []client.Content{
				{
					Type: "text",
					Text: req.Description,
				},
			},
		},
	}
	data := client.CreateIssueRequest{
		IssueDefaultFields: client.IssueDefaultFields{
			IssueType: client.IssueType{
				ID: jiraIssueTypeID,
			},
			Project: client.Project{
				Key: jiraProjectKey,
			},
			Summary: req.Summary,
			Description: client.Description{
				Type:               "doc",
				Version:            1,
				DescriptionContent: DescriptionContents,
			},
		},
		CustomFields: mapCustomField,
	}
	return c.jiraClient.CreateIssue(ctx, &data)
}
