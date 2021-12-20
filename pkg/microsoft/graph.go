package microsoft

import (
	"fmt"
	"github.com/guonaihong/gout"
	"github.com/tidwall/gjson"
)

// GetUserInfo return newRefreshToken, infoJSON and error
func GetUserInfo(id, secret, refresh string) (string, string, error) {
	var (
		content string
	)
	newRefresh, access, err := GetToken(id, secret, refresh)
	if err != nil {
		return "", "", err
	}
	err = gout.GET(graphURL + "/v1.0/me").
		SetHeader(gout.H{
			"Authorization": access,
		}).
		BindBody(&content).
		Do()
	if err != nil {
		return "", "", err
	}

	if gjson.Get(content, "id").String() != "" {
		return newRefresh, content, nil
	}
	return "", "", fmt.Errorf(content)
}

func GetOutlookMails(id, secret, refresh string) (string, error) {
	var content string
	newRefresh, access, err := GetToken(id, secret, refresh)
	if err != nil {
		return "", err
	}

	err = gout.GET(graphURL + "/v1.0/me/messages").
		SetHeader(gout.H{
			"Authorization": access,
		}).
		BindBody(&content).
		Do()
	if err != nil {
		return "", err
	}
	// 这里的.需要转义，否则会按路径的方式解析
	if gjson.Get(content, "@odata\\.context").String() != "" {
		return newRefresh, nil
	}
	return "", fmt.Errorf(gjson.Get(content, "error").String())
}
