package fabo

import (
	"time"

	"o.o/api/fabo/fbmessaging/fb_customer_conversation_type"
	"o.o/api/top/types/common"
	"o.o/api/top/types/etc/status3"
	"o.o/capi/dot"
	"o.o/capi/filter"
	"o.o/common/jsonx"
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
	ShopID       dot.ID              `json:"shop_id"`
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
	ShopID       dot.ID              `json:"shop_id"`
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
	ExternalPermissions  []string            `json:"external_permissions"`
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
	ExternalPermissions  []string            `json:"external_permissions"`
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

type ListCustomerConversationsRequest struct {
	Paging *common.CursorPaging        `json:"paging"`
	Filter *CustomerConversationFilter `json:"filter"`
}

func (m *ListCustomerConversationsRequest) String() string { return jsonx.MustMarshalToString(m) }

type CustomerConversationFilter struct {
	FbPageIDs        filter.IDs                                                   `json:"fb_page_ids"`
	FbExternalUserID dot.NullString                                               `json:"fb_external_user_id"`
	IsRead           dot.NullBool                                                 `json:"is_read"`
	Type             fb_customer_conversation_type.NullFbCustomerConversationType `json:"type"`
}

type FbCustomerConversation struct {
	ID                     dot.ID                                                   `json:"id"`
	FbPageID               dot.ID                                                   `json:"fb_page_id"`
	ExternalID             string                                                   `json:"external_id"`
	ExternalUserID         string                                                   `json:"external_user_id"`
	ExternalUserName       string                                                   `json:"external_user_name"`
	IsRead                 bool                                                     `json:"is_read"`
	Type                   fb_customer_conversation_type.FbCustomerConversationType `json:"type"`
	PostAttachments        []*PostAttachment                                        `json:"post_attachments"`
	ExternalUserPictureURL string                                                   `json:"external_user_picture_url"`
	LastMessage            string                                                   `json:"last_message"`
	LastMessageAt          time.Time                                                `json:"last_message_at"`
	CreatedAt              time.Time                                                `json:"created_at"`
	UpdatedAt              time.Time                                                `json:"updated_at"`
}

type PostAttachment struct {
	Media *PostAttachmentMedia `json:"media"`
	Type  string               `json:"type"`
}

type PostAttachmentMedia struct {
	Height int    `json:"height"`
	Width  int    `json:"width"`
	Src    string `json:"src"`
}

type FbCustomerConversationsResponse struct {
	CustomerConversations []*FbCustomerConversation `json:"fb_customer_conversations"`
	Paging                *common.CursorPageInfo    `json:"paging"`
}

func (m *FbCustomerConversationsResponse) String() string { return jsonx.MustMarshalToString(m) }

type ListMessagesRequest struct {
	Paging *common.CursorPaging `json:"paging"`
	Filter *MessageFilter       `json:"filter"`
}

func (m *ListMessagesRequest) String() string { return jsonx.MustMarshalToString(m) }

type MessageFilter struct {
	FbExternalConversationIDs []string `json:"fb_external_conversation_ids"`
}

type FbMessagesResponse struct {
	FbMessages []*FbExternalMessage   `json:"fb_messages"`
	Paging     *common.CursorPageInfo `json:"paging"`
}

func (m *FbMessagesResponse) String() string { return jsonx.MustMarshalToString(m) }

type FbExternalMessage struct {
	ID                     dot.ID                 `json:"id"`
	FbConversationID       dot.ID                 `json:"fb_conversation_id"`
	ExternalConversationID string                 `json:"external_conversation_id"`
	FbPageID               dot.ID                 `json:"fb_page_id"`
	ExternalID             string                 `json:"external_id"`
	ExternalMessage        string                 `json:"external_message"`
	ExternalTo             []*FbObjectTo          `json:"external_to"`
	ExternalFrom           *FbObjectFrom          `json:"external_from"`
	ExternalAttachments    []*FbMessageAttachment `json:"external_attachments"`
	ExternalCreatedTime    time.Time              `json:"external_created_time"`
	CreatedAt              time.Time              `json:"created_at"`
	UpdatedAt              time.Time              `json:"updated_at"`
}

type FbObjectTo struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type FbObjectFrom struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
type FbMessageAttachment struct {
	ID        string                        `json:"id"`
	ImageData *FbMessageAttachmentImageData `json:"image_data"`
	MimeType  string                        `json:"mime_type"`
	Name      string                        `json:"name"`
	Size      int                           `json:"size"`
}

type FbMessageAttachmentImageData struct {
	Width           int    `json:"width"`
	Height          int    `json:"height"`
	MaxWidth        int    `json:"max_width"`
	MaxHeight       int    `json:"max_height"`
	URL             string `json:"url"`
	PreviewURL      string `json:"preview_url"`
	ImageType       int    `json:"image_type"`
	RenderAsSticker bool   `json:"render_as_sticker"`
}
