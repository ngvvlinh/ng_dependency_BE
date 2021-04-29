package client

import (
	"o.o/backend/pkg/common/cmreflect"
	"o.o/common/jsonx"
	"o.o/common/l"
)

type CreateIssueResponse struct {
	ID   string `json:"id"`
	Key  string `json:"key"`
	Self string `json:"self"`
}

type ErrorResponse struct {
	Errors        map[string]string `json:"errors"`
	ErrorMessages []string          `json:"errorMessages"`
}

type IssueDefaultFields struct {
	Summary     string      `json:"summary"`
	Description Description `json:"description"`
	IssueType   IssueType   `json:"issuetype"`
	Project     Project     `json:"project"`
}

type CreateIssueRequest struct {
	IssueDefaultFields
	CustomFields map[string]string
}

func (r *CreateIssueRequest) String() string { return jsonx.MustMarshalToString(r) }

func (f CreateIssueRequest) ToJiraBodyRequest() (map[string]interface{}, error) {
	data, err := cmreflect.EncodeStructToMap(f.IssueDefaultFields, "json")
	if err != nil {
		ll.Error("encode error", l.Error(err))
	}

	for k, v := range f.CustomFields {
		data[k] = v
	}
	return map[string]interface{}{
		"fields": data,
	}, nil
}

type IssueType struct {
	ID string `json:"id"`
}

func (r *IssueType) String() string { return jsonx.MustMarshalToString(r) }

type Project struct {
	Key string `json:"key"`
}

func (r *Project) String() string { return jsonx.MustMarshalToString(r) }

type Description struct {
	Type               string               `json:"type"`
	Version            int                  `json:"version"`
	DescriptionContent []DescriptionContent `json:"content"`
}

type DescriptionContent struct {
	Type    string    `json:"type"`
	Content []Content `json:"content"`
}

type Content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func URL(baseUrl, path string) string {
	return baseUrl + path
}
