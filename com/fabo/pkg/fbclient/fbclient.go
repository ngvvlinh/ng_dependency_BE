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
)

var encoder = schema.NewEncoder()

func init() {
	encoder.SetAliasTag("url")
}

type RequestMethod string

const (
	GET  RequestMethod = "GET"
	POST RequestMethod = "POST"
)

type AppConfig struct {
	ID          string `yaml:"id"`
	Secret      string `yaml:"secret"`
	AccessToken string `yaml:"access_token"`
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
	if err := f.sendGetRequest(path, params, &tok); err != nil {
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
	if err := f.sendGetRequest(path, params, &me); err != nil {
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
	if err := f.sendGetRequest(path, params, &accounts); err != nil {
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
	if err := f.sendGetRequest(path, params, &tok); err != nil {
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
	if err := f.sendGetRequest(path, params, &tok); err != nil {
		return nil, err
	}

	return &tok, nil
}

func (f *FbClient) CallAPIListFeeds(accessToken, pageID string, pagination *model.FacebookPagingRequest) (*model.PublishedPostsResponse, error) {
	params := &ListFeedsParams{
		AccessToken: accessToken,
		Fields:      "id,created_time,from,full_picture,icon,is_expired,is_hidden,is_popular,is_published,message,story,permalink_url,shares,status_type,updated_time,picture,attachments{media_type,media,type,subattachments}",
		DateFormat:  UnixDateFormat,
	}

	if pagination != nil {
		pagination.ApplyQueryParams(true, DefaultLimitGetPosts, params)
	}

	path := "/me/feed"
	var publishedPostsResponse model.PublishedPostsResponse
	if err := f.sendGetRequest(path, params, &publishedPostsResponse); err != nil {
		return nil, err
	}

	return &publishedPostsResponse, nil
}

func (f *FbClient) CallAPIListPublishedPosts(accessToken, pageID string, pagination *model.FacebookPagingRequest) (*model.PublishedPostsResponse, error) {
	params := &ListPublishedPostsParams{
		AccessToken: accessToken,
		Fields:      "id,created_time,from,full_picture,icon,is_expired,is_hidden,is_popular,is_published,message,story,permalink_url,shares,status_type,updated_time,picture,attachments{media_type,media,type,subattachments}",
		DateFormat:  UnixDateFormat,
	}

	if pagination != nil {
		pagination.ApplyQueryParams(true, DefaultLimitGetPosts, params)
	}

	path := "/me/published_posts"
	var publishedPostsResponse model.PublishedPostsResponse
	if err := f.sendGetRequest(path, params, &publishedPostsResponse); err != nil {
		return nil, err
	}

	return &publishedPostsResponse, nil
}

func (f *FbClient) CallAPIGetPost(accessToken, postID string) (*model.Post, error) {
	params := &GetPostParams{
		AccessToken: accessToken,
		Fields:      "id,created_time,from,full_picture,icon,is_expired,is_hidden,is_popular,is_published,message,permalink_url,shares,status_type,updated_time,picture,attachments{media_type,type,subattachments}",
		DateFormat:  UnixDateFormat,
	}

	path := fmt.Sprintf("/%s", postID)
	var post model.Post
	if err := f.sendGetRequest(path, params, &post); err != nil {
		return nil, err
	}

	return &post, nil
}

func (f *FbClient) CallAPIListComments(accessToken, postID string, pagination *model.FacebookPagingRequest) (*model.CommentsResponse, error) {
	limit := DefaultLimitGetComments
	if pagination != nil && pagination.Limit.Valid {
		limit = pagination.Limit.Int
	}

	params := &ListCommentsParams{
		AccessToken: accessToken,
		Fields:      fmt.Sprintf("comments.filter(stream).limit(%d){message,attachment,id,created_time,comment_count,parent,from{id,name,email,first_name,last_name,picture},is_hidden}", limit),
		DateFormat:  UnixDateFormat,
	}

	if pagination != nil {
		pagination.ApplyQueryParams(false, DefaultLimitGetPosts, params)
	}

	path := fmt.Sprintf("/%s", postID)
	var commentsResponse model.CommentsResponse
	if err := f.sendGetRequest(path, params, &commentsResponse); err != nil {
		return nil, err
	}

	return &commentsResponse, nil
}

func (f *FbClient) CallAPIListCommentsByPostIDs(accessToken string, postIDs []string) (*model.CommentsByPostIDsResponse, error) {
	if len(postIDs) == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "postIDs must no be empty")
	}

	params := &ListCommentByPostIDsParams{
		AccessToken: accessToken,
		IDs:         strings.Join(postIDs, ","),
		Filter:      "stream",
		Limit:       fmt.Sprintf("%d", DefaultLimitGetComments),
		Fields:      "message,id,from{id,name,email,first_name,last_name,picture},attachment,comment_count,parent,created_time",
		DateFormat:  UnixDateFormat,
	}

	path := "/comments"
	var commentsByPostIDsResponse model.CommentsByPostIDsResponse
	if err := f.sendGetRequest(path, params, &commentsByPostIDsResponse); err != nil {
		return nil, err
	}

	return &commentsByPostIDsResponse, nil
}

func (f *FbClient) CallAPIListConversations(accessToken, pageID string, pagination *model.FacebookPagingRequest) (*model.ConversationsResponse, error) {
	defaultPaging := DefaultLimitGetConversations
	if pagination != nil && pagination.Limit.Valid {
		defaultPaging = pagination.Limit.Int
	}

	params := &ListConversationsParams{
		AccessToken: accessToken,
		Fields:      fmt.Sprintf("conversations.limit(%d){id,message_count,updated_time,link,senders}", defaultPaging),
		DateFormat:  UnixDateFormat,
	}

	if pagination != nil {
		pagination.ApplyQueryParams(false, defaultPaging, params)
	}

	path := fmt.Sprintf("/%s", pageID)
	var conversationsResponse model.ConversationsResponse
	if err := f.sendGetRequest(path, params, &conversationsResponse); err != nil {
		return nil, err
	}

	return &conversationsResponse, nil
}

func (f *FbClient) CallAPIGetConversationByUserID(accessToken, pageID, userID string) (*model.Conversations, error) {
	if userID == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "user_id must not be null")
	}

	params := &GetConversationByUserIDParams{
		AccessToken: accessToken,
		Fields:      "id,message_count,updated_time,link,senders",
		DateFormat:  UnixDateFormat,
		UserID:      userID,
	}
	path := fmt.Sprintf("/%s/conversations", pageID)
	var conversations model.Conversations
	if err := f.sendGetRequest(path, params, &conversations); err != nil {
		return nil, err
	}

	return &conversations, nil
}

