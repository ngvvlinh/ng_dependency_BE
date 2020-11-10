package fabo

import (
	"strings"
	"time"

	"o.o/api/fabo/fbmessaging/fb_comment_action"
	"o.o/api/fabo/fbmessaging/fb_customer_conversation_type"
	"o.o/api/top/external/types"
	"o.o/api/top/types/common"
	"o.o/api/top/types/etc/connection_type"
	"o.o/api/top/types/etc/customer_type"
	"o.o/api/top/types/etc/gender"
	"o.o/api/top/types/etc/status3"
	"o.o/capi/dot"
	"o.o/capi/filter"
	"o.o/common/jsonx"
	"o.o/common/xerrors"
)

type CustomerWithFbUserAvatars struct {
	ExternalId   dot.NullString `json:"external_id"`
	ExternalCode dot.NullString `json:"external_code"`

	Id        dot.ID                         `json:"id"`
	ShopId    dot.ID                         `json:"shop_id"`
	FullName  dot.NullString                 `json:"full_name"`
	Code      dot.NullString                 `json:"code"`
	Note      dot.NullString                 `json:"note"`
	Phone     dot.NullString                 `json:"phone"`
	Email     dot.NullString                 `json:"email"`
	Gender    gender.NullGender              `json:"gender"`
	Type      customer_type.NullCustomerType `json:"type"`
	Birthday  dot.NullString                 `json:"birthday"`
	CreatedAt dot.Time                       `json:"created_at"`
	UpdatedAt dot.Time                       `json:"updated_at"`
	Status    status3.NullStatus             `json:"status"`
	Deleted   bool                           `json:"deleted"`
	FbUsers   []*FbUser                      `json:"fb_users"`
}

type ListCustomersWithFbUsersResponse struct {
	Customers []*CustomerWithFbUserAvatars `json:"customers"`
	Paging    *common.PageInfo             `json:"paging"`
}

func (m *ListCustomersWithFbUsersResponse) String() string {
	return jsonx.MustMarshalToString(m)
}

type ListCustomersWithFbUsersRequest struct {
	Filters []*common.Filter `json:"filters"`
	Paging  *common.Paging   `json:"paging"`
	GetAll  dot.NullBool     `json:"get_all"`
}

func (m *ListCustomersWithFbUsersRequest) String() string {
	return jsonx.MustMarshalToString(m)
}

type CreateFbUserCustomerRequest struct {
	ExternalID string `json:"external_id"`
	CustomerID dot.ID `json:"customer_id"`
}

func (m *CreateFbUserCustomerRequest) String() string {
	return jsonx.MustMarshalToString(m)
}

type FbUserWithCustomer struct {
	ExternalID   string              `json:"external_id"`
	ExternalInfo *ExternalFbUserInfo `json:"external_info"`
	Status       status3.Status      `json:"status"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
	CustomerID   dot.ID              `json:"customer_id"`
	Customer     *types.Customer     `json:"customer"`
	TagIDS       []dot.ID            `json:"tag_ids"`
}

func (m *FbUserWithCustomer) String() string {
	return jsonx.MustMarshalToString(m)
}

type GetFbUserRequest struct {
	ExternalID string `json:"external_id"`
}

func (m *GetFbUserRequest) String() string {
	return jsonx.MustMarshalToString(m)
}

type GetFbUserResponse struct {
	FbUser   *FbUser         `json:"fb_users"`
	Customer *types.Customer `json:"customer"`
}

func (m *GetFbUserResponse) String() string {
	return jsonx.MustMarshalToString(m)
}

type ListFbUsersRequest struct {
	CustomerID dot.NullID `json:"customer_id"`
}

func (m *ListFbUsersRequest) String() string { return jsonx.MustMarshalToString(m) }

type ListFbUsersResponse struct {
	FbUsers []*FbUserWithCustomer `json:"fb_users"`
}

func (m *ListFbUsersResponse) String() string { return jsonx.MustMarshalToString(m) }

type ConnectPagesRequest struct {
	AccessToken string `json:"access_token"`
}

func (m *ConnectPagesRequest) Censor() {
	m.AccessToken = "..."
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
	ExternalID   string              `json:"external_id"`
	ExternalInfo *ExternalFbUserInfo `json:"external_info"`
	Status       status3.Status      `json:"status"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
}

