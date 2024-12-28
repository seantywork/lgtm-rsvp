package auth

import (
	"fmt"
	"os"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"

	pkgutils "our-wedding-rsvp/pkg/utils"
)

var DEBUG bool = false

var USE_OAUTH2 bool = false

var ADMINS = make(map[string]string)

func GenerateStateAuthCookie(c *gin.Context) string {

	state, _ := pkgutils.GetRandomHex(64)

	session := sessions.Default(c)

	session.Set("RSVP", state)
	session.Save()

	return state
}

func RegisterAdmins(admins map[string]string) error {

	err := os.RemoveAll("./data/admin")

	if err != nil {
		return fmt.Errorf("failed to remove data/admin")
	}

	err = os.MkdirAll("./data/admin", 0755)

	if err != nil {

		return fmt.Errorf("failed to create data/admin")
	}

	for k, v := range admins {

		ADMINS[k] = v

		name := "./data/admin/" + k + ".json"

		err := os.WriteFile(name, []byte("{}"), 0644)

		if err != nil {

			return fmt.Errorf("failed to create data/admin: %s: %s", k, err.Error())
		}

	}

	return nil

}

/*

func OauthGoogleLogin(c *gin.Context) {

	my_key, my_type, _ := WhoAmI(c)

	if my_key != "" && my_type != "" {

		fmt.Printf("oauth login: already logged in\n")

		c.JSON(http.StatusBadRequest, com.SERVER_RE{Status: "error", Reply: "already logged in"})

		return

	}

	if !USE_OAUTH2 {

		c.Redirect(302, "/signin/idiot")

		return
	}

	oauth_state := GenerateStateAuthCookie(c)

	u := GoogleOauthConfig.AuthCodeURL(oauth_state)

	c.Redirect(302, u)

}

func OauthGoogleCallback(c *gin.Context) {

	my_key, my_type, _ := WhoAmI(c)

	if my_key != "" && my_type != "" {

		fmt.Printf("oauth callback: already logged in\n")

		c.JSON(http.StatusBadRequest, com.SERVER_RE{Status: "error", Reply: "already logged in"})

		return

	}

	session := sessions.Default(c)

	var session_id string

	v := session.Get("SOLIAGAIN")

	if v == nil {
		fmt.Printf("access auth failed: %s\n", "session id not found")
		return
	} else {
		session_id = v.(string)
	}

	state := c.Request.FormValue("state")

	if state == "" {
		fmt.Printf("access auth failed: %s\n", "form state not found")
		return
	}

	if state != session_id {
		fmt.Printf("access auth failed: %s\n", "value not matching")
		c.Redirect(302, "/signin")
		return
	}

	data, err := GetUserDataFromGoogle(c.Request.FormValue("code"))
	if err != nil {
		fmt.Printf("access auth failed: %s\n", err.Error())
		c.Redirect(302, "/signin")
		return
	}

	var oauth_struct OAuthStruct

	err = json.Unmarshal(data, &oauth_struct)

	if err != nil {
		fmt.Printf("access auth failed: %s\n", err.Error())
		c.Redirect(302, "/signin")
		return
	}

	if !oauth_struct.VERIFIED_EMAIL {
		fmt.Printf("access auth failed: %s\n", err.Error())
		c.Redirect(302, "/signin")
		return
	}

	if err != nil {
		fmt.Printf("access auth failed: %s\n", err.Error())
		c.Redirect(302, "/signin")
		return
	}

	as, err := dbquery.GetByIdFromAdmin(oauth_struct.EMAIL)

	if as == nil {

		fmt.Printf("access auth failed: %s", err.Error())

		c.Redirect(302, "/signin")

		return

	}

	err = dbquery.MakeSessionForAdmin(session_id, oauth_struct.EMAIL)

	if err != nil {

		fmt.Printf("make session failed for admin: %s", err.Error())

		c.Redirect(302, "/signin")

		return

	}

	fmt.Println("oauth sign in success")

	c.Redirect(302, "/")

}

*/
