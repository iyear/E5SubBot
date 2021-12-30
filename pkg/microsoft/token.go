package microsoft

import (
	"fmt"
	"github.com/guonaihong/gout"
	"github.com/tidwall/gjson"
)

// GetTokenWithCode return refreshToken
func GetTokenWithCode(id, secret, code string) (string, error) {
	var content string
	err := gout.POST(apiURL + "/common/oauth2/v2.0/token").
		SetWWWForm(gout.H{
			"client_id":     id,
			"client_secret": secret,
			"grant_type":    "authorization_code",
			"scope":         scope,
			"code":          code,
			"redirect_uri":  redirect,
		}).
		BindBody(&content).Do()

	if err != nil {
		return "", err
	}

	if gjson.Get(content, "token_type").String() == "Bearer" {
		return gjson.Get(content, "refresh_token").String(), nil
	}
	return "", fmt.Errorf("wrong token type")
}

// GetToken return new refreshToken, accessToken and error
func GetToken(id, secret, refresh string) (string, string, error) {
	var content string
	err := gout.POST(apiURL + "/common/oauth2/v2.0/token").
		SetWWWForm(gout.H{
			"client_id":     id,
			"client_secret": secret,
			"grant_type":    "refresh_token",
			"scope":         scope,
			"refresh_token": refresh,
			"redirect_uri":  redirect,
		}).
		BindBody(&content).
		Do()
	if err != nil {
		return "", "", err
	}

	if gjson.Get(content, "token_type").String() == "Bearer" {
		return gjson.Get(content, "refresh_token").String(),
			gjson.Get(content, "access_token").String(),
			nil
	}
	return "", "", fmt.Errorf("wrong token type")
}
