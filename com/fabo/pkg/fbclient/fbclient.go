package fbclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"o.o/backend/com/fabo/pkg/fbclient/model"
	cc "o.o/backend/pkg/common/config"
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
	appInfo              AppConfig
	apiInfo              ApiInfo
	facebookErrorService *FacebookErrorService
}

func New(_appInfo AppConfig) *FbClient {
	return &FbClient{
		appInfo: _appInfo,
		apiInfo: ApiInfo{
			Host:    "https://graph.facebook.com",
			Version: "v6.0",
		},
		facebookErrorService: NewFacebookErrorService(),
	}
}

func (f *FbClient) Ping() error {
	URL, err := url.Parse(fmt.Sprintf("%s/oauth/access_token", f.apiInfo.Url()))
	if err != nil {
		return err
	}

	query, err := url.ParseQuery(URL.RawQuery)
	if err != nil {
		return err
	}

	query.Add(ClientIDKey, f.appInfo.ID)
	query.Add(ClientSecret, f.appInfo.Secret)
	query.Add(GrantType, ClientCredentials)
	URL.RawQuery = query.Encode()
	resp, err := http.Get(URL.String())
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	if err := f.facebookErrorService.HandleErrorFacebookAPI(body, URL.String()); err != nil {
		return err
	}

	var tok model.Token
	if err := json.Unmarshal(body, &tok); err != nil {
		return err
	}

	f.appInfo.AccessToken = tok.AccessToken

	return nil
}

func (f *FbClient) CallAPIGetMe(accessToken string) (*model.Me, error) {
	URL, err := url.Parse(fmt.Sprintf("%s/me", f.apiInfo.Url()))
	if err != nil {
		return nil, err
	}

	query, err := url.ParseQuery(URL.RawQuery)
	if err != nil {
		return nil, err
	}

	query.Add(AccessToken, accessToken)
	query.Add(Fields, "id,name,last_name,first_name,short_name,picture")
	URL.RawQuery = query.Encode()
	resp, err := http.Get(URL.String())
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if err := f.facebookErrorService.HandleErrorFacebookAPI(body, URL.String()); err != nil {
		return nil, err
	}

	var me model.Me
	if err := json.Unmarshal(body, &me); err != nil {
		return nil, err
	}

	return &me, nil
}

// TODO: add pagination
func (f *FbClient) CallAPIGetAccounts(accessToken string) (*model.AccountsResponse, error) {
	URL, err := url.Parse(fmt.Sprintf("%s/me", f.apiInfo.Url()))
	if err != nil {
		return nil, err
	}

	query, err := url.ParseQuery(URL.RawQuery)
	if err != nil {
		return nil, err
	}

	query.Add(AccessToken, accessToken)
	query.Add(Fields, "accounts{access_token,category,category_list,name,id,tasks,description,about,fan_count,picture},permissions")
	query.Add(DateFormat, UnixDateFormat)
	URL.RawQuery = query.Encode()
	resp, err := http.Get(URL.String())
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if err := f.facebookErrorService.HandleErrorFacebookAPI(body, URL.String()); err != nil {
		return nil, err
	}

	var accounts model.AccountsResponse

	if err := json.Unmarshal(body, &accounts); err != nil {
		return nil, err
	}
	return &accounts, nil
}

func (f *FbClient) CallAPIGetLongLivedAccessToken(accessToken string) (*model.Token, error) {
	URL, err := url.Parse(fmt.Sprintf("%s/oauth/access_token", f.apiInfo.Url()))
	if err != nil {
		return nil, err
	}

	query, err := url.ParseQuery(URL.RawQuery)
	if err != nil {
		return nil, err
	}

	query.Add(GrantType, GrantTypeFBExchangeToken)
	query.Add(FBExchangeToken, accessToken)
	query.Add(ClientIDKey, f.appInfo.ID)
	query.Add(ClientSecret, f.appInfo.Secret)

	URL.RawQuery = query.Encode()
	resp, err := http.Get(URL.String())
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if err := f.facebookErrorService.HandleErrorFacebookAPI(body, URL.String()); err != nil {
		return nil, err
	}

	var tok model.Token

	if err := json.Unmarshal(body, &tok); err != nil {
		return nil, err
	}
	return &tok, nil
}

