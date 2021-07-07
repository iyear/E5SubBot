package main

import (
	"errors"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	MsApiUrl    string = "https://login.microsoftonline.com"
	MsGraUrl    string = "https://graph.microsoft.com"
	redirectUri string = "http://localhost/e5sub"
	scope       string = "openid offline_access mail.read user.read"
)

func MSGetAuthUrl(cid string) string {
	return "https://login.microsoftonline.com/common/oauth2/v2.0/authorize?client_id=" + cid + "&response_type=code&redirect_uri=" + url.QueryEscape(redirectUri) + "&response_mode=query&scope=" + url.QueryEscape(scope)
}
func MSGetReAppUrl() string {
	ru := "https://developer.microsoft.com/en-us/graph/quick-start?appID=_appId_&appName=_appName_&redirectUrl=http://localhost:8000&platform=option-windowsuniversal"
	deeplink := "/quickstart/graphIO?publicClientSupport=false&appName=e5sub&redirectUrl=http://localhost/e5sub&allowImplicitFlow=false&ru=" + url.QueryEscape(ru)
	app_url := "https://apps.dev.microsoft.com/?deepLink=" + url.QueryEscape(deeplink)
	return app_url
}

//return access_token and refresh_token
func MSFirGetToken(code, cid, cse string) (access string, refresh string, Error error) {
	var r http.Request
	client := &http.Client{}
	r.ParseForm()
	r.Form.Add("client_id", cid)
	r.Form.Add("client_secret", cse)
	r.Form.Add("grant_type", "authorization_code")
	r.Form.Add("scope", scope)
	r.Form.Add("code", code)
	r.Form.Add("redirect_uri", redirectUri)
	body := strings.NewReader(r.Form.Encode())
	req, err := http.NewRequest("POST", MsApiUrl+"/common/oauth2/v2.0/token", body)
	if err != nil {
		logger.Println(err)
		return "", "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		logger.Println(err)
		return "", "", err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Println(err)
		return "", "", err
	}
	if gjson.Get(string(content), "token_type").String() == "Bearer" {
		return gjson.Get(string(content), "access_token").String(), gjson.Get(string(content), "refresh_token").String(), nil
	}
	return "", "", errors.New(string(content))
}

//return access_token and new refresh token
func MSGetToken(refreshtoken, cid, cse string) (access string, newRefreshToken string, Error error) {
	var r http.Request
	client := &http.Client{}
	r.ParseForm()
	r.Form.Add("client_id", cid)
	r.Form.Add("client_secret", cse)
	r.Form.Add("grant_type", "refresh_token")
	r.Form.Add("scope", scope)
	r.Form.Add("refresh_token", refreshtoken)
	r.Form.Add("redirect_uri", redirectUri)
	body := strings.NewReader(r.Form.Encode())
	//fmt.Println(body)
	req, err := http.NewRequest("POST", MsApiUrl+"/common/oauth2/v2.0/token", body)
	if err != nil {
		logger.Println(err)
		return "", "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		logger.Println(err)
		return "", "", err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Println(err)
		return "", "", err
	}
	//fmt.Println(string(content))
	//fmt.Println(gjson.Get(string(content), "access_token").String())
	if gjson.Get(string(content), "token_type").String() == "Bearer" {
		return gjson.Get(string(content), "access_token").String(), gjson.Get(string(content), "refresh_token").String(), nil
	}
	return "", "", errors.New(string(content))
}

//Get User's Information
func MSGetUserInfo(accesstoken string) (json string, Error error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", MsGraUrl+"/v1.0/me", nil)
	if err != nil {
		logger.Println(err)
		return "", err
	}
	req.Header.Set("Authorization", accesstoken)
	resp, err := client.Do(req)
	if err != nil {
		logger.Println(err)
		return "", err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Println(err)
		return "", err
	}
	if gjson.Get(string(content), "id").String() != "" {
		//fmt.Println("UserName: " + gjson.Get(string(content), "displayName").String())
		return string(content), nil
	}
	return "", errors.New(string(content))
}

func OutLookGetMails(accesstoken string) (bool, error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", MsGraUrl+"/v1.0/me/messages", nil)
	if err != nil {
		logger.Println(err)
		return false, err
	}
	req.Header.Set("Authorization", accesstoken)
	resp, err := client.Do(req)
	if err != nil {
		logger.Println(err)
		return false, err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Println(err)
		return false, err
	}
	//fmt.Println(string(content))
	//这里的.需要转义，否则会按路径的方式解析
	if gjson.Get(string(content), "@odata\\.context").String() != "" {
		return true, nil
	}
	return false, errors.New(string(content))
}
