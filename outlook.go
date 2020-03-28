package main

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	MsApiUrl string = "https://login.microsoftonline.com"
	MsGraUrl string = "https://graph.microsoft.com"
)

var (
	cliId       string
	redirectUri string
	scope       string
	authUrl     string
)

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	authUrl = viper.GetString("auth_url")
	cliId = GetURLValue(authUrl, "client_id")
	redirectUri, _ = url.QueryUnescape(GetURLValue(authUrl, "redirect_uri"))
	scope, _ = url.QueryUnescape(GetURLValue(authUrl, "scope"))
	//refreshtoken := "xxxx"
	//fmt.Println(MSGetUserInfo(MSGetToken(refreshtoken,"user.read mail.read")))
	//code := "xxx"
	//fmt.Println(MSFirGetToken(code))
}

//return access_token and refresh_token
func MSFirGetToken(code string) (access string, refresh string) {
	var r http.Request
	client := &http.Client{}
	r.ParseForm()
	r.Form.Add("client_id", cliId)
	r.Form.Add("grant_type", "authorization_code")
	r.Form.Add("scope", scope)
	r.Form.Add("code", code)
	r.Form.Add("redirect_uri", redirectUri)
	body := strings.NewReader(r.Form.Encode())
	req, err := http.NewRequest("POST", MsApiUrl+"/common/oauth2/v2.0/token", body)
	resp, err := client.Do(req)
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
	}
	if gjson.Get(string(content), "token_type").String() == "Bearer" {
		return gjson.Get(string(content), "access_token").String(), gjson.Get(string(content), "refresh_token").String()
	} else {
		return "", ""
	}
	return "", ""
}

//return access_token
func MSGetToken(refreshtoken string) (access string) {
	var r http.Request
	client := &http.Client{}
	r.ParseForm()
	r.Form.Add("client_id", cliId)
	r.Form.Add("grant_type", "refresh_token")
	r.Form.Add("scope", scope)
	r.Form.Add("refresh_token", refreshtoken)
	r.Form.Add("redirect_uri", redirectUri)
	body := strings.NewReader(r.Form.Encode())
	fmt.Println(body)
	req, err := http.NewRequest("POST", MsApiUrl+"/common/oauth2/v2.0/token", body)
	resp, err := client.Do(req)
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
	}
	fmt.Println(string(content))
	//fmt.Println(gjson.Get(string(content), "access_token").String())
	if gjson.Get(string(content), "token_type").String() == "Bearer" {
		return gjson.Get(string(content), "access_token").String()
	} else {
		return ""
	}
	return ""
}

//Get User's Information
func MSGetUserInfo(accesstoken string) (json string) {
	client := http.Client{}
	//r.Header.Set("Host","graph.microsoft.com")
	req, err := http.NewRequest("GET", MsGraUrl+"/v1.0/me", nil)
	if err != nil {
		fmt.Println("MSGetUserInfo ERROR ", err.Error())
		return ""
	}
	req.Header.Set("Authorization", accesstoken)
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)
	if gjson.Get(string(content), "id").String() != "" {
		fmt.Println("UserName: " + gjson.Get(string(content), "displayName").String())
		return string(content)
	}
	return ""
}

func OutLookGetMails(accesstoken string) bool {
	client := http.Client{}
	//r.Header.Set("Host","graph.microsoft.com")
	req, err := http.NewRequest("GET", MsGraUrl+"/v1.0/me/messages", nil)
	if err != nil {
		fmt.Println("MSGetMils ERROR ", err.Error())
		return false
	}
	req.Header.Set("Authorization", accesstoken)
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)
	if gjson.Get(string(content), "@odata.context").String() != "" {
		return true
	}
	return false
}
