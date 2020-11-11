package fbclient

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/schema"
	"gopkg.in/resty.v1"

	"o.o/backend/com/fabo/pkg/fbclient/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/httpreq"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/metrics"
)

var encoder = schema.NewEncoder()

func init() {
	encoder.SetAliasTag("url")
}

type RequestMethod string

const (
	GET    RequestMethod = "GET"
	POST   RequestMethod = "POST"
	DELETE RequestMethod = "DELETE"
)

type AppConfig struct {
	ID          string `yaml:"id"`
	Secret      string `yaml:"secret"`
	AccessToken string `yaml:"access_token"`
	Source      string `yaml:"source"`
}

func (c *AppConfig) MustLoadEnv(prefix ...string) {
	p := "ET_FACEBOOK_APP"
	if len(prefix) > 0 {
		p = prefix[0]
	}
	cc.EnvMap{
		p + "_ID":           &c.ID,
		p + "_SECRET":       &c.Secret,
		p + "_ACCESS_TOKEN": &c.AccessToken,
		p + "_SOURCE":       &c.Source,
	}.MustLoad()
}

type ApiInfo struct {
	Host    string
	Version string
}

func (api ApiInfo) Url() string {
	return fmt.Sprintf("%s/%s", api.Host, api.Version)
}

type FbClient struct {
	rclient              *httpreq.Resty
	appInfo              AppConfig
	apiInfo              ApiInfo
	facebookErrorService *FacebookErrorService
}

func New(_appInfo AppConfig) *FbClient {
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	rcfg := httpreq.RestyConfig{Client: client}
	return &FbClient{
		appInfo: _appInfo,
		apiInfo: ApiInfo{
			Host:    "https://graph.facebook.com",
			Version: "v7.0",
		},
		facebookErrorService: NewFacebookErrorService(),
		rclient:              httpreq.NewResty(rcfg),
	}
}

func (f *FbClient) Ping() error {
	params := &PingParams{
		ClientID:     f.appInfo.ID,
		ClientSecret: f.appInfo.Secret,
		GrantType:    ClientCredentials,
	}

	path := "/oauth/access_token"
	var tok model.Token
	if err := f.sendGetRequest(path, "", params, &tok); err != nil {
		return err
	}

	f.appInfo.AccessToken = tok.AccessToken
	return nil
}

func (f *FbClient) CallAPIGetMe(accessToken string) (*model.Me, error) {
	params := &GetMeParams{
		AccessToken: accessToken,
		Fields:      "id,name,last_name,first_name,short_name,picture",
	}

	path := "/me"
	var me model.Me
	if err := f.sendGetRequest(path, "", params, &me); err != nil {
		return nil, err
	}

	return &me, nil
}

// TODO: add pagination
func (f *FbClient) CallAPIGetAccounts(accessToken string) (*model.AccountsResponse, error) {
	params := &GetAccountsParams{
		AccessToken: accessToken,
		Fields:      "accounts{access_token,category,category_list,name,id,tasks,description,about,fan_count,picture},permissions",
		DateFormat:  UnixDateFormat,
	}

	path := "/me"
	var accounts model.AccountsResponse
	if err := f.sendGetRequest(path, "", params, &accounts); err != nil {
		return nil, err
	}

	return &accounts, nil
}

func (f *FbClient) CallAPIGetLongLivedAccessToken(accessToken string) (*model.Token, error) {
	params := &GetLongLivedAccessTokenParams{
		GrantType:       GrantTypeFBExchangeToken,
		FBExchangeToken: accessToken,
		ClientID:        f.appInfo.ID,
		ClientSecret:    f.appInfo.Secret,
	}

	path := "/oauth/access_token"
	var tok model.Token
	if err := f.sendGetRequest(path, "", params, &tok); err != nil {
		return nil, err
	}

	return &tok, nil
}

func (f *FbClient) CallAPICheckAccessToken(accessToken string) (*model.UserToken, error) {
	params := &DebugTokenParams{
		AccessToken: f.appInfo.AccessToken,
		InputToken:  accessToken,
	}

	path := "/debug_token"
	var tok model.UserToken
	if err := f.sendGetRequest(path, "", params, &tok); err != nil {
		return nil, err
	}

	return &tok, nil
}

func (f *FbClient) CallAPIListFeeds(req *ListFeedsRequest) (*model.PublishedPostsResponse, error) {
	params := &ListFeedsParams{
		AccessToken: req.AccessToken,
		Fields:      "id,created_time,from,full_picture,icon,is_expired,is_hidden,is_popular,is_published,message,story,permalink_url,shares,status_type,updated_time,picture,attachments{media_type,media,type,subattachments}",
		DateFormat:  UnixDateFormat,
	}

	if req.Pagination != nil {
		req.Pagination.ApplyQueryParams(true, DefaultLimitGetPosts, params)
	}

	path := "/me/feed"
	var publishedPostsResponse model.PublishedPostsResponse
	if err := f.sendGetRequest(path, req.PageID, params, &publishedPostsResponse); err != nil {
		return nil, err
	}

	return &publishedPostsResponse, nil
}