func (m *FbUserCombined) String() string { return jsonx.MustMarshalToString(m) }

type FbUser struct {
	ExternalID   string              `json:"external_id"`
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
	ShopID               dot.ID              `json:"shop_id"`
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
	ShopID               dot.ID              `json:"shop_id"`
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
	ExternalIDs    []string       `json:"ids"`
	NewExternalIDs filter.Strings `json:"external_id"`
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

type GetCustomerConversationByIDRequest struct {
	ID dot.ID `json:"id"`
}

func (m *GetCustomerConversationByIDRequest) String() string { return jsonx.MustMarshalToString(m) }

type SearchCustomerConversationRequest struct {
	Text string `json:"text"`
}

func (m *SearchCustomerConversationRequest) String() string { return jsonx.MustMarshalToString(m) }

type CustomerConversationFilter struct {
	// New
	ExternalPageID filter.Strings `json:"external_page_id"`
	ExternalUserID dot.NullString `json:"external_user_id"`

	// Old
	FbPageIDs        filter.IDs     `json:"fb_page_id"`
	FbExternalUserID dot.NullString `json:"fb_external_user_id"`

	IsRead dot.NullBool                                                 `json:"is_read"`
	Type   fb_customer_conversation_type.NullFbCustomerConversationType `json:"type"`
}

type FbCustomerConversation struct {
	ID                         dot.ID                 `json:"id"`
	ExternalPageID             string                 `json:"external_page_id"`
	ExternalID                 string                 `json:"external_id"`
	ExternalUserTags           []dot.ID               `json:"external_user_tags"`
	ExternalUserID             string                 `json:"external_user_id"`
	ExternalUserName           string                 `json:"external_user_name"`
	ExternalFrom               *FbObjectFrom          `json:"external_from"`
	IsRead                     bool                   `json:"is_read"`
	Type                       string                 `json:"type"`
	ExternalPostAttachments    []*PostAttachment      `json:"external_post_attachments"`
	ExternalCommentAttachment  *CommentAttachment     `json:"external_comment_attachment"`
	ExternalMessageAttachments []*FbMessageAttachment `json:"external_message_attachments"`
	ExternalUserPictureURL     string                 `json:"external_user_picture_url"`
	LastMessage                string                 `json:"last_message"`
	LastMessageAt              time.Time              `json:"last_message_at"`
	LastCustomerMessageAt      time.Time              `json:"last_customer_message_at"`
	CreatedAt                  time.Time              `json:"created_at"`
	UpdatedAt                  time.Time              `json:"updated_at"`
	Customer                   *types.Customer        `json:"customer"`
}

type PostAttachment struct {
	MediaType      string               `json:"media_type"`
	Type           string               `json:"type"`
	Media          *MediaPostAttachment `json:"media"`
	SubAttachments []*SubAttachment     `json:"sub_attachments"`
}

type MediaPostAttachment struct {
	Image *ImageMediaPostAttachment `json:"image"`
}

type ImageMediaPostAttachment struct {
	Height int    `json:"height"`
	Width  int    `json:"width"`
	Src    string `json:"src"`
}

type SubAttachment struct {
	Media  *MediaDataSubAttachment  `json:"media"`
	Target *TargetDataSubAttachment `json:"target"`
	Type   string                   `json:"type"`
	URL    string                   `json:"url"`
}

type MediaDataSubAttachment struct {
	Height int    `json:"height"`
	Width  int    `json:"width"`
	Src    string `json:"src"`
}

type TargetDataSubAttachment struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type FbCustomerConversationsResponse struct {
	CustomerConversations []*FbCustomerConversation `json:"fb_customer_conversations"`
	Paging                *common.CursorPageInfo    `json:"paging"`
}

func (m *FbCustomerConversationsResponse) String() string { return jsonx.MustMarshalToString(m) }

type SearchFbCustomerConversationsResponse struct {
	CustomerConversations []*FbCustomerConversation `json:"fb_customer_conversations"`
}

func (m *SearchFbCustomerConversationsResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetCustomerConversationByIDResponse struct {
	Conversation *FbCustomerConversation `json:"conversation"`
}

func (m *GetCustomerConversationByIDResponse) String() string { return jsonx.MustMarshalToString(m) }

type ListMessagesRequest struct {
	Paging *common.CursorPaging `json:"paging"`
	Filter *MessageFilter       `json:"filter"`
}

func (m *ListMessagesRequest) String() string { return jsonx.MustMarshalToString(m) }

type MessageFilter struct {
	FbExternalConversationIDs filter.Strings `json:"fb_external_conversation_ids"`
	ExternalConversationID    filter.Strings `json:"external_conversation_id"`
}

type FbMessagesResponse struct {
	FbMessages []*FbExternalMessage   `json:"fb_messages"`
	Paging     *common.CursorPageInfo `json:"paging"`
}

func (m *FbMessagesResponse) String() string { return jsonx.MustMarshalToString(m) }

type FbExternalMessage struct {
	ID                     dot.ID                 `json:"id"`
	ExternalConversationID string                 `json:"external_conversation_id"`
	ExternalPageID         string                 `json:"external_page_id"`
	ExternalID             string                 `json:"external_id"`
	ExternalMessage        string                 `json:"external_message"`
	ExternalSticker        string                 `json:"external_sticker"`
	ExternalTo             []*FbObjectTo          `json:"external_to"`
	ExternalFrom           *FbObjectFrom          `json:"external_from"`
	ExternalAttachments    []*FbMessageAttachment `json:"external_attachments"`
	ExternalCreatedTime    time.Time              `json:"external_created_time"`
	CreatedAt              time.Time              `json:"created_at"`
	UpdatedAt              time.Time              `json:"updated_at"`
}

func (m *FbExternalMessage) String() string { return jsonx.MustMarshalToString(m) }

type FbObjectTo struct {
	ID                     string `json:"id"`
	Name                   string `json:"name"`
	Email                  string `json:"email"`
	ExternalUserPictureURL string `json:"external_user_picture_url"`
}

type FbObjectFrom struct {
	ID                     string `json:"id"`
	Name                   string `json:"name"`
	Email                  string `json:"email"`
	ExternalUserPictureURL string `json:"external_user_picture_url"`
}
type FbMessageAttachment struct {
	ID        string                            `json:"id"`
	ImageData *FbMessageAttachmentImageData     `json:"image_data"`
	MimeType  string                            `json:"mime_type"`
	Name      string                            `json:"name"`
	Size      int                               `json:"size"`
	VideoData *FbMessageDataAttachmentVideoData `json:"video_data"`
	FileURL   string                            `json:"file_url"`
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

type FbMessageDataAttachmentVideoData struct {
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	Length     int    `json:"length"`
	VideoType  int    `json:"video_type"`
	URL        string `json:"url"`
	PreviewURL string `json:"preview_url"`
	Rotation   int    `json:"rotation"`
}

type ListCommentsByExternalPostIDRequest struct {
	Filter *CommentFilter       `json:"filter"`
	Paging *common.CursorPaging `json:"paging"`
}

func (m *ListCommentsByExternalPostIDRequest) String() string { return jsonx.MustMarshalToString(m) }

type CommentFilter struct {
	ExternalPostID string `json:"external_post_id"`
	ExternalUserID string `json:"external_user_id"`
}

type ListCommentsByExternalPostIDResponse struct {
	FbPost                          *FbExternalPost     `json:"fb_post"`
	LatestCustomerFbExternalComment *FbExternalComment  `json:"latest_customer_fb_external_comment"`
	FbComments                      *FbCommentsResponse `json:"fb_comments"`
}

func (m *ListCommentsByExternalPostIDResponse) String() string { return jsonx.MustMarshalToString(m) }

type FbCommentsResponse struct {
	FbComments []*FbExternalComment   `json:"data"`
	Paging     *common.CursorPageInfo `json:"paging"`
}

type FbExternalPost struct {
	ID                  dot.ID            `json:"id"`
	ExternalPageID      string            `json:"external_page_id"`
	ExternalID          string            `json:"external_id"`
	ExternalParentID    string            `json:"external_parent_id"`
	ExternalFrom        *FbObjectFrom     `json:"external_from"`
	ExternalPicture     string            `json:"external_picture"`
	ExternalIcon        string            `json:"external_icon"`
	ExternalMessage     string            `json:"external_message"`
	ExternalAttachments []*PostAttachment `json:"external_attachments"`
	ExternalCreatedTime time.Time         `json:"external_created_time"`
	CreatedAt           time.Time         `json:"created_at"`
	UpdatedAt           time.Time         `json:"updated_at"`

	ExternalParent *FbExternalPost `json:"external_parent"`
}

type FbExternalComment struct {
	ID                   dot.ID             `json:"id"`
	ExternalPostID       string             `json:"external_post_id"`
	ExternalPageID       string             `json:"external_page_id"`
	ExternalID           string             `json:"external_id"`
	ExternalUserID       string             `json:"external_user_id"`
	ExternalParentID     string             `json:"external_parent_id"`
	ExternalParentUserID string             `json:"external_parent_user_id"`
	ExternalMessage      string             `json:"external_message"`
	ExternalCommentCount int                `json:"external_comment_count"`
	ExternalParent       *FbExternalComment `json:"external_parent"`
	ExternalFrom         *FbObjectFrom      `json:"external_from"`
	ExternalAttachment   *CommentAttachment `json:"external_attachment"`
	ExternalCreatedTime  time.Time          `json:"external_created_time"`
	IsLiked              bool               `json:"is_liked"`
	IsHidden             bool               `json:"is_hidden"`
	IsPrivateReplied     bool               `json:"is_private_replied"`
	CreatedAt            time.Time          `json:"created_at"`
	UpdatedAt            time.Time          `json:"updated_at"`
}

func (m *FbExternalComment) String() string { return jsonx.MustMarshalToString(m) }

type FbObjectParent struct {
	CreatedTime time.Time     `json:"created_time"`
	From        *FbObjectFrom `json:"from"`
	Message     string        `json:"message"`
	ID          string        `json:"id"`
}

type CommentAttachment struct {
	Media  *ImageMediaDataSubAttachment `json:"media"`
	Target *TargetDataSubAttachment     `json:"target"`
	Title  string                       `json:"title"`
	Type   string                       `json:"type"`
	URL    string                       `json:"url"`
}

type ImageMediaDataSubAttachment struct {
	Image *MediaDataSubAttachment `json:"image"`
}

type UpdateReadStatusRequest struct {
	CustomerConversationID dot.ID `json:"customer_conversation_id"`
	Read                   bool   `json:"read"`
}

func (m *UpdateReadStatusRequest) String() string { return jsonx.MustMarshalToString(m) }

type SendMessageRequest struct {
	ExternalPageID         string                     `json:"external_page_id"`
	ExternalConversationID string                     `json:"external_conversation_id"`
	Message                *MessageSendMessageRequest `json:"message"`
}

func (m *SendMessageRequest) String() string { return jsonx.MustMarshalToString(m) }

type MessageSendMessageRequest struct {
	// Type: text or image
	Type string `json:"type"`
	Text string `json:"text"`
	URL  string `json:"url"`
}

type SendCommentRequest struct {
	ExternalUserID string `json:"external_user_id"`
	ExternalPageID string `json:"external_page_id"`
	ExternalID     string `json:"external_id"` // post_id, comment_id -- deprecated
	ExternalPostID string `json:"external_post_id"`
	Message        string `json:"message"`
	AttachmentURL  string `json:"attachment_url"`
}

func (m *SendCommentRequest) Validate() error {
	if m.ExternalPageID == "" {
		return xerrors.Errorf(xerrors.FailedPrecondition, nil, "missing external_page_id")
	}
	// TODO(nakhoa): uncomment following section
	//if m.ExternalUserID == "" {
	//	return xerrors.Errorf(xerrors.FailedPrecondition, nil, "missing external_user_id")
	//}
	if m.ExternalPostID == "" {
		return xerrors.Errorf(xerrors.FailedPrecondition, nil, "missing external_post_id")
	}

	m.Message = strings.TrimSpace(m.Message)
	if m.Message == "" && m.AttachmentURL == "" {
		return xerrors.Errorf(xerrors.FailedPrecondition, nil, "missing message content")
	}
	return nil
}

func (m *SendCommentRequest) String() string { return jsonx.MustMarshalToString(m) }

type LikeOrUnLikeCommentRequest struct {
	ExternalPageID    string                            `json:"external_page_id"`
	ExternalCommentID string                            `json:"external_comment_id"`
	Action            fb_comment_action.FbCommentAction `json:"action"`
}

func (m *LikeOrUnLikeCommentRequest) String() string {
	return jsonx.MustMarshalToString(m)
}

type HideOrUnHideCommentRequest struct {
	ExternalPageID    string                            `json:"external_page_id"`
	ExternalCommentID string                            `json:"external_comment_id"`
	Action            fb_comment_action.FbCommentAction `json:"action"`
}

func (m *HideOrUnHideCommentRequest) String() string {
	return jsonx.MustMarshalToString(m)
}

type SendPrivateReplyRequest struct {
	ExternalPageID    string `json:"external_page_id"`
	ExternalCommentID string `json:"external_comment_id"`
	Message           string `json:"message"`
}

func (m *SendPrivateReplyRequest) String() string {
	return jsonx.MustMarshalToString(m)
}

type CreatePostRequest struct {
	ExternalPageID string `json:"external_page_id"`
	Message        string `json:"message"`
}

func (m *CreatePostRequest) String() string {
	return jsonx.MustMarshalToString(m)
}

type CreatePostResponse struct {
	ExternalPostID string `json:"external_post_id"`
	ExternalURL    string `json:"external_url"`
}

func (m *CreatePostResponse) String() string {
	return jsonx.MustMarshalToString(m)
}

type CheckPagePermissionsRequest struct {
	ExternalPageIDS []string `json:"external_page_ids"`
}

func (c *CheckPagePermissionsRequest) String() string {
	return jsonx.MustMarshalToString(c)
}

type CheckPagePermissionsResponse struct {
	PageMissingRoles map[string][]string `json:"page_missing_roles"`
}

func (c *CheckPagePermissionsResponse) String() string {
	return jsonx.MustMarshalToString(c)
}

type CreateFbShopTagRequest struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

func (c *CreateFbShopTagRequest) Validate() error {
	if c.Name == "" {
		return xerrors.Error(xerrors.InvalidArgument, "missing name", nil)
	}
	if c.Color == "" {
		return xerrors.Error(xerrors.InvalidArgument, "missing color", nil)
	}

	return nil
}

func (c *CreateFbShopTagRequest) String() string {
	return jsonx.MustMarshalToString(c)
}

type FbShopUserTag struct {
	ID     dot.ID `json:"id"`
	Name   string `json:"name"`
	Color  string `json:"color"`
	ShopID dot.ID `json:"shop_id"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (c *FbShopUserTag) String() string {
	return jsonx.MustMarshalToString(c)
}

type ListFbShopTagResponse struct {
	ShopTags []*FbShopUserTag `json:"result"`
}

func (c *ListFbShopTagResponse) String() string {
	return jsonx.MustMarshalToString(c)
}

type CreateShopTagResponse struct {
	ID     dot.ID `json:"id"`
	Name   string `json:"name"`
	Color  string `json:"color"`
	ShopID dot.ID `json:"shop_id"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (c *CreateShopTagResponse) String() string {
	return jsonx.MustMarshalToString(c)
}

type DeleteFbShopTagRequest struct {
	ID dot.ID `json:"id"`
}

func (c *DeleteFbShopTagRequest) Validate() error {
	if c.ID == 0 {
		return xerrors.Error(xerrors.InvalidArgument, "missing id", nil)
	}
	return nil
}

func (c *DeleteFbShopTagRequest) String() string {
	return jsonx.MustMarshalToString(c)
}

type UpdateFbShopTagRequest struct {
	ID    dot.ID `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

func (c *UpdateFbShopTagRequest) Validate() error {
	if c.ID == 0 {
		return xerrors.Error(xerrors.InvalidArgument, "missing id", nil)
	}
	if c.Name == "" && c.Color == "" {
		return xerrors.Error(xerrors.InvalidArgument, "missing name or color", nil)
	}
	return nil
}

func (c *UpdateFbShopTagRequest) String() string {
	return jsonx.MustMarshalToString(c)
}

type UpdateUserTagResponse struct {
	TagIDs []dot.ID `json:"tag_ids"`
}

func (c *UpdateUserTagResponse) String() string {
	return jsonx.MustMarshalToString(c)
}

type UpdateUserTagsRequest struct {
	TagIDs           []dot.ID `json:"tag_ids"`
	FbExternalUserID dot.ID   `json:"fb_external_user_id"`
}

func (c *UpdateUserTagsRequest) Validate() error {
	if c.FbExternalUserID == 0 {
		return xerrors.Error(xerrors.InvalidArgument, "missing fb_external_user_id", nil)
	}
	return nil
}

func (c *UpdateUserTagsRequest) String() string {
	return jsonx.MustMarshalToString(c)
}

type CustomerReturnRateRequest struct {
	Phone string `json:"phone"`
}

func (r *CustomerReturnRateRequest) String() string {
	return jsonx.MustMarshalToString(r)
}

type CustomerReturnRateResponse struct {
	CustomerReturnRateExtendeds []*CustomerReturnRateExtended `json:"customer_return_rates"`
}

func (r *CustomerReturnRateResponse) String() string {
	return jsonx.MustMarshalToString(r)
}

type CustomerReturnRateExtended struct {
	ConnectionID     dot.ID                           `json:"connection_id"`
	ConnectionName   string                           `json:"connection_name"`
	ConnectionMethod connection_type.ConnectionMethod `json:"connection_method"`

	CustomerReturnRate *CustomerReturnRate `json:"customer_return_rate"`
}

func (r *CustomerReturnRateExtended) String() string {
	return jsonx.MustMarshalToString(r)
}

type CustomerReturnRate struct {
	Level     string  `json:"level"`
	LevelCode string  `json:"level_code"`
	Rate      float64 `json:"rate"`
}

func (r *CustomerReturnRate) String() string {
	return jsonx.MustMarshalToString(r)
}

type MessageTemplateVariable struct {
	Code  string `json:"code"`
	Label string `json:"label"`
}

type MessageTemplateVariableResponse struct {
	Variables []*MessageTemplateVariable `json:"variables"`
}

func (r *MessageTemplateVariableResponse) String() string {
	return jsonx.MustMarshalToString(r)
}

type MessageTemplate struct {
	ID        dot.ID    `json:"id"`
	Template  string    `json:"template"`
	ShortCode string    `json:"short_code"`
	ShopID    dot.ID    `json:"shop_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (r *MessageTemplate) String() string {
	return jsonx.MustMarshalToString(r)
}

type MessageTemplateResponse struct {
	Templates []*MessageTemplate `json:"templates"`
}

func (r *MessageTemplateResponse) String() string {
	return jsonx.MustMarshalToString(r)
}

type CreateMessageTemplateRequest struct {
	Template  string `json:"template"`
	ShortCode string `json:"short_code"`
}

func (r *CreateMessageTemplateRequest) String() string {
	return jsonx.MustMarshalToString(r)
}

type UpdateMessageTemplateRequest struct {
	ID        dot.ID         `json:"id"`
	Template  dot.NullString `json:"template"`
	ShortCode dot.NullString `json:"short_code"`
}

func (r *UpdateMessageTemplateRequest) String() string {
	return jsonx.MustMarshalToString(r)
}

type DeleteMessageTemplateRequest struct {
	ID dot.ID `json:"id"`
}

func (r *DeleteMessageTemplateRequest) String() string {
	return jsonx.MustMarshalToString(r)
}
