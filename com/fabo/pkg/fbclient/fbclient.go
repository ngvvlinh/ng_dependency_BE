package fbclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"o.o/backend/com/fabo/pkg/fbclient/model"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/extservice/telebot"
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
	bot                  *telebot.Channel
}

func New(_appInfo AppConfig, _bot *telebot.Channel) *FbClient {
	return &FbClient{
		appInfo: _appInfo,
		apiInfo: ApiInfo{
			Host:    "https://graph.facebook.com",
			Version: "v6.0",
		},
		bot:                  _bot,
		facebookErrorService: NewFacebookErrorService(_bot),
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

func (f *FbClient) CallAPIPublishedPosts(pageID, accessToken string, pagination *model.FacebookPagingRequest) (*model.PublishedPostsResponse, error) {
	URL, err := url.Parse(fmt.Sprintf("%s/%s", f.apiInfo.Url(), pageID))
	if err != nil {
		return nil, err
	}

	query, err := url.ParseQuery(URL.RawQuery)
	if err != nil {
		return nil, err
	}

	defaultPaging := 100
	if pagination != nil && pagination.Limit.Valid {
		defaultPaging = pagination.Limit.Int
	}

	query.Add(AccessToken, accessToken)
	query.Add(Fields, fmt.Sprintf("published_posts.limit(%d){id,created_time,from,full_picture,icon,is_expired,is_hidden,is_popular,is_published,message,permalink_url,shares,status_type,updated_time,picture,attachments{media_type,type,subattachments}}", defaultPaging))
	query.Add(DateFormat, UnixDateFormat)

	URL.RawQuery = query.Encode()
	resp, err := http.Get(pagination.AddQueryParams(URL.String(), false, 0))
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

func (f *FbClient) CallAPIGetPost(pageID, postID, accessToken string) (*model.Post, error) {
	URL, err := url.Parse(fmt.Sprintf("%s/%s_%s", f.apiInfo.Url(), pageID, postID))
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

func (f *FbClient) CallAPIGetComments(parentID, ID, accessToken string, pagination *model.FacebookPagingRequest) (*model.CommentsResponse, error) {
	URL, err := url.Parse(fmt.Sprintf("%s/%s_%s/comments", f.apiInfo.Url(), parentID, ID))
	if err != nil {
		return nil, err
	}

	query, err := url.ParseQuery(URL.RawQuery)
	if err != nil {
		return nil, err
	}

	query.Add(AccessToken, accessToken)
	query.Add(Fields, "message,id,from,attachment,comment_count,parent,created_time")
	query.Add(DateFormat, UnixDateFormat)

	URL.RawQuery = query.Encode()
	resp, err := http.Get(pagination.AddQueryParams(URL.String(), true, 100))
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

func (f *FbClient) CallAPIGetConversations(pageID, accessToken string, pagination *model.FacebookPagingRequest) (*model.ConversationsResponse, error) {
	URL, err := url.Parse(fmt.Sprintf("%s/%s", f.apiInfo.Url(), pageID))
	if err != nil {
		return nil, err
	}

	query, err := url.ParseQuery(URL.RawQuery)
	if err != nil {
		return nil, err
	}

	defaultPaging := 100
	if pagination != nil && pagination.Limit.Valid {
		defaultPaging = pagination.Limit.Int
	}

	query.Add(AccessToken, accessToken)
	query.Add(Fields, fmt.Sprintf("conversations.limit(%d){message_count,id,updated_time,link}", defaultPaging))
	query.Add(DateFormat, UnixDateFormat)

	URL.RawQuery = query.Encode()
	resp, err := http.Get(pagination.AddQueryParams(URL.String(), false, 100))
	if err != nil {
		return nil, err
	}

	fmt.Println(URL.String())
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

func (f *FbClient) CallAPIGetMessages(conversationID, accessToken string, pagination *model.FacebookPagingRequest) (*model.MessagesResponse, error) {
	URL, err := url.Parse(fmt.Sprintf("%s/%s", f.apiInfo.Url(), conversationID))
	if err != nil {
		return nil, err
	}

	defaultPaging := 100
	if pagination != nil && pagination.Limit.Valid {
		defaultPaging = pagination.Limit.Int
	}

	query, err := url.ParseQuery(URL.RawQuery)
	if err != nil {
		return nil, err
	}

	query.Add(AccessToken, accessToken)
	query.Add(Fields, fmt.Sprintf("messages.limit(%d){from,to,message,created_time,id}", defaultPaging))
	query.Add(DateFormat, UnixDateFormat)

	URL.RawQuery = query.Encode()
	resp, err := http.Get(pagination.AddQueryParams(URL.String(), false, 100))
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
