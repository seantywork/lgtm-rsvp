package server

import (
	"github.com/gin-gonic/gin"

	pkgauth "our-wedding-rsvp/pkg/auth"
)

func getIndex(c *gin.Context) {

	c.HTML(200, "index.html", gin.H{})

}

func getSignin(c *gin.Context) {

	if pkgauth.USE_OAUTH2 {

		c.Redirect(302, "/api/oauth2/google/signin")

		return

	}

	c.HTML(200, "signin.html", gin.H{})
}

func getRead(c *gin.Context) {

	c.HTML(200, "read.html", gin.H{})
}

func getWrite(c *gin.Context) {

	if !pkgauth.Is0(c) {
		c.HTML(200, "index.html", gin.H{})
		return
	}

	c.HTML(200, "write.html", gin.H{})
}
