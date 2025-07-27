package server

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	pkgauth "lgtm-rsvp/pkg/auth"
	pkgdb "lgtm-rsvp/pkg/db"
	pkgglob "lgtm-rsvp/pkg/glob"
	pkgserverapi "lgtm-rsvp/pkg/server/api"
)

func getIndex(c *gin.Context) {

	miName := pkgserverapi.GetMainImage().Name

	c.HTML(200, "index.html", gin.H{
		"main_image":       miName,
		"title":            pkgglob.G_CONF.Title,
		"groom":            pkgglob.G_CONF.Groom,
		"bride":            pkgglob.G_CONF.Bride,
		"comment":          pkgglob.G_CONF.Comment,
		"message":          pkgglob.G_CONF.Message,
		"google_comment_y": pkgserverapi.GOOGLE_COMMENT_Y,
		"kakao_share_y":    pkgserverapi.KAKAO_SHARE_Y,
	})

}

func getSignin(c *gin.Context) {

	if pkgauth.USE_OAUTH2 {

		c.Redirect(302, "/api/oauth2/google/signin")

		return

	}

	c.HTML(200, "signin.html", gin.H{})
}

func getRead(c *gin.Context) {

	watchId := c.Param("storyId")

	if !pkgauth.VerifyDefaultValue(watchId) {

		log.Printf("get article: illegal: %s\n", watchId)

		c.JSON(http.StatusBadRequest, pkgserverapi.SERVER_RESP{Status: "error", Reply: "invalid format"})

		return

	}
	c.HTML(200, "read.html", gin.H{
		"article_code": watchId,
		"title":        pkgglob.G_CONF.Title,
	})

}

func getWrite(c *gin.Context) {

	if !pkgauth.Is0(c, nil, nil) {
		log.Printf("get write: illegal\n")

		c.JSON(http.StatusBadRequest, pkgserverapi.SERVER_RESP{Status: "error", Reply: "you're not admin"})

		return
	}

	c.HTML(200, "write.html", gin.H{
		"title": pkgglob.G_CONF.Title,
	})
}

func Logout(c *gin.Context) {

	userId := ""

	sessionId := ""

	if !pkgauth.Is0(c, &userId, &sessionId) {

		log.Printf("logout: not logged in\n")

		c.JSON(http.StatusBadRequest, pkgserverapi.SERVER_RESP{Status: "error", Reply: "not logged in"})

		return

	}

	err := pkgdb.SetAdminSessionId(userId, sessionId, false)

	if err != nil {

		log.Printf("logout: error: %v\n", err)

		c.JSON(http.StatusBadRequest, pkgserverapi.SERVER_RESP{Status: "error", Reply: "invalid"})

		return

	}

	c.JSON(http.StatusOK, pkgserverapi.SERVER_RESP{Status: "success", Reply: "logged out"})

}

func DeleteStory(c *gin.Context) {

	if !pkgauth.Is0(c, nil, nil) {

		log.Printf("delete: not allowed\n")

		c.JSON(http.StatusBadRequest, pkgserverapi.SERVER_RESP{Status: "error", Reply: "invalid"})

		return

	}

	storyId := c.Param("storyId")

	err := pkgdb.DeleteStoryById(storyId)

	if err != nil {
		log.Printf("delete: error: %v\n", err)

		c.JSON(http.StatusInternalServerError, pkgserverapi.SERVER_RESP{Status: "error", Reply: "internal error"})

		return
	}

	c.JSON(http.StatusOK, pkgserverapi.SERVER_RESP{Status: "success", Reply: "deleted"})

}

func GetComment(c *gin.Context) {

	c.HTML(200, "comment.html", gin.H{
		"title": pkgglob.G_CONF.Title,
	})
}

func GetCommentSudo(c *gin.Context) {

	if !pkgauth.Is0(c, nil, nil) {

		log.Printf("comment sudo: not allowed\n")

		c.JSON(http.StatusBadRequest, pkgserverapi.SERVER_RESP{Status: "error", Reply: "invalid"})

		return

	}

	c.HTML(200, "comment_sudo.html", gin.H{
		"title": pkgglob.G_CONF.Title,
	})
}
