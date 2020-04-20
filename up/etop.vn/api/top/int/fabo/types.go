package fabo

import (
	"time"

	"etop.vn/api/top/types/common"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/capi/dot"
	"etop.vn/common/jsonx"
)

type InitSessionRequest struct {
	AccessToken string `json:"access_token"`
}

func (m *InitSessionRequest) String() string { return jsonx.MustMarshalToString(m) }

type InitSessionResponse struct {
	FbUser  *FbUserCombined   `json:"fb_user"`
	FbPages []*FbPageCombined `json:"fb_pages"`
}

func (m *InitSessionResponse) String() string { return jsonx.MustMarshalToString(m) }

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
	Status               status3.Status      `json:"status"`
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
	Status               status3.Status      `json:"status"`
	CreatedAt            time.Time           `json:"created_at"`
	UpdatedAt            time.Time           `json:"updated_at"`
}

type ExternalCategory struct {
	ID   dot.ID `json:"id"`
	Name string `json:"name"`
}

func (m *FbPage) String() string { return jsonx.MustMarshalToString(m) }

type RemoveFbPagesRequest struct {
	IDs []dot.ID `json:"ids"`
}

func (m *RemoveFbPagesRequest) String() string { return jsonx.MustMarshalToString(m) }

type ListFbPagesRequest struct {
	Paging  *common.Paging   `json:"paging"`
	Filters []*common.Filter `json:"filters"`
}

func (m *ListFbPagesRequest) String() string { return jsonx.MustMarshalToString(m) }

type FbPagesResponse struct {
	FbPages []*FbPage        `json:"fb_pages"`
	Paging  *common.PageInfo `json:"paging"`
}

func (m *FbPagesResponse) String() string { return jsonx.MustMarshalToString(m) }
