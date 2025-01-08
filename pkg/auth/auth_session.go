package auth

import (
	"log"
	pkgdb "our-wedding-rsvp/pkg/db"
	pkgglob "our-wedding-rsvp/pkg/glob"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Is0(c *gin.Context, userId *string, sessionId *string) bool {

	session := sessions.Default(c)

	var session_id string

	v := session.Get(pkgglob.G_CONF.SessionStore)

	if v == nil {
		return false

	} else {

		session_id = v.(string)

		if session_id == "" {
			return false
		}

	}

	s, err := pkgdb.GetAdminBySessionId(session_id)

	if err != nil {
		log.Printf("is0: err: %v\n", err)
		return false
	}

	if s == nil {
		return false
	}

	if userId != nil {
		*userId = s.Id
	}

	if sessionId != nil {

		*sessionId = session_id
	}

	return true
}