func (f *FbClient) CallAPICheckAccessToken(accessToken string) (*model.UserToken, error) {
	URL, err := url.Parse(fmt.Sprintf("%s/debug_token", f.apiInfo.Url()))
	if err != nil {
		return nil, err
	}

	query, err := url.ParseQuery(URL.RawQuery)
	if err != nil {
		return nil, err
	}

	query.Add(AccessToken, f.appInfo.AccessToken)
	query.Add(InputToken, accessToken)

	URL.RawQuery = query.Encode()
	resp, err := http.Get(URL.String())
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if err := f.facebookErrorService.HandleErrorFacebookAPI(body, URL.String()); err != nil {
		return nil, err
	}

	var tok model.UserToken

	if err := json.Unmarshal(body, &tok); err != nil {
		return nil, err
	}

	return &tok, nil
}

func (f *FbClient) CallAPIListPublishedPosts(accessToken, pageID string, pagination *model.FacebookPagingRequest) (*model.PublishedPostsResponse, error) {
	URL, err := url.Parse(fmt.Sprintf("%s/%s/published_posts", f.apiInfo.Url(), pageID))
	if err != nil {
		return nil, err
	}

	query, err := url.ParseQuery(URL.RawQuery)
	if err != nil {
		return nil, err
	}

	query.Add(AccessToken, accessToken)
	query.Add(Fields, "id,created_time,from,full_picture,icon,is_expired,is_hidden,is_popular,is_published,message,story,permalink_url,shares,status_type,updated_time,picture,attachments{media_type,type,subattachments}")
	query.Add(DateFormat, UnixDateFormat)

	URL.RawQuery = query.Encode()
	resp, err := http.Get(pagination.AddQueryParams(URL.String(), true, DefaultLimitGetPosts))
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if err := f.facebookErrorService.HandleErrorFacebookAPI(body, URL.String()); err != nil {
		return nil, err
	}

	var publishedPostsResponse model.PublishedPostsResponse

	if err := json.Unmarshal(body, &publishedPostsResponse); err != nil {
		return nil, err
	}

	return &publishedPostsResponse, nil
}

func (f *FbClient) CallAPIListPostsByIDs(accessToken string, postIDs []string) (*model.PublishedPostsByIDsResponse, error) {
	// TODO: Ngoc add specific error
	if len(postIDs) == 0 {
		return nil, nil
	}

	URL, err := url.Parse(fmt.Sprintf("%s/", f.apiInfo.Url()))
	if err != nil {
		return nil, err
	}

	query, err := url.ParseQuery(URL.RawQuery)
	if err != nil {
		return nil, err
	}

	query.Add(AccessToken, accessToken)
	query.Add(IDs, strings.Join(postIDs, ","))
	query.Add(Fields, "id,created_time,from,icon,updated_time,picture,name")
	query.Add(DateFormat, UnixDateFormat)

	URL.RawQuery = query.Encode()
	fmt.Println(URL.String())
	resp, err := http.Get(URL.String())
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if err := f.facebookErrorService.HandleErrorFacebookAPI(body, URL.String()); err != nil {
		return nil, err
	}

	var publishedPostsByIDsResponse model.PublishedPostsByIDsResponse

	if err := json.Unmarshal(body, &publishedPostsByIDsResponse); err != nil {
		return nil, err
	}

	return &publishedPostsByIDsResponse, nil
}

func (f *FbClient) CallAPIGetPost(postID, accessToken string) (*model.Post, error) {
	URL, err := url.Parse(fmt.Sprintf("%s/%s", f.apiInfo.Url(), postID))
	if err != nil {
		return nil, err
	}

	query, err := url.ParseQuery(URL.RawQuery)
	if err != nil {
		return nil, err
	}

	query.Add(AccessToken, accessToken)
	query.Add(Fields, "id,created_time,from,full_picture,icon,is_expired,is_hidden,is_popular,is_published,message,permalink_url,shares,status_type,updated_time,picture,attachments{media_type,type,subattachments}")
	query.Add(DateFormat, UnixDateFormat)

	URL.RawQuery = query.Encode()
	resp, err := http.Get(URL.String())
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if err := f.facebookErrorService.HandleErrorFacebookAPI(body, URL.String()); err != nil {
		return nil, err
	}

	var post model.Post

	if err := json.Unmarshal(body, &post); err != nil {
		return nil, err
	}

	return &post, nil
}

