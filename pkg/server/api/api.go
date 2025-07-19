package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	pkgglob "our-wedding-rsvp/pkg/glob"

	"github.com/gin-gonic/gin"
)

type CLIENT_REQ struct {
	Data string `json:"data"`
}

type SERVER_RESP struct {
	Status string `json:"status"`
	Reply  string `json:"reply"`
}

type ApiJSON struct {
	GoogleComment string `json:"google_comment"`
	KakaoShare    string `json:"kakao_share"`
}

var API_JSON ApiJSON

var USE_GOOGLE_COMMENT bool = false
var USE_KAKAO_SHARE bool = false

var GOOGLE_COMMENT_EL = `
<div class="ww-section bg-light" id="comment">
	<div class="ww-photo-gallery">
		<div class="container">
			<div class="col text-center">
				<h2 class="h1 text-center pb-3 ww-title" style="font-family: 'Noto Serif KR', serif;">축하메시지</h2><br>
				<div class="row">
					<div class="col text-center">
						<button class="btn btn-primary btn-submit" type="submit" onclick="location.href='/comment'">메시지 남기러 가기</button>
					</div>
				</div>
				<br>
				<div id="comment-rows"></div>
				<br>
			</div>
		</div>
	</div>
</div>
`

var KAKAO_SHARE_EL = `
<a id="kakaotalk-sharing-btn" href="javascript:;">
	<img src="https://developers.kakao.com/assets/img/about/logos/kakaotalksharing/kakaotalk_sharing_btn_medium.png"
		alt="카카오톡 공유 보내기 버튼" />
</a> 
`

func InitAPI() error {

	fileb, err := os.ReadFile("api.json")

	if err != nil {

		return fmt.Errorf("failed to init api: %s", err.Error())
	}

	mj := ApiJSON{}

	err = json.Unmarshal(fileb, &mj)

	if err != nil {

		return fmt.Errorf("failed to init api: unmarshal: %s", err.Error())
	}

	API_JSON = mj

	if API_JSON.GoogleComment != "" {
		log.Printf("using google comment\n")
		USE_GOOGLE_COMMENT = true
	} else {
		log.Printf("not using google comment\n")
		GOOGLE_COMMENT_EL = ""
	}

	if API_JSON.KakaoShare != "" {
		log.Printf("using kakao share\n")
		USE_KAKAO_SHARE = true
	} else {
		log.Printf("not using kakao share\n")
		KAKAO_SHARE_EL = ""
	}
	return nil
}

func GetKakaoShare(c *gin.Context) {

	c.JSON(http.StatusOK, SERVER_RESP{Status: "success", Reply: API_JSON.KakaoShare})

	return
}

func GetGiftPage(c *gin.Context) {

	c.JSON(http.StatusOK, SERVER_RESP{Status: "success", Reply: pkgglob.G_CONF.GiftPage})

	return
}