func (f *FbClient) CallAPIGetPost(req *GetPostRequest) (*model.Post, error) {
	params := &GetPostParams{
		AccessToken: req.AccessToken,
		Fields:      "id,created_time,from,full_picture,icon,is_expired,is_hidden,is_popular,is_published,message,story,permalink_url,shares,status_type,updated_time,picture,attachments{target,media_type,media,type,subattachments}",
		DateFormat:  UnixDateFormat,
	}

	path := fmt.Sprintf("/%s", req.PostID)
	var post model.Post
	if err := f.sendGetRequest(path, req.PageID, params, &post); err != nil {
		return nil, err
	}

	return &post, nil
}

func (f *FbClient) CallAPIListComments(req *ListCommentsRequest) (*model.CommentsResponse, error) {
	pagination := req.Pagination

	limit := DefaultLimitGetComments
	if pagination != nil && pagination.Limit.Valid {
		limit = pagination.Limit.Int
	}

	params := &ListCommentsParams{
		AccessToken: req.AccessToken,
		Fields:      fmt.Sprintf("comments.filter(stream).limit(%d){message,attachment,id,created_time,comment_count,parent,from{id,name,email,first_name,last_name,picture},is_hidden}", limit),
		DateFormat:  UnixDateFormat,
	}

	if pagination != nil {
		pagination.ApplyQueryParams(false, DefaultLimitGetPosts, params)
	}

	path := fmt.Sprintf("/%s", req.PostID)
	var commentsResponse model.CommentsResponse
	if err := f.sendGetRequest(path, req.PageID, params, &commentsResponse); err != nil {
		return nil, err
	}

	return &commentsResponse, nil
}

func (f *FbClient) CallAPIListConversations(req *ListConversationsRequest) (*model.ConversationsResponse, error) {
	pagination := req.Pagination

	defaultPaging := DefaultLimitGetConversations
	if pagination != nil && pagination.Limit.Valid {
		defaultPaging = pagination.Limit.Int
	}

	params := &ListConversationsParams{
		AccessToken: req.AccessToken,
		Fields:      fmt.Sprintf("conversations.limit(%d){id,message_count,updated_time,link,senders}", defaultPaging),
		DateFormat:  UnixDateFormat,
	}

	if pagination != nil {
		pagination.ApplyQueryParams(false, defaultPaging, params)
	}

	path := fmt.Sprintf("/%s", req.PageID)
	var conversationsResponse model.ConversationsResponse
	if err := f.sendGetRequest(path, req.PageID, params, &conversationsResponse); err != nil {
		return nil, err
	}

	return &conversationsResponse, nil
}

func (f *FbClient) CallAPIGetConversationByUserID(req *GetConversationByUserIDRequest) (*model.Conversations, error) {
	if req.UserID == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "user_id must not be null")
	}

	params := &GetConversationByUserIDParams{
		AccessToken: req.AccessToken,
		Fields:      "id,message_count,updated_time,link,senders",
		DateFormat:  UnixDateFormat,
		UserID:      req.UserID,
	}
	path := fmt.Sprintf("/%s/conversations", req.PageID)
	var conversations model.Conversations
	if err := f.sendGetRequest(path, req.PageID, params, &conversations); err != nil {
		return nil, err
	}

	return &conversations, nil
}

func (f *FbClient) CallAPIListMessages(req *ListMessagesRequest) (*model.MessagesResponse, error) {
	pagination := req.Pagination

	defaultPaging := DefaultLimitGetMessages
	if pagination != nil && pagination.Limit.Valid {
		defaultPaging = pagination.Limit.Int
	}

	params := &ListMessagesParams{
		AccessToken: req.AccessToken,
		Fields:      fmt.Sprintf("messages.limit(%d){id,from,to,message,sticker,created_time,attachments{id,image_data,mime_type,name,size,video_data,file_url},shares{description,id,link,name}}", defaultPaging),
		DateFormat:  UnixDateFormat,
	}

	if pagination != nil {
		pagination.ApplyQueryParams(false, DefaultLimitGetMessages, params)
	}

	path := fmt.Sprintf("/%s", req.ConversationID)
	var messagesResponse model.MessagesResponse
	if err := f.sendGetRequest(path, req.PageID, params, &messagesResponse); err != nil {
		return nil, err
	}

	return &messagesResponse, nil
}

