package fbclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"etop.vn/backend/com/fabo/pkg/fbclient/model"
	cc "etop.vn/backend/pkg/common/config"
	"etop.vn/backend/pkg/common/extservice/telebot"
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

func (f *FbClient) CallAPIGetMe(accessToken string) (*model.Me, error) {
	URL, err := url.Parse(fmt.Sprintf("%s/me", f.apiInfo.Url()))
	if err != nil {
		return nil, err
	}

	query, err := url.ParseQuery(URL.RawQuery)
	if err != nil {
		return nil, err
	}

	query.Add(AccessTokenKey, accessToken)
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

func (f *FbClient) CallAPIGetAccounts(accessToken string) (*model.Accounts, error) {
	URL, err := url.Parse(fmt.Sprintf("%s/me", f.apiInfo.Url()))
	if err != nil {
		return nil, err
	}

	query, err := url.ParseQuery(URL.RawQuery)
	if err != nil {
		return nil, err
	}

	query.Add(AccessTokenKey, accessToken)
	query.Add(Fields, "accounts{access_token,category,category_list,name,id,tasks,description,about,fan_count,picture}")
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

	var accounts model.Accounts

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
	query.Add(ClientSecretKey, f.appInfo.Secret)

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

	query.Add(AccessTokenKey, f.appInfo.AccessToken)
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