func (f *FbClient) CallAPIListMessages(accessToken, conversationID string, pagination *model.FacebookPagingRequest) (*model.MessagesResponse, error) {
	defaultPaging := DefaultLimitGetMessages
	if pagination != nil && pagination.Limit.Valid {
		defaultPaging = pagination.Limit.Int
	}

	params := &ListMessagesParams{
		AccessToken: accessToken,
		Fields:      fmt.Sprintf("messages.limit(%d){id,from,to,message,sticker,created_time,attachments{id,image_data,mime_type,name,size,video_data,file_url}}", defaultPaging),
		DateFormat:  UnixDateFormat,
	}

	if pagination != nil {
		pagination.ApplyQueryParams(false, DefaultLimitGetMessages, params)
	}

	path := fmt.Sprintf("/%s", conversationID)
	var messagesResponse model.MessagesResponse
	if err := f.sendGetRequest(path, params, &messagesResponse); err != nil {
		return nil, err
	}

	return &messagesResponse, nil
}

func (f *FbClient) CallAPIGetMessage(accessToken, messageID string) (*model.MessageData, error) {
	params := &GetMessageParams{
		AccessToken: accessToken,
		Fields:      fmt.Sprintf("id,from,to,message,sticker,created_time,attachments{id,image_data,mime_type,name,size,video_data,file_url}"),
		DateFormat:  UnixDateFormat,
	}

	path := fmt.Sprintf("/%s", messageID)
	var message model.MessageData
	if err := f.sendGetRequest(path, params, &message); err != nil {
		return nil, err
	}

	return &message, nil
}

