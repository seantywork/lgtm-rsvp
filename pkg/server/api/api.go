package api

import (
	"encoding/json"
	"fmt"
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
	AppKey   string `json:"app_key"`
	MailPass string `json:"mail_pass"`
}

var API_JSON ApiJSON

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

	return nil
}

func GetAppKey(c *gin.Context) {

	c.JSON(http.StatusOK, SERVER_RESP{Status: "success", Reply: API_JSON.AppKey})

	return
}

func GetGiftPage(c *gin.Context) {

	c.JSON(http.StatusOK, SERVER_RESP{Status: "success", Reply: pkgglob.G_CONF.GiftPage})

	return
}
