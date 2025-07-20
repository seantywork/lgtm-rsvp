package auth

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	pkgglob "lgtm-rsvp/pkg/glob"

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

func GetOAuthJSON() (OAuthJSON, error) {

	var cj OAuthJSON

	cj.Web.ClientID = pkgglob.G_CONF.Admin.OAuth.ClientId
	cj.Web.ProjectID = pkgglob.G_CONF.Admin.OAuth.ProjectId
	cj.Web.AuthURI = pkgglob.G_CONF.Admin.OAuth.AuthUri
	cj.Web.TokenURI = pkgglob.G_CONF.Admin.OAuth.TokenUri
	cj.Web.AuthProviderX509CertURL = pkgglob.G_CONF.Admin.OAuth.AuthProviderX509CertUrl
	cj.Web.ClientSecret = pkgglob.G_CONF.Admin.OAuth.ClientSecret
	cj.Web.RedirectUris = pkgglob.G_CONF.Admin.OAuth.RidirectUris

	return cj, nil

}

func GenerateGoogleOauthConfig() (*oauth2.Config, error) {

	google_oauth_config := &oauth2.Config{
		ClientID:     OAUTH_JSON.Web.ClientID,
		ClientSecret: OAUTH_JSON.Web.ClientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	found := 0

	redirectlen := len(OAUTH_JSON.Web.RedirectUris)

	for i := 0; i < redirectlen; i++ {

		if strings.HasPrefix(OAUTH_JSON.Web.RedirectUris[i], pkgglob.G_CONF.Url) {

			google_oauth_config.RedirectURL = OAUTH_JSON.Web.RedirectUris[i]

			found = 1

			break
		}

	}

	if found == 0 {
		return nil, fmt.Errorf("failed to find redirect url")
	}

	log.Println(google_oauth_config.RedirectURL)

	return google_oauth_config, nil

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
