package fabo

import (
	"time"

	"etop.vn/api/top/types/common"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/capi/dot"
	"etop.vn/common/jsonx"
)

type ConnectPagesRequest struct {
	AccessToken string `json:"access_token"`
}

func (m *ConnectPagesRequest) String() string { return jsonx.MustMarshalToString(m) }

type ConnectPagesResponse struct {
	FbUser       *FbUserCombined   `json:"fb_user"`
	FbPages      []*FbPageCombined `json:"fb_pages"`
	FbErrorPages []*FbErrorPage    `json:"fb_error_pages"`
}

func (m *ConnectPagesResponse) String() string { return jsonx.MustMarshalToString(m) }

type FbErrorPage struct {
	ExternalID       string `json:"external_id"`
	ExternalName     string `json:"external_name"`
	ExternalImageURL string `json:"external_image_url"`
	Reason           string `json:"reason"`
}

func (m *FbErrorPage) String() string { return jsonx.MustMarshalToString(m) }

type FbUserCombined struct {
	ID           dot.ID              `json:"id"`
	ExternalID   string              `json:"external_id"`
	UserID       dot.ID              `json:"user_id"`
	ExternalInfo *ExternalFbUserInfo `json:"external_info"`
	Status       status3.Status      `json:"status"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
}

func (m *FbUserCombined) String() string { return jsonx.MustMarshalToString(m) }

type FbUser struct {
	ID           dot.ID              `json:"id"`
	ExternalID   string              `json:"external_id"`
	UserID       dot.ID              `json:"user_id"`
	ExternalInfo *ExternalFbUserInfo `json:"external_info"`
	Status       status3.Status      `json:"status"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
}

func (m *FbUser) String() string { return jsonx.MustMarshalToString(m) }

type ExternalFbUserInfo struct {
	Name      string `json:"name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	ShortName string `json:"short_name"`
	ImageURL  string `json:"image_url"`
}

type FbPageCombined struct {
	ID                   dot.ID              `json:"id"`
	ExternalID           string              `json:"external_id"`
	FbUserID             dot.ID              `json:"fb_user_id"`
	ShopID               dot.ID              `json:"shop_id"`
	UserID               dot.ID              `json:"user_id"`
	ExternalName         string              `json:"external_name"`
	ExternalCategory     string              `json:"external_category"`
	ExternalCategoryList []*ExternalCategory `json:"external_category_list"`
	ExternalTasks        []string            `json:"external_tasks"`
	ExternalImageURL     string              `json:"external_image_url"`
	Status               status3.Status      `json:"status"`
	ConnectionStatus     status3.Status      `json:"connection_status"`
	CreatedAt            time.Time           `json:"created_at"`
	UpdatedAt            time.Time           `json:"updated_at"`
}

func (m *FbPageCombined) String() string { return jsonx.MustMarshalToString(m) }

type FbPage struct {
	ID                   dot.ID              `json:"id"`
	ExternalID           string              `json:"external_id"`
	FbUserID             dot.ID              `json:"fb_user_id"`
	ShopID               dot.ID              `json:"shop_id"`
	UserID               dot.ID              `json:"user_id"`
	ExternalName         string              `json:"external_name"`
	ExternalCategory     string              `json:"external_category"`
	ExternalCategoryList []*ExternalCategory `json:"external_category_list"`
	ExternalTasks        []string            `json:"external_tasks"`
	ExternalImageURL     string              `json:"external_image_url"`
	Status               status3.Status      `json:"status"`
	ConnectionStatus     status3.Status      `json:"connection_status"`
	CreatedAt            time.Time           `json:"created_at"`
	UpdatedAt            time.Time           `json:"updated_at"`
}

type ExternalCategory struct {
	ID   dot.ID `json:"id"`
	Name string `json:"name"`
}

func (m *FbPage) String() string { return jsonx.MustMarshalToString(m) }

type RemovePagesRequest struct {
	IDs []dot.ID `json:"ids"`
}

func (m *RemovePagesRequest) String() string { return jsonx.MustMarshalToString(m) }

type ListPagesRequest struct {
	Paging  *common.Paging   `json:"paging"`
	Filters []*common.Filter `json:"filters"`
}

func (m *ListPagesRequest) String() string { return jsonx.MustMarshalToString(m) }

type ListPagesResponse struct {
	FbPages []*FbPage        `json:"fb_pages"`
	Paging  *common.PageInfo `json:"paging"`
}

func (m *ListPagesResponse) String() string { return jsonx.MustMarshalToString(m) }
