package server

import (
	"fmt"
	"html/template"
	pkgserverapi "lgtm-rsvp/pkg/server/api"
	"os"
	"path/filepath"

	pkgauth "lgtm-rsvp/pkg/auth"
	pkgglob "lgtm-rsvp/pkg/glob"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func CreateServerFromConfig() (*gin.Engine, error) {

	genserver := gin.Default()

	funcMap := make(template.FuncMap)

	funcMap["HTMLnoesc"] = func(s string) template.HTML {
		return template.HTML(s)
	}

	genserver.SetFuncMap(funcMap)

	store := sessions.NewCookieStore([]byte(pkgglob.G_CONF.SessionStore))

	so := sessions.Options{
		Path: "/",
	}

	store.Options(so)

	genserver.Use(sessions.Sessions(pkgglob.G_CONF.SessionStore, store))

	genserver.SetTrustedProxies(nil)

	err := pkgauth.InitAuth()

	if err != nil {

		return nil, err
	}

	err = pkgserverapi.InitAPI()

	if err != nil {

		return nil, err
	}

	if pkgserverapi.USE_GOOGLE_COMMENT {

		reterr := make(chan error)

		go pkgserverapi.StartMailer(reterr)

		re := <-reterr

		if re != nil {

			return nil, re
		}
	}

	err = configureServer(genserver)

	if err != nil {

		return nil, err
	}

	return genserver, nil

}

func configureServer(e *gin.Engine) error {

	albumPath := pkgglob.G_ALBUM_PATH

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

	e.GET("/story/r/:storyId", getRead)

	e.GET("/story/w", getWrite)

	e.GET("/story/r/:storyId/delete", DeleteStory)

	e.GET("/api/oauth2/google/signin", pkgserverapi.OauthGoogleLogin)

	e.GET("/api/oauth2/google/callback", pkgserverapi.OauthGoogleCallback)

	e.POST("/api/signin", pkgserverapi.Login)

	e.POST("/api/story/upload", pkgserverapi.UploadStory)

	e.GET("/api/story/download/:storyId", pkgserverapi.DownloadStoryById)

	e.POST("/api/media/upload", pkgserverapi.UploadStoryMedia)

	e.GET("/api/media/download/c/:mediaId", pkgserverapi.DownloadStoryMediaById)

	e.GET("/api/story/list", pkgserverapi.GetStoryList)

	e.GET("/api/image/list", pkgserverapi.GetImageList)

	e.GET("/api/gift", pkgserverapi.GetGiftPage)

	if pkgserverapi.USE_GOOGLE_COMMENT {

		e.GET("/comment", GetComment)

		e.GET("/comment/sudo", GetCommentSudo)

		e.GET("/api/comment/list", pkgserverapi.GetApprovedComments)

		e.POST("/api/comment/register", pkgserverapi.RegisterComment)

		e.GET("/api/comment/approve/:commentId", pkgserverapi.ApproveComment)

		e.POST("/api/comment/sudo/:cmd", pkgserverapi.CommentSudoCmd)
	}

	if pkgserverapi.USE_KAKAO_SHARE {

		e.GET("/api/appkey", pkgserverapi.GetKakaoShare)
	}

	return nil
}
