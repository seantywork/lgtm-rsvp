package auth

import (
	"log"
	pkgdb "our-wedding-rsvp/pkg/db"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Is0(c *gin.Context) bool {

	session := sessions.Default(c)

	var session_id string

	v := session.Get("RSVP")

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

	return true
}