func (f *FbClient) CallAPIListComments(accessToken, postID string, pagination *model.FacebookPagingRequest) (*model.CommentsResponse, error) {
	URL, err := url.Parse(fmt.Sprintf("%s/%s", f.apiInfo.Url(), postID))
	if err != nil {
		return nil, err
	}

	query, err := url.ParseQuery(URL.RawQuery)
	if err != nil {
		return nil, err
	}

	limit := DefaultLimitGetComments
	if pagination != nil && pagination.Limit.Valid {
		limit = pagination.Limit.Int
	}

	query.Add(AccessToken, accessToken)
	query.Add(Fields, fmt.Sprintf("comments.filter(stream).limit(%d){message,attachment,id,created_time,comment_count,parent,from,is_hidden}", limit))
	query.Add(DateFormat, UnixDateFormat)

	URL.RawQuery = query.Encode()
	resp, err := http.Get(pagination.AddQueryParams(URL.String(), false, DefaultLimitGetComments))
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if err := f.facebookErrorService.HandleErrorFacebookAPI(body, URL.String()); err != nil {
		return nil, err
	}

	var commentsResponse model.CommentsResponse

	if err := json.Unmarshal(body, &commentsResponse); err != nil {
		return nil, err
	}

	return &commentsResponse, nil
}

func (f *FbClient) CallAPIListCommentsByPostIDs(accessToken string, postIDs []string) (*model.CommentsByPostIDsResponse, error) {
	// TODO: Ngoc add specific error
	if len(postIDs) == 0 {
		return nil, nil
	}
	URL, err := url.Parse(fmt.Sprintf("%s/comments", f.apiInfo.Url()))
	if err != nil {
		return nil, err
	}

	query, err := url.ParseQuery(URL.RawQuery)
	if err != nil {
		return nil, err
	}

	query.Add(AccessToken, accessToken)
	query.Add(IDs, strings.Join(postIDs, ","))
	query.Add(Filter, "stream")
	query.Add(Limit, fmt.Sprintf("%d", DefaultLimitGetComments))
	query.Add(Fields, "message,id,from,attachment,comment_count,parent,created_time")
	query.Add(DateFormat, UnixDateFormat)

	URL.RawQuery = query.Encode()
	fmt.Println(URL.String())
	resp, err := http.Get(URL.String())
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if err := f.facebookErrorService.HandleErrorFacebookAPI(body, URL.String()); err != nil {
		return nil, err
	}

	var commentsByPostIDsResponse model.CommentsByPostIDsResponse

	if err := json.Unmarshal(body, &commentsByPostIDsResponse); err != nil {
		return nil, err
	}

	return &commentsByPostIDsResponse, nil
}

func (f *FbClient) CallAPIListCommentSummaries(accessToken string, postIDs []string) (map[string]*model.CommentSummary, error) {
	URL, err := url.Parse(fmt.Sprintf("%s/comments", f.apiInfo.Url()))
	if err != nil {
		return nil, err
	}

	query, err := url.ParseQuery(URL.RawQuery)
	if err != nil {
		return nil, err
	}

	query.Add(AccessToken, accessToken)
	// TODO: Ngoc add specific error
	if len(postIDs) <= 0 {
		return nil, nil
	}
	query.Add(IDs, strings.Join(postIDs, ","))
	query.Add(Limit, "0")
	query.Add(Filter, "stream")
	query.Add(Summary, "true")

	URL.RawQuery = query.Encode()
	resp, err := http.Get(URL.String())
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if err := f.facebookErrorService.HandleErrorFacebookAPI(body, URL.String()); err != nil {
		return nil, err
	}

	var commentsSummariesResponse map[string]*model.CommentSummary

	if err := json.Unmarshal(body, &commentsSummariesResponse); err != nil {
		return nil, err
	}

	return commentsSummariesResponse, nil
}

