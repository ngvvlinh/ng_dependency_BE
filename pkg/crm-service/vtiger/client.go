package vtiger

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/schema"
)

var encoder = schema.NewEncoder()

func init() {
	encoder.SetAliasTag("json")
}

// SendRequest use api
func (c *Client) SendRequest(path string, body interface{}) (*http.Response, error) {
	var bodyData []byte
	if body == nil {
		bodyData = []byte("{}")
	} else {
		var err error
		bodyData, err = json.Marshal(body)
		if err != nil {
			panic(err) // should panic because error is unexpected
		}
	}

	requestURL := c.BaseURL + path

	log.Println(requestURL)
	request, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(bodyData))
	if err != nil {
		return nil, err
	}

	request.Header.Add("Authorization", "Bearer "+c.Token)
	request.Header.Add("Content-Type", "application/json")

	resp, err := c.httpClient.Do(request)
	return resp, err
}

// GetMD5Hash string to md5
func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

// GetSessionKey which use in excute vtiger API
func GetSessionKey(vtigerService string, vtigerUser string, vtigerAccessKey string) (*VtigerSessionResult, error) {

	// get vtiger token
	requestURL := fmt.Sprintf("%v?operation=getchallenge&username=%v", vtigerService, vtigerUser)
	fmt.Println(requestURL)
	resp, err := http.Get(requestURL)
	if err != nil {
		fmt.Print("err 1")
		return nil, err
	}
	defer resp.Body.Close()
	bodyToken, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Print("err 2")
		panic(err.Error())
	}

	bodyString := string(bodyToken)

	fmt.Println("body string :" + bodyString)

	var responseVtigerMap VtigerServiceTokenResponse

	err = json.Unmarshal(bodyToken, &responseVtigerMap)

	if err != nil {
		fmt.Print("err 3")
		return nil, err
	}

	// get session
	accessKey := GetMD5Hash(responseVtigerMap.Result.Token + vtigerAccessKey)
	bodySessionRequest := BodySessionRequest{
		Operation: "login",
		Username:  vtigerUser,
		AccessKey: accessKey,
	}

	fmt.Println("token :" + accessKey)

	data := url.Values{}
	if err := encoder.Encode(bodySessionRequest, data); err != nil {
		panic("loi")
	}

	sessionRequest, err := http.NewRequest("POST", vtigerService, strings.NewReader(data.Encode()))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	sessionRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{Timeout: 10 * time.Second}

	// send request and get response about session
	sessionResponse, err := client.Do(sessionRequest)
	bodySession, err := ioutil.ReadAll(sessionResponse.Body)
	if err != nil {
		return nil, err
	}
	defer sessionResponse.Body.Close()
	var vtigerSessionResponse VtigerSessionResponse
	err = json.Unmarshal(bodySession, &vtigerSessionResponse)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// return sessionName use in access vtiger
	return vtigerSessionResponse.Result, nil
}

// SendRequestVtigerValue send request to vtiger
func (v VtigerClient) SendRequestVtigerValue(path string, body interface{}, method string) (*VtigerResponse, error) {

	var bodyData []byte
	if body == nil {
		bodyData = []byte("{}")
	} else {
		var err error
		bodyData, err = json.Marshal(body)
		if err != nil {
			panic(err) // should panic because error is unexpected
		}
	}

	requestURL := v.BaseURL + path

	log.Println("request url : " + requestURL)
	var request *http.Request
	var err error
	if strings.ToUpper(method) == "GET" {
		request, err = http.NewRequest("GET", requestURL, nil)
	} else {
		request, err = http.NewRequest(method, requestURL, bytes.NewBuffer(bodyData))
	}

	if err != nil {
		return nil, err
	}

	resp, err := v.httpClient.Do(request)

	bodyResponse, err := ioutil.ReadAll(resp.Body)

	log.Println("Vtiger response:" + string(bodyResponse))

	var vtigerResponse *VtigerResponse

	err = json.Unmarshal([]byte(bodyResponse), &vtigerResponse)

	return vtigerResponse, nil
}
