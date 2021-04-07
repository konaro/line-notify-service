package linenotify

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func GetAuthorizeUrl(clientId, redirectUri string) string {
	const endpoint = "https://notify-bot.line.me/oauth/authorize"

	query := url.Values{
		"response_type": {"code"},
		"client_id":     {clientId},
		"redirect_uri":  {redirectUri},
		"scope":         {"notify"},
		"state":         {"nonce-123"},
		"response_mode": {"form_post"},
	}

	return endpoint + "?" + query.Encode()
}

func GetAccessToken(code, clientId, clientSecret, redirectUri string) []byte {
	const endpoint = "https://notify-bot.line.me/oauth/token"

	client := &http.Client{}

	form := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {redirectUri},
		"client_id":     {clientId},
		"client_secret": {clientSecret},
	}

	req, _ := http.NewRequest("POST", endpoint, strings.NewReader(form.Encode()))

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))

	resp, err := client.Do(req)

	if err != nil {
		return nil
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil
	}

	return body
}