func (f *FbClient) CallAPIGetMessage(req *GetMessageRequest) (*model.MessageData, error) {
	params := &GetMessageParams{
		AccessToken: req.AccessToken,
		Fields:      fmt.Sprintf("id,from,to,message,sticker,created_time,attachments{id,image_data,mime_type,name,size,video_data,file_url},shares{description,id,link,name}"),
		DateFormat:  UnixDateFormat,
	}

	path := fmt.Sprintf("/%s", req.MessageID)
	var message model.MessageData
	if err := f.sendGetRequest(path, req.PageID, params, &message); err != nil {
		return nil, err
	}

	return &message, nil
}

func (f *FbClient) CallAPICreateSubscribedApps(req *CreateSubscribedAppsRequest) (*model.SubscribedAppResponse, error) {
	params := &CreateSubscribedAppsParams{
		AccessToken:      req.AccessToken,
		SubscribedFields: strings.Join(req.Fields, ","),
	}

	path := "/me/subscribed_apps"
	var subscribedAddResponse model.SubscribedAppResponse
	if err := f.sendPostRequest(path, req.PageID, params, &subscribedAddResponse); err != nil {
		return nil, err
	}
	return &subscribedAddResponse, nil
}

func (f *FbClient) CallAPISendMessage(req *SendMessageRequest) (*model.SendMessageResponse, error) {
	recipient, err := json.Marshal(req.SendMessageArgs.Recipient)
	if err != nil {
		return nil, err
	}
	message, err := json.Marshal(req.SendMessageArgs.Message)
	if err != nil {
		return nil, err
	}

	params := &SendMessageParams{
		AccessToken: req.AccessToken,
		Recipient:   string(recipient),
		Message:     string(message),
		Tag:         req.SendMessageArgs.Tag,
	}

	path := "/me/messages"
	var sendMessageResponse model.SendMessageResponse
	if err = f.sendPostRequest(path, req.PageID, params, &sendMessageResponse); err != nil {
		return nil, err
	}
	return &sendMessageResponse, nil
}

func (f *FbClient) CallAPISendComment(req *SendCommentRequest) (*model.SendCommentResponse, error) {
	sendCommentArgs := req.SendCommentArgs
	params := &SendCommentParams{
		AccessToken:   req.AccessToken,
		Message:       sendCommentArgs.Message,
		AttachmentURL: sendCommentArgs.AttachmentURL,
	}

	path := fmt.Sprintf("/%s/comments", sendCommentArgs.ID)
	var sendCommentResponse model.SendCommentResponse
	if err := f.sendPostRequest(path, req.PageID, params, &sendCommentResponse); err != nil {
		return nil, err
	}

	return &sendCommentResponse, nil
}

func (f *FbClient) CallAPILikeComment(req *LikeCommentRequest) (*model.CommonResponse, error) {
	params := &LikeCommentParams{
		AccessToken: req.AccessToken,
	}

	path := fmt.Sprintf("/%s/likes", req.CommentID)
	var likeCommentResponse model.CommonResponse
	if err := f.sendPostRequest(path, req.PageID, params, &likeCommentResponse); err != nil {
		return nil, err
	}

	return &likeCommentResponse, nil
}

func (f *FbClient) CallAPIUnLikeComment(req *UnLikeCommentRequest) (*model.CommonResponse, error) {
	params := &UnLikeCommentParams{
		AccessToken: req.AccessToken,
	}

	path := fmt.Sprintf("/%s/likes", req.CommentID)
	var unLikeCommentResponse model.CommonResponse
	if err := f.sendDeleteRequest(path, req.PageID, params, &unLikeCommentResponse); err != nil {
		return nil, err
	}

	return &unLikeCommentResponse, nil
}

func (f *FbClient) CallAPIHideAndUnHideComment(req *HideOrUnHideCommentRequest) (*model.CommonResponse, error) {
	params := &HideOrUnHideCommentParams{
		AccessToken: req.AccessToken,
		IsHidden:    req.IsHidden,
	}

	path := fmt.Sprintf("/%s", req.CommentID)
	var hideAndUnHideCommentResponse model.CommonResponse
	if err := f.sendPostRequest(path, req.PageID, params, &hideAndUnHideCommentResponse); err != nil {
		return nil, err
	}

	return &hideAndUnHideCommentResponse, nil
}

func (f *FbClient) CallAPICreatePost(req *CreatePostRequest) (*model.CreatePostResponse, error) {
	content := req.Content
	params := &CreatePostParams{
		AccessToken: req.AccessToken,
		Message:     content.Message,
	}

	path := fmt.Sprintf("/%s/feed", req.PageID)
	var response *model.CreatePostResponse
	if err := f.sendPostRequest(path, req.PageID, params, &response); err != nil {
		return nil, err
	}

	return response, nil
}

