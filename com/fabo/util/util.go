package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"etop.vn/backend/cmd/fabo/config"
	"etop.vn/backend/com/fabo/api"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/extservice/telebot"
)

var (
	apiInfo config.ApiInfo
	appInfo config.AppInfo
	bot     *telebot.Channel
)

// TODO: check errors

func New(_apiInfo config.ApiInfo, _appInfo config.AppInfo, _bot *telebot.Channel) {
	apiInfo = _apiInfo
	appInfo = _appInfo
	bot = _bot
}

func CallAPIGetMe(accessToken string) (*api.Me, error) {
	URL, err := url.Parse(fmt.Sprintf("%s/me", apiInfo.Url()))
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

	var me api.Me
	if err := json.Unmarshal(body, &me); err != nil {
		return nil, err
	}

	return &me, nil
}

func CallAPIGetAccounts(accessToken string) (*api.Accounts, error) {
	URL, err := url.Parse(fmt.Sprintf("%s/me", apiInfo.Url()))
	if err != nil {
		return nil, err
	}

	query, err := url.ParseQuery(URL.RawQuery)
	if err != nil {
		return nil, err
	}

	query.Add(AccessTokenKey, accessToken)
	query.Add(Fields, "accounts{access_token,category,category_list,name,id,tasks,description,about,fan_count}")
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

	var accounts api.Accounts

	if err := json.Unmarshal(body, &accounts); err != nil {
		return nil, err
	}
	return &accounts, nil
}

func CallAPIGetLongLivedAccessToken(accessToken string) (*api.Token, error) {
	URL, err := url.Parse(fmt.Sprintf("%s/oauth/access_token", apiInfo.Url()))
	if err != nil {
		return nil, err
	}

	query, err := url.ParseQuery(URL.RawQuery)
	if err != nil {
		return nil, err
	}

	query.Add(GrantType, GrantTypeFBExchangeToken)
	query.Add(FBExchangeToken, accessToken)
	query.Add(ClientIDKey, appInfo.AppID)
	query.Add(ClientSecretKey, appInfo.AppSecret)

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

	var tok api.Token

	if err := json.Unmarshal(body, &tok); err != nil {
		return nil, err
	}
	return &tok, nil
}

func CallAPICheckAccessToken(accessToken string) (*api.UserToken, error) {
	URL, err := url.Parse(fmt.Sprintf("%s/debug_token", apiInfo.Url()))
	if err != nil {
		return nil, err
	}

	query, err := url.ParseQuery(URL.RawQuery)
	if err != nil {
		return nil, err
	}

	query.Add(AccessTokenKey, appInfo.AppAccessToken)
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

	var tok api.UserToken

	if err := json.Unmarshal(body, &tok); err != nil {
		return nil, err
	}

	if tok.Data.Error != nil {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "").
			WithMetaM(map[string]string{
				"code":    fmt.Sprintf("%v", tok.Data.Error.Code),
				"message": tok.Data.Error.Message,
			})
	}
	return &tok, nil
}
