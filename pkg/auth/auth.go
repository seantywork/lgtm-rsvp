package auth

import (
	"fmt"
	"log"

	pkgdb "our-wedding-rsvp/pkg/db"

	pkgglob "our-wedding-rsvp/pkg/glob"

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