func (f *FbClient) CallAPIListConversations(accessToken, pageID, userID string, pagination *model.FacebookPagingRequest) (*model.ConversationsResponse, error) {
	URL, err := url.Parse(fmt.Sprintf("%s/%s", f.apiInfo.Url(), pageID))
	if err != nil {
		return nil, err
	}

	query, err := url.ParseQuery(URL.RawQuery)
	if err != nil {
		return nil, err
	}

	defaultPaging := DefaultLimitGetConversations
	if pagination != nil && pagination.Limit.Valid {
		defaultPaging = pagination.Limit.Int
	}

	query.Add(AccessToken, accessToken)
	query.Add(Fields, fmt.Sprintf("conversations.limit(%d){id,message_count,updated_time,link,senders}", defaultPaging))
	query.Add(DateFormat, UnixDateFormat)

	if userID != "" {
		query.Add(UserID, userID)
	}

	URL.RawQuery = query.Encode()
	resp, err := http.Get(pagination.AddQueryParams(URL.String(), false, defaultPaging))
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if err := f.facebookErrorService.HandleErrorFacebookAPI(body, URL.String()); err != nil {
		return nil, err
	}

	var conversationsResponse model.ConversationsResponse

	if err := json.Unmarshal(body, &conversationsResponse); err != nil {
		return nil, err
	}

	return &conversationsResponse, nil
}

func (f *FbClient) CallAPIListMessages(accessToken, conversationID string, pagination *model.FacebookPagingRequest) (*model.MessagesResponse, error) {
	URL, err := url.Parse(fmt.Sprintf("%s/%s", f.apiInfo.Url(), conversationID))
	if err != nil {
		return nil, err
	}

	defaultPaging := DefaultLimitGetMessages
	if pagination != nil && pagination.Limit.Valid {
		defaultPaging = pagination.Limit.Int
	}

	query, err := url.ParseQuery(URL.RawQuery)
	if err != nil {
		return nil, err
	}

	query.Add(AccessToken, accessToken)
	query.Add(Fields, fmt.Sprintf("messages.limit(%d){id,from,to,message,sticker,created_time,attachments{id,image_data,mime_type,name,size,video_data,file_url}}", defaultPaging))
	query.Add(DateFormat, UnixDateFormat)

	URL.RawQuery = query.Encode()
	resp, err := http.Get(pagination.AddQueryParams(URL.String(), false, DefaultLimitGetMessages))
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if err := f.facebookErrorService.HandleErrorFacebookAPI(body, URL.String()); err != nil {
		return nil, err
	}

	var messagesResponse model.MessagesResponse

	if err := json.Unmarshal(body, &messagesResponse); err != nil {
		return nil, err
	}

	return &messagesResponse, nil
}

func (f *FbClient) CallAPIGetMessage(accessToken, messageID string) (*model.MessageData, error) {
	URL, err := url.Parse(fmt.Sprintf("%s/%s", f.apiInfo.Url(), messageID))
	if err != nil {
		return nil, err
	}

	query, err := url.ParseQuery(URL.RawQuery)
	if err != nil {
		return nil, err
	}

	query.Add(AccessToken, accessToken)
	query.Add(Fields, fmt.Sprintf("id,from,to,message,sticker,created_time,attachments{id,image_data,mime_type,name,size,video_data,file_url}"))
	query.Add(DateFormat, UnixDateFormat)

	URL.RawQuery = query.Encode()
	resp, err := http.Get(URL.String())
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if err := f.facebookErrorService.HandleErrorFacebookAPI(body, URL.String()); err != nil {
		return nil, err
	}

	var message model.MessageData

	if err := json.Unmarshal(body, &message); err != nil {
		return nil, err
	}

	return &message, nil
}