func (f *FbClient) CallAPICreateSubscribedApps(accessToken string, fields []string) (*model.SubscribedAppResponse, error) {
	params := &CreateSubscribedAppsParams{
		AccessToken:      accessToken,
		SubscribedFields: strings.Join(fields, ","),
	}

	path := "/me/subscribed_apps"
	var subscribedAddResponse model.SubscribedAppResponse
	if err := f.sendPostRequest(path, params, &subscribedAddResponse); err != nil {
		return nil, err
	}

	return &subscribedAddResponse, nil
}

func (f *FbClient) CallAPISendMessage(accessToken string, sendMessageRequest *model.SendMessageRequest) (*model.SendMessageResponse, error) {
	recipient, err := json.Marshal(sendMessageRequest.Recipient)
	if err != nil {
		return nil, err
	}
	message, err := json.Marshal(sendMessageRequest.Message)
	if err != nil {
		return nil, err
	}

	params := &SendMessageParams{
		AccessToken: accessToken,
		Recipient:   string(recipient),
		Message:     string(message),
	}

	path := "/me/messages"
	var sendMessageResponse model.SendMessageResponse
	if err := f.sendPostRequest(path, params, &sendMessageResponse); err != nil {
		return nil, err
	}

	return &sendMessageResponse, nil
}

func (f *FbClient) CallAPISendComment(accessToken string, sendCommentRequest *model.SendCommentRequest) (*model.SendCommentResponse, error) {
	params := &SendCommentParams{
		AccessToken:   accessToken,
		Message:       sendCommentRequest.Message,
		AttachmentURL: sendCommentRequest.AttachmentURL,
	}

	path := fmt.Sprintf("/%s/comments", sendCommentRequest.ID)
	var sendCommentResponse model.SendCommentResponse
	if err := f.sendPostRequest(path, params, &sendCommentResponse); err != nil {
		return nil, err
	}

	return &sendCommentResponse, nil
}

func (f *FbClient) CallAPICreatePost(accessToken string, pageID string, request *model.CreatePostRequest) (*model.CreatePostResponse, error) {
	params := &CreatePostParams{
		AccessToken: accessToken,
		Message:     request.Message,
	}

	path := fmt.Sprintf("/%s/feed", pageID)
	var response *model.CreatePostResponse
	if err := f.sendPostRequest(path, params, &response); err != nil {
		return nil, err
	}

	return response, nil
}

func (f *FbClient) CallAPICommentByID(accessToken, commentID string) (*model.Comment, error) {
	params := &GetCommentByIDParams{
		AccessToken: accessToken,
		Fields:      "message,attachment,id,created_time,comment_count,parent,from{id,name,email,first_name,last_name,picture},is_hidden",
		DateFormat:  UnixDateFormat,
	}

	path := fmt.Sprintf("/%s", commentID)
	var comment model.Comment
	if err := f.sendGetRequest(path, params, &comment); err != nil {
		return nil, err
	}

	return &comment, nil
}

func (f *FbClient) CallAPIGetProfileByPSID(accessToken, PSID string) (*model.Profile, error) {
	params := &GetProfileByPISDParams{
		AccessToken: accessToken,
		Fields:      "id,name,first_name,last_name,profile_pic",
	}

	path := fmt.Sprintf("/%s", PSID)
	var profile model.Profile
	if err := f.sendGetRequest(path, params, &profile); err != nil {
		return nil, err
	}

	return &profile, nil
}

func (f *FbClient) sendGetRequest(path string, params, resp interface{}) error {
	return f.sendRequest(GET, path, params, resp)
}

func (f *FbClient) sendPostRequest(path string, params, resp interface{}) error {
	return f.sendRequest(POST, path, params, resp)
}

func (f *FbClient) sendRequest(method RequestMethod, path string, params, resp interface{}) error {
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
	default:
		return cm.Errorf(cm.Internal, nil, "unsupported method %v", method)
	}
	if err != nil {
		return err
	}

	status := res.StatusCode()

	switch {
	case status >= 200 && status < 300:
		if err := json.Unmarshal(res.Body(), resp); err != nil {
			return cm.Errorf(cm.ExternalServiceError, err, "")
		}
	case status >= 400:
		if err := f.facebookErrorService.HandleErrorFacebookAPI(res, res.Request.URL); err != nil {
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
