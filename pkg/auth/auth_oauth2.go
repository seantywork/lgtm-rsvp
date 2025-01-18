package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	pkgglob "our-wedding-rsvp/pkg/glob"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type OAuthJSON struct {
	Web struct {
		ClientID                string   `json:"client_id"`
		ProjectID               string   `json:"project_id"`
		AuthURI                 string   `json:"auth_uri"`
		TokenURI                string   `json:"token_uri"`
		AuthProviderX509CertURL string   `json:"auth_provider_x509_cert_url"`
		ClientSecret            string   `json:"client_secret"`
		RedirectUris            []string `json:"redirect_uris"`
	} `json:"web"`
}

const OauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

type OAuthStruct struct {
	ID             string `json:"id"`
	EMAIL          string `json:"email"`
	VERIFIED_EMAIL bool   `json:"verified_email"`
	PICTURE        string `json:"picture"`
}

var OAUTH_JSON OAuthJSON

var GoogleOauthConfig *oauth2.Config

func InitAuth() error {

	if !USE_OAUTH2 {
		return nil
	}
	oj, err := GetOAuthJSON()

	if err != nil {
		return err
	}

	OAUTH_JSON = oj

	GoogleOauthConfig = GenerateGoogleOauthConfig()

	return nil
}

func GetOAuthJSON() (OAuthJSON, error) {

	var cj OAuthJSON

	file_byte, err := os.ReadFile("oauth.json")

	if err != nil {
		return cj, fmt.Errorf("failed to read: oauth.json: %s", err.Error())
	}

	err = json.Unmarshal(file_byte, &cj)

	if err != nil {
		return cj, fmt.Errorf("failed to unmarshal oauth json: %s", err.Error())
	}

	return cj, nil

}

func GenerateGoogleOauthConfig() *oauth2.Config {

	google_oauth_config := &oauth2.Config{
		ClientID:     OAUTH_JSON.Web.ClientID,
		ClientSecret: OAUTH_JSON.Web.ClientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	google_oauth_config.RedirectURL = pkgglob.G_CONF.Url

	log.Println(google_oauth_config.RedirectURL)

	return google_oauth_config

}

func GetUserDataFromGoogle(code string) ([]byte, error) {

	token, err := GoogleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}
	response, err := http.Get(OauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}
	return contents, nil
}
