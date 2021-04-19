package linenotify

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/konaro/line-notify-service/model"
)

var client = &http.Client{}

const authorizeEndpoint = "https://notify-bot.line.me/oauth/authorize"
const tokenEndpoint = "https://notify-bot.line.me/oauth/token"
const notifyEndpoint = "https://notify-api.line.me/api/notify"

func GetAuthorizeUrl(clientId, redirectUri string) string {
	query := url.Values{
		"response_type": {"code"},
		"client_id":     {clientId},
		"redirect_uri":  {redirectUri},
		"scope":         {"notify"},
		"state":         {"nonce-123"},
		"response_mode": {"form_post"},
	}

	return authorizeEndpoint + "?" + query.Encode()
}

func GetAccessToken(code, clientId, clientSecret, redirectUri string) []byte {
	form := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {redirectUri},
		"client_id":     {clientId},
		"client_secret": {clientSecret},
	}

	req, _ := http.NewRequest("POST", tokenEndpoint, strings.NewReader(form.Encode()))

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

func RevokeAccessToken(token string) ([]byte, error) {
	req, _ := http.NewRequest("POST", notifyEndpoint+"/revoke", nil)

	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}

func PushNotification(message model.Notify, token string) ([]byte, error) {
	form := url.Values{
		"message": {message.Message},
	}
	req, _ := http.NewRequest("POST", notifyEndpoint, strings.NewReader(form.Encode()))

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))

	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}
