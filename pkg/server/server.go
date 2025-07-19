package server

import (
	"fmt"
	"os"
	pkgserverapi "our-wedding-rsvp/pkg/server/api"
	"path/filepath"
	"strings"

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

	if !strings.HasPrefix(pkgglob.G_CONF.Album.Addr, "public/") {
		return fmt.Errorf("album addr needs to start with 'public' prefix")
	}

	albumPath := filepath.Join(pkgglob.G_CONF.Album.Addr, "album")

	if _, err := os.Stat(albumPath); err != nil {
		return fmt.Errorf("album addr not found: %v", err)
	}

	de, err := os.ReadDir(albumPath)

	if err != nil {
		return fmt.Errorf("failed to readdir: %v", err)
	}

	paths := make([]string, 0)

	delen := len(de)

	if delen < 3 {
		return fmt.Errorf("at least three images should exist: title, groom, bride")
	}

	for i := 0; i < delen; i++ {

		if de[i].IsDir() {
			return fmt.Errorf("album should not contain subdir")
		}

		path := filepath.Join(albumPath, de[i].Name())

		paths = append(paths, path)

	}

	pkgserverapi.AddImageList(paths)

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

	e.GET("/api/image/list", pkgserverapi.GetImageList)

	e.GET("/api/appkey", pkgserverapi.GetAppKey)

	e.GET("/api/gift", pkgserverapi.GetGiftPage)

	return nil
}
