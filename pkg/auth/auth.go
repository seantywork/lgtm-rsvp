package auth

import (
	"fmt"
	"log"

	pkgdb "lgtm-rsvp/pkg/db"

	pkgglob "lgtm-rsvp/pkg/glob"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"

	pkgutils "lgtm-rsvp/pkg/utils"
)

var USE_OAUTH2 bool = false

var ADMINS = make(map[string]string)

func InitAuth() error {

	if pkgglob.G_CONF.Admin.OAuth != nil {
		USE_OAUTH2 = true
	} else {
		return nil
	}

	oj, err := GetOAuthJSON()

	if err != nil {
		return err
	}

	OAUTH_JSON = oj

	GoogleOauthConfig, err = GenerateGoogleOauthConfig()

	if err != nil {
		return err
	}

	return nil
}

func GenerateStateAuthCookie(c *gin.Context) string {

	state, _ := pkgutils.GetRandomHex(64)

	session := sessions.Default(c)

	session.Set(pkgglob.G_CONF.SessionStore, state)
	err := session.Save()

	if err != nil {
		log.Printf("cookie gen failed: %s\n", err.Error())
	}

	return state
}

func RegisterAdmin(id string, pw string) error {

	err := pkgdb.UpsertAdmin(id, pw)

	if err != nil {
		return fmt.Errorf("failed to register admin")
	}

	return nil

}
