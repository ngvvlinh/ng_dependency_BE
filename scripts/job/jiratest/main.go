package main

import (
	"encoding/base64"

	"github.com/k0kubun/pp"
	"gopkg.in/resty.v1"
	jiraclient "o.o/backend/pkg/integration/jira/client"
	"o.o/common/l"
)

type A struct {
	*B
	CustomFields map[string]string
}

type B struct {
	Name string `json:"name"`
	ID   string `json:"id"`
	Type BType  `json:"type"`
}

type BType struct {
	TypeName string `json:"type_name"`
}

var ll = l.New()

func main() {
	myData := jiraclient.CreateIssueRequest{
		IssueDefaultFields: jiraclient.IssueDefaultFields{
			Summary: "Tuan test",
			IssueType: jiraclient.IssueType{
				ID: "10059",
			},
			Project: jiraclient.Project{
				Key: "BET",
			},
		},
		CustomFields: map[string]string{
			"customfield_10062": "Tuan Phan",
			"customfield_10070": "tuan@etop.vn",
		},
	}

	body, err := myData.ToJiraBodyRequest()
	if err != nil {
		ll.Error("ToJiraBodyRequest error", l.Error(err))
	}
	pp.Println("body :: ", body)
	token := base64.StdEncoding.EncodeToString([]byte("kimhai.ngvan@gmail.com:63aiaN3wd5OgsSW77VYjCB17"))
	req := resty.R().
		SetHeader("Authorization", "Basic "+token).
		SetBody(body)

	url := "https://e-b2b.atlassian.net/rest/api/3/issue"
	res, err := req.Post(url)

	if err != nil {
		ll.Error("create issue error", l.Error(err))
	}
	pp.Println("result :: ", string(res.Body()))
}
