package api

import (
	"log"
	"net/http"

	pkgglob "lgtm-rsvp/pkg/glob"

	"github.com/gin-gonic/gin"
)

type CLIENT_REQ struct {
	Data string `json:"data"`
}

type SERVER_RESP struct {
	Status string `json:"status"`
	Reply  string `json:"reply"`
}

var USE_GOOGLE_COMMENT bool = false
var USE_KAKAO_SHARE bool = false

var GOOGLE_COMMENT_Y = "y"

var KAKAO_SHARE_Y = "y"

func InitAPI() error {

	if pkgglob.G_CONF.Api.GoogleComment != nil {
		log.Printf("using google comment\n")
		USE_GOOGLE_COMMENT = true
	} else {
		log.Printf("not using google comment\n")
		GOOGLE_COMMENT_Y = ""
	}

	if pkgglob.G_CONF.Api.KakaoShare != nil {
		log.Printf("using kakao share\n")
		USE_KAKAO_SHARE = true
	} else {
		log.Printf("not using kakao share\n")
		KAKAO_SHARE_Y = ""
	}
	return nil
}

func GetKakaoShare(c *gin.Context) {

	c.JSON(http.StatusOK, SERVER_RESP{Status: "success", Reply: pkgglob.G_CONF.Title + ":" + *pkgglob.G_CONF.Api.KakaoShare})

	return
}

func GetGiftPage(c *gin.Context) {

	c.JSON(http.StatusOK, SERVER_RESP{Status: "success", Reply: pkgglob.G_CONF.GiftPage})

	return
}
