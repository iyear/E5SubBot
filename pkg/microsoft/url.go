package microsoft

import (
	"fmt"
	"net/url"
)

const (
	apiURL   string = "https://login.microsoftonline.com"
	graphURL string = "https://graph.microsoft.com"
	redirect string = "http://localhost/e5sub"
	scope    string = "openid offline_access mail.read user.read"
)

func GetAuthURL(clientID string) string {
	return fmt.Sprintf(
		"https://login.microsoftonline.com/common/oauth2/v2.0/authorize?client_id=%s&response_type=code&redirect_uri=%s&response_mode=query&scope=%s",
		clientID,
		url.QueryEscape(redirect),
		url.QueryEscape(scope),
	)
}

func GetRegURL() string {
	ru := "https://developer.microsoft.com/en-us/graph/quick-start?appID=_appId_&appName=_appName_&redirectUrl=http://localhost:8000&platform=option-windowsuniversal"
	deeplink := fmt.Sprintf("/quickstart/graphIO?publicClientSupport=false&appName=e5sub&redirectUrl=%s&allowImplicitFlow=false&ru=%s", redirect, url.QueryEscape(ru))
	appUrl := fmt.Sprintf("https://apps.dev.microsoft.com/?deepLink=%s", url.QueryEscape(deeplink))
	return appUrl
}
