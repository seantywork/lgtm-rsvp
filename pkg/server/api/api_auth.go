package api

import (
	"encoding/json"
	"log"
	"net/http"

	pkgauth "lgtm-rsvp/pkg/auth"
	pkgdb "lgtm-rsvp/pkg/db"
	pkgglob "lgtm-rsvp/pkg/glob"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

type UserLogin struct {
	Id         string `json:"id"`
	Passphrase string `json:"passphrase"`
}

func OauthGoogleLogin(c *gin.Context) {

	if pkgauth.Is0(c, nil, nil) {

		log.Printf("oauth login: already logged in\n")

		c.JSON(http.StatusBadRequest, SERVER_RESP{Status: "error", Reply: "already logged in"})

		return

	}

	oauth_state := pkgauth.GenerateStateAuthCookie(c)

	u := pkgauth.GoogleOauthConfig.AuthCodeURL(oauth_state)

	c.Redirect(302, u)

}

func OauthGoogleCallback(c *gin.Context) {

	if pkgauth.Is0(c, nil, nil) {

		log.Printf("oauth callback: already logged in\n")

		c.JSON(http.StatusBadRequest, SERVER_RESP{Status: "error", Reply: "already logged in"})

		return

	}

	session := sessions.Default(c)

	var session_id string

	v := session.Get(pkgglob.G_CONF.SessionStore)

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
		log.Printf("access auth failed: not verified email\n")
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

func Login(c *gin.Context) {

	if pkgauth.USE_OAUTH2 {
		log.Printf("user login: use oauth2\n")

		c.JSON(http.StatusBadRequest, SERVER_RESP{Status: "error", Reply: "invalid"})

		return
	}

	var req CLIENT_REQ

	var u_login UserLogin

	if err := c.BindJSON(&req); err != nil {

		log.Printf("user login: failed to bind: %s\n", err.Error())

		c.JSON(http.StatusBadRequest, SERVER_RESP{Status: "error", Reply: "invalid format"})

		return
	}

	err := json.Unmarshal([]byte(req.Data), &u_login)

	if err != nil {

		log.Printf("user login: failed to unmarshal: %s\n", err.Error())

		c.JSON(http.StatusBadRequest, SERVER_RESP{Status: "error", Reply: "invalid format"})

		return

	}

	if !pkgauth.VerifyDefaultValue(u_login.Id) {

		log.Printf("user login: not valid id: %s\n", u_login.Id)

		c.JSON(http.StatusBadRequest, SERVER_RESP{Status: "error", Reply: "invalid format"})

		return
	}

	us, err := pkgdb.GetAdminById(u_login.Id)

	if err != nil {

		log.Printf("user login: failed to get from user: %s\n", err.Error())

		c.JSON(http.StatusForbidden, SERVER_RESP{Status: "error", Reply: "id doesn't exist"})

		return
	}

	if !us.Pw.Valid {
		log.Printf("user login: failed to get admin pw: null\n")

		c.JSON(http.StatusForbidden, SERVER_RESP{Status: "error", Reply: "id invalid"})

		return
	}

	if us.Pw.String != u_login.Passphrase {

		log.Printf("user login: passphrase: %s", "not matching")

		c.JSON(http.StatusForbidden, SERVER_RESP{Status: "error", Reply: "passphrase not matching"})

		return

	}

	session_key := pkgauth.GenerateStateAuthCookie(c)

	err = pkgdb.SetAdminSessionId(us.Id, session_key, true)

	if err != nil {

		log.Printf("user login: failed to set session: %s", err.Error())

		c.JSON(http.StatusInternalServerError, SERVER_RESP{Status: "error", Reply: "failed to login"})

		return

	}

	c.JSON(http.StatusOK, SERVER_RESP{Status: "success", Reply: "logged in"})

}
