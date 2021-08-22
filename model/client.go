package model

import (
	"github.com/guonaihong/gout"
	"github.com/iyear/E5SubBot/config"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"net/url"
)

type Client struct {
	ID           int    `gorm:"unique;primaryKey;not null"`
	TgId         int64  `gorm:"not null"`
	RefreshToken string `gorm:"not null"`
	MsId         string `gorm:"not null"`
	Uptime       int64  `gorm:"autoUpdateTime;not null"`
	Alias        string `gorm:"not null"`
	ClientId     string `gorm:"not null"`
	ClientSecret string `gorm:"not null"`
	Other        string
}
type ErrClient struct {
	*Client
	Err error
}

const (
	msApiUrl    string = "https://login.microsoftonline.com"
	msGraUrl    string = "https://graph.microsoft.com"
	redirectUri string = "http://localhost/e5sub"
	scope       string = "openid offline_access mail.read user.read"
)

func (c *Client) TableName() string {
	return config.Table
}
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

func (c *Client) GetTokenWithCode(code string) error {
	var content string
	err := gout.POST(msApiUrl + "/common/oauth2/v2.0/token").
		SetWWWForm(gout.H{
			"client_id":     c.ClientId,
			"client_secret": c.ClientSecret,
			"grant_type":    "authorization_code",
			"scope":         scope,
			"code":          code,
			"redirect_uri":  redirectUri,
		}).
		BindBody(&content).Do()

	if err != nil {
		return err
	}

	if gjson.Get(content, "token_type").String() == "Bearer" {
		c.RefreshToken = gjson.Get(content, "refresh_token").String()
		return nil
	}
	return errors.New(content)
}

// getToken return accessToken and error
func (c *Client) getToken() (string, error) {

	var content string
	err := gout.POST(msApiUrl + "/common/oauth2/v2.0/token").
		SetWWWForm(gout.H{
			"client_id":     c.ClientId,
			"client_secret": c.ClientSecret,
			"grant_type":    "refresh_token",
			"scope":         scope,
			"refresh_token": c.RefreshToken,
			"redirect_uri":  redirectUri,
		}).
		BindBody(&content).
		Do()
	if err != nil {
		return "", err
	}

	if gjson.Get(content, "token_type").String() == "Bearer" {
		c.RefreshToken = gjson.Get(content, "refresh_token").String()
		return gjson.Get(content, "access_token").String(), nil
	}
	return "", errors.New(gjson.Get(content, "error").String())
}

// GetUserInfo return infoJSON and error
func (c *Client) GetUserInfo() (string, error) {
	var (
		content     string
		err         error
		accessToken string
	)
	if accessToken, err = c.getToken(); err != nil {
		return "", err
	}
	err = gout.GET(msGraUrl + "/v1.0/me").
		SetHeader(gout.H{
			"Authorization": accessToken,
		}).
		BindBody(&content).
		Do()
	if err != nil {
		return "", err
	}

	if gjson.Get(content, "id").String() != "" {
		return content, nil
	}
	return "", errors.New(content)
}

func (c *Client) GetOutlookMails() error {

	var content string
	var err error
	var accessToken string
	if accessToken, err = c.getToken(); err != nil {
		return err
	}
	err = gout.GET(msGraUrl + "/v1.0/me/messages").
		SetHeader(gout.H{
			"Authorization": accessToken,
		}).
		BindBody(&content).
		Do()
	if err != nil {
		return err
	}
	// 这里的.需要转义，否则会按路径的方式解析
	if gjson.Get(content, "@odata\\.context").String() != "" {
		return nil
	}
	return errors.New(gjson.Get(content, "error").String())
}