func (f *FbClient) CallAPICommentByID(req *GetCommentByIDRequest) (*model.Comment, error) {
	params := &GetCommentByIDParams{
		AccessToken: req.AccessToken,
		Fields:      "message,attachment,id,created_time,comment_count,parent,from{id,name,email,first_name,last_name,picture},is_hidden",
		DateFormat:  UnixDateFormat,
	}

	path := fmt.Sprintf("/%s", req.CommentID)
	var comment model.Comment
	if err := f.sendGetRequest(path, req.PageID, params, &comment); err != nil {
		return nil, err
	}

	return &comment, nil
}

func (f *FbClient) CallAPIGetProfileByPSID(req *GetProfileRequest) (*model.Profile, error) {
	if req.ProfileDefault == nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Profile mustn't be null")
	}
	params := &GetProfileByPISDParams{
		AccessToken: req.AccessToken,
		Fields:      "id,name,first_name,last_name,profile_pic",
	}

	path := fmt.Sprintf("/%s", req.PSID)
	var profile model.Profile
	if err := f.sendGetRequest(path, req.PageID, params, &profile); err == nil {
		return &profile, nil
	}

	linkImage := fmt.Sprintf("https://graph.facebook.com/%s/picture?height=200&width=200&type=normal", req.PSID)
	resp, err := http.Get(linkImage)
	if err != nil || (resp != nil && resp.StatusCode != 200) {
		req.ProfileDefault.ProfilePic = DefaultFaboImage
	} else {
		req.ProfileDefault.ProfilePic = linkImage
	}

	return req.ProfileDefault, nil
}

func (f *FbClient) sendGetRequest(path, pageID string, params, resp interface{}) error {
	return f.sendRequest(GET, path, pageID, params, resp)
}

func (f *FbClient) sendPostRequest(path, pageID string, params, resp interface{}) error {
	return f.sendRequest(POST, path, pageID, params, resp)
}

func (f *FbClient) sendDeleteRequest(path, pageID string, params, resp interface{}) error {
	return f.sendRequest(DELETE, path, pageID, params, resp)
}

func (f *FbClient) sendRequest(method RequestMethod, path, pageID string, params, resp interface{}) error {
	t0 := time.Now()

	queryString := url.Values{}
	if params != nil {
		err := encoder.Encode(params, queryString)
		if err != nil {
			return cm.Error(cm.Internal, "", err)
		}
	}

	var res *resty.Response
	var err error

	_url := f.apiInfo.Url() + path
	req := f.rclient.R().
		SetQueryString(queryString.Encode())

	switch method {
	case GET:
		res, err = req.Get(_url)
	case POST:
		res, err = req.Post(_url)
	case DELETE:
		res, err = req.Delete(_url)
	default:
		return cm.Errorf(cm.Internal, nil, "unsupported method %v", method)
	}
	if err != nil {
		return err
	}

	status := res.StatusCode()

	d := time.Now().Sub(t0)
	metrics.FaboEgressRequest(req.RawRequest.URL, status, d, f.appInfo.Source, pageID)
	switch {
	case status >= 200 && status < 300:
		if err = json.Unmarshal(res.Body(), resp); err != nil {
			return cm.Errorf(cm.ExternalServiceError, err, "")
		}
	case status >= 400:
		if err = f.facebookErrorService.HandleErrorFacebookAPI(res, res.Request.URL); err != nil {
			return err
		}
	}
	return nil
}

func GetRole(tasks []string) FacebookRole {
	var hasAdvertise, hasAnalyze, hasCreateContent, hasManage, hasModerate bool
	for _, task := range tasks {
		switch task {
		case "ADVERTISE":
			hasAdvertise = true
		case "ANALYZE":
			hasAnalyze = true
		case "CREATE_CONTENT":
			hasCreateContent = true
		case "MANAGE":
			hasManage = true
		case "MODERATE":
			hasModerate = true
		}
	}

	if isBool(hasAdvertise, hasAnalyze, hasCreateContent, hasManage, hasModerate) {
		return ADMIN
	}
	if isBool(hasAdvertise, hasAnalyze, hasCreateContent, hasModerate) {
		return EDITOR
	}
	if isBool(hasAdvertise, hasAnalyze, hasModerate) {
		return MODERATOR
	}
	if isBool(hasAdvertise, hasAnalyze) {
		return ADVERTISER
	}
	if isBool(hasAnalyze) {
		return ANALYST
	}

	return UNKNOWN
}

func isBool(A ...bool) bool {
	for _, a := range A {
		if !a {
			return false
		}
	}
	return true
}
