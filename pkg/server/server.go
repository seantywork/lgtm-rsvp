package server

import (
	pkgserverapi "our-wedding-rsvp/pkg/server/api"

	pkgauth "our-wedding-rsvp/pkg/auth"
	pkgglob "our-wedding-rsvp/pkg/glob"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func CreateServerFromConfig() (*gin.Engine, error) {

	genserver := gin.Default()

	store := sessions.NewCookieStore([]byte(pkgglob.G_CONF.SessionStore))

	so := sessions.Options{
		Path: "/",
	}

	store.Options(so)

	genserver.Use(sessions.Sessions(pkgglob.G_CONF.SessionStore, store))

	pkgauth.USE_OAUTH2 = pkgglob.G_CONF.Admin.UseOauth2

	err := pkgauth.InitAuth()

	if err != nil {

		return nil, err
	}

	err = pkgserverapi.InitAPI()

	if err != nil {

		return nil, err
	}

	reterr := make(chan error)

	go pkgserverapi.StartMailer(reterr)

	re := <-reterr

	if re != nil {

		return nil, err
	}

	err = configureServer(genserver)

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

	e.GET("/comment", GetComment)

	e.GET("/story/r/:storyId", getRead)

	e.GET("/story/w", getWrite)

	e.GET("/story/r/:storyId/delete", DeleteStory)

	e.GET("/api/oauth2/google/signin", pkgserverapi.OauthGoogleLogin)

	e.GET("/api/oauth2/google/callback", pkgserverapi.OauthGoogleCallback)

	e.POST("/api/signin", pkgserverapi.Login)

	e.POST("/api/story/upload", pkgserverapi.UploadStory)

	e.GET("/api/story/download/:storyId", pkgserverapi.DownloadStoryById)

	e.GET("/api/comment/list", pkgserverapi.GetApprovedComments)

	e.POST("/api/comment/register", pkgserverapi.RegisterComment)

	e.GET("/api/comment/approve/:commentId", pkgserverapi.ApproveComment)

	e.POST("/api/media/upload", pkgserverapi.UploadStoryMedia)

	e.GET("/api/media/download/c/:mediaId", pkgserverapi.DownloadStoryMediaById)

	e.GET("/api/story/list", pkgserverapi.GetStoryList)

	e.GET("/api/appkey", pkgserverapi.GetAppKey)

	e.GET("/api/gift", pkgserverapi.GetGiftPage)

	return nil
}