func (f *FbClient) CallAPICreateSubscribedApps(accessToken string, fields []string) (*model.SubcribedAppResponse, error) {
	URL, err := url.Parse(fmt.Sprintf("%s/me/subscribed_apps", f.apiInfo.Url()))
	if err != nil {
		return nil, err
	}

	query, err := url.ParseQuery(URL.RawQuery)
	if err != nil {
		return nil, err
	}

	query.Add(AccessToken, accessToken)
	query.Add(SubscribedFields, strings.Join(fields, ","))
	URL.RawQuery = query.Encode()

	resp, err := http.Post(URL.String(), "", nil)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if err := f.facebookErrorService.HandleErrorFacebookAPI(body, URL.String()); err != nil {
		return nil, err
	}

	var subcribedAddResponse model.SubcribedAppResponse

	if err := json.Unmarshal(body, &subcribedAddResponse); err != nil {
		return nil, err
	}

	return &subcribedAddResponse, nil
}

// TODO: handle errors
func (f *FbClient) CallAPISendMessage(accessToken string, sendMessageRequest *model.SendMessageRequest) (*model.SendMessageResponse, error) {
	URL, err := url.Parse(fmt.Sprintf("%s/me/messages", f.apiInfo.Url()))
	if err != nil {
		return nil, err
	}

	query, err := url.ParseQuery(URL.RawQuery)
	if err != nil {
		return nil, err
	}

	recipient, err := json.Marshal(sendMessageRequest.Recipient)
	if err != nil {
		return nil, err
	}
	message, err := json.Marshal(sendMessageRequest.Message)
	if err != nil {
		return nil, err
	}

	query.Add(AccessToken, accessToken)
	query.Add(Recipient, string(recipient))
	query.Add(Message, string(message))
	URL.RawQuery = query.Encode()

	fmt.Println(URL.String())
	resp, err := http.Post(URL.String(), "", nil)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	fmt.Println(query.Encode())

	if err := f.facebookErrorService.HandleErrorFacebookAPI(body, URL.String()); err != nil {
		return nil, err
	}

	var sendMessageResponse model.SendMessageResponse

	if err := json.Unmarshal(body, &sendMessageResponse); err != nil {
		return nil, err
	}

	return &sendMessageResponse, nil
}

func (f *FbClient) CallAPISendComment(accessToken string, sendCommentRequest *model.SendCommentRequest) (*model.SendCommentResponse, error) {
	URL, err := url.Parse(fmt.Sprintf("%s/%s/comments", f.apiInfo.Url(), sendCommentRequest.ID))
	if err != nil {
		return nil, err
	}

	query, err := url.ParseQuery(URL.RawQuery)
	if err != nil {
		return nil, err
	}

	query.Add(AccessToken, accessToken)
	if sendCommentRequest.Message != "" {
		query.Add(Message, sendCommentRequest.Message)
	}
	if sendCommentRequest.AttachmentURL != "" {
		query.Add(AttachmentURL, sendCommentRequest.AttachmentURL)
	}
	URL.RawQuery = query.Encode()

	resp, err := http.Post(URL.String(), "", nil)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if err := f.facebookErrorService.HandleErrorFacebookAPI(body, URL.String()); err != nil {
		return nil, err
	}

	var sendCommentResponse model.SendCommentResponse

	if err := json.Unmarshal(body, &sendCommentResponse); err != nil {
		return nil, err
	}

	return &sendCommentResponse, nil
}

func (f *FbClient) CallAPICommentByID(accessToken, commentID string) (*model.Comment, error) {
	URL, err := url.Parse(fmt.Sprintf("%s/%s", f.apiInfo.Url(), commentID))
	if err != nil {
		return nil, err
	}

	query, err := url.ParseQuery(URL.RawQuery)
	if err != nil {
		return nil, err
	}

	query.Add(AccessToken, accessToken)
	query.Add(Fields, "message,attachment,id,created_time,comment_count,parent,from,is_hidden")
	query.Add(DateFormat, UnixDateFormat)

	URL.RawQuery = query.Encode()
	resp, err := http.Get(URL.String())
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if err := f.facebookErrorService.HandleErrorFacebookAPI(body, URL.String()); err != nil {
		return nil, err
	}

	var comment model.Comment

	if err := json.Unmarshal(body, &comment); err != nil {
		return nil, err
	}

	return &comment, nil
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
