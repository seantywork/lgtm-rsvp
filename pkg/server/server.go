package server

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"

	pkgserverapi "our-wedding-rsvp/pkg/server/api"
)

func CreateServerFromConfig() (*gin.Engine, error) {

	genserver := gin.Default()

	store := sessions.NewCookieStore([]byte("RSVP"))

	genserver.Use(sessions.Sessions("RSVP", store))

	err := configureServer(genserver)

	if err != nil {

		return nil, err
	}

	return genserver, nil

}

func configureServer(e *gin.Engine) error {

	e.LoadHTMLGlob("view/*")

	e.Static("/public", "./public")

	e.GET("/", getIndex)

	e.GET("/signin", getSignin)

	e.GET("/signout", Logout)

	e.GET("/story/r/:storyId", getRead)

	e.GET("/story/w", getWrite)

	e.GET("/story/w/:storyId/delete", DeleteStory)

	e.GET("/api/oauth2/google/signin", pkgserverapi.OauthGoogleLogin)

	e.GET("/oauth2/google/callback", pkgserverapi.OauthGoogleCallback)

	e.POST("/api/signin", pkgserverapi.Login)

	e.POST("/api/story/upload", pkgserverapi.UploadStory)

	e.POST("/api/story/download/:storyId", pkgserverapi.DownloadStoryById)

	e.POST("/api/media/upload", pkgserverapi.UploadStoryMedia)

	e.GET("/api/media/download/c/:mediaId", pkgserverapi.DownloadStoryMediaById)

	e.GET("/api/story/list", pkgserverapi.GetStoryList)

	return nil
}
