package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	pkgauth "our-wedding-rsvp/pkg/auth"
	pkgdb "our-wedding-rsvp/pkg/db"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

type CLIENT_REQ struct {
	Data string `json:"data"`
}

type SERVER_RESP struct {
	Status string `json:"status"`
	Reply  string `json:"reply"`
}

func OauthGoogleLogin(c *gin.Context) {

	if pkgauth.Is0(c) {

		log.Printf("oauth login: already logged in\n")

		c.JSON(http.StatusBadRequest, SERVER_RESP{Status: "error", Reply: "already logged in"})

		return

	}

	oauth_state := pkgauth.GenerateStateAuthCookie(c)

	u := pkgauth.GoogleOauthConfig.AuthCodeURL(oauth_state)

	c.Redirect(302, u)

}

func OauthGoogleCallback(c *gin.Context) {

	if pkgauth.Is0(c) {

		fmt.Printf("oauth callback: already logged in\n")

		c.JSON(http.StatusBadRequest, SERVER_RESP{Status: "error", Reply: "already logged in"})

		return

	}

	session := sessions.Default(c)

	var session_id string

	v := session.Get("RSVP")

	if v == nil {
		log.Printf("access auth failed: %s\n", "session id not found")
		return
	} else {
		session_id = v.(string)
	}

	state := c.Request.FormValue("state")

	if state == "" {
		log.Printf("access auth failed: %s\n", "form state not found")
		return
	}

	if state != session_id {
		log.Printf("access auth failed: %s\n", "value not matching")
		c.Redirect(302, "/")
		return
	}

	data, err := pkgauth.GetUserDataFromGoogle(c.Request.FormValue("code"))
	if err != nil {
		log.Printf("access auth failed: %s\n", err.Error())
		c.Redirect(302, "/")
		return
	}

	var oauth_struct pkgauth.OAuthStruct

	err = json.Unmarshal(data, &oauth_struct)

	if err != nil {
		log.Printf("access auth failed: %s\n", err.Error())
		c.Redirect(302, "/")
		return
	}

	if !oauth_struct.VERIFIED_EMAIL {
		log.Printf("access auth failed: %s\n", err.Error())
		c.Redirect(302, "/")
		return
	}

	if err != nil {
		log.Printf("access auth failed: %s\n", err.Error())
		c.Redirect(302, "/")
		return
	}

	as, err := pkgdb.GetAdminById(oauth_struct.EMAIL)

	if as == nil {

		log.Printf("access auth failed: %s", err.Error())

		c.Redirect(302, "/")

		return

	}

	err = pkgdb.SetAdminSessionId(oauth_struct.EMAIL, session_id, true)

	if err != nil {

		log.Printf("make session failed for admin: %s", err.Error())

		c.Redirect(302, "/")

		return

	}

	log.Println("oauth sign in success")

	c.Redirect(302, "/")

}
