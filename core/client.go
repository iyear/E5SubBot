package core

import (
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	TgId         int64  `gorm:"column:tg_id"`
	RefreshToken string `gorm:"column:refresh_token"`
	MsId         string `gorm:"column:ms_id"`
	Uptime       int64  `gorm:"column:uptime"`
	Alias        string `gorm:"column:alias"`
	ClientId     string `gorm:"column:client_id"`
	ClientSecret string `gorm:"column:client_secret"`
	Other        string `gorm:"column:other"`
}

const (
	msApiUrl    string = "https://login.microsoftonline.com"
	msGraUrl    string = "https://graph.microsoft.com"
	redirectUri string = "http://localhost/e5sub"
	scope       string = "openid offline_access mail.read user.read"
)

func NewClient(clientId string, clientSecret string) *Client {
	return &Client{
		ClientId:     clientId,
		ClientSecret: clientSecret,
	}
}
func GetMSAuthUrl(clientId string) string {
	return "https://login.microsoftonline.com/common/oauth2/v2.0/authorize?client_id=" + clientId + "&response_type=code&redirect_uri=" + url.QueryEscape(redirectUri) + "&response_mode=query&scope=" + url.QueryEscape(scope)
}
func GetMSRegisterAppUrl() string {
	ru := "https://developer.microsoft.com/en-us/graph/quick-start?appID=_appId_&appName=_appName_&redirectUrl=http://localhost:8000&platform=option-windowsuniversal"
	deeplink := "/quickstart/graphIO?publicClientSupport=false&appName=e5sub&redirectUrl=http://localhost/e5sub&allowImplicitFlow=false&ru=" + url.QueryEscape(ru)
	appUrl := "https://apps.dev.microsoft.com/?deepLink=" + url.QueryEscape(deeplink)
	return appUrl
}

//return access_token and refresh_token
func (c *Client) GetTokenWithCode(code string) (error error) {
	var r http.Request
	client := &http.Client{}
	r.ParseForm()
	r.Form.Add("client_id", c.ClientId)
	r.Form.Add("client_secret", c.ClientSecret)
	r.Form.Add("grant_type", "authorization_code")
	r.Form.Add("scope", scope)
	r.Form.Add("code", code)
	r.Form.Add("redirect_uri", redirectUri)
	body := strings.NewReader(r.Form.Encode())
	req, err := http.NewRequest("POST", msApiUrl+"/common/oauth2/v2.0/token", body)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if gjson.Get(string(content), "token_type").String() == "Bearer" {
		c.RefreshToken = gjson.Get(string(content), "refresh_token").String()
		return nil
	}
	return errors.New(string(content))
}

//return access_token and new refresh token
func (c *Client) getToken() (access string) {
	var r http.Request
	client := &http.Client{}
	r.ParseForm()
	r.Form.Add("client_id", c.ClientId)
	r.Form.Add("client_secret", c.ClientSecret)
	r.Form.Add("grant_type", "refresh_token")
	r.Form.Add("scope", scope)
	r.Form.Add("refresh_token", c.RefreshToken)
	r.Form.Add("redirect_uri", redirectUri)
	body := strings.NewReader(r.Form.Encode())
	//fmt.Println(body)
	req, err := http.NewRequest("POST", msApiUrl+"/common/oauth2/v2.0/token", body)
	if err != nil {
		return ""
	}
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	//fmt.Println(string(content))
	//fmt.Println(gjson.Get(string(content), "access_token").String())
	if gjson.Get(string(content), "token_type").String() == "Bearer" {
		c.RefreshToken = gjson.Get(string(content), "refresh_token").String()
		return gjson.Get(string(content), "access_token").String()
	}
	return ""
}

//Get User's Information
func (c *Client) GetUserInfo() (json string, error error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", msGraUrl+"/v1.0/me", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", c.getToken())
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if gjson.Get(string(content), "id").String() != "" {
		//fmt.Println("UserName: " + gjson.Get(string(content), "displayName").String())
		return string(content), nil
	}
	return "", errors.New(string(content))
}

func (c *Client) GetOutlookMails() error {
	client := http.Client{}
	req, err := http.NewRequest("GET", msGraUrl+"/v1.0/me/messages", nil)

	if err != nil {
		return err
	}
	req.Header.Set("Authorization", c.getToken())
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	//fmt.Println(string(content))
	//这里的.需要转义，否则会按路径的方式解析
	if gjson.Get(string(content), "@odata\\.context").String() != "" {
		return nil
	}
	return errors.New(gjson.Get(string(content), "error").String())
}
