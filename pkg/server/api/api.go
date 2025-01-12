package api

import (
	"encoding/json"
	"fmt"
	"os"
)

type CLIENT_REQ struct {
	Data string `json:"data"`
}

type SERVER_RESP struct {
	Status string `json:"status"`
	Reply  string `json:"reply"`
}

type MailJSON struct {
	Pass string `json:"pass"`
}

var MAIL_JSON MailJSON

func InitAPI() error {

	fileb, err := os.ReadFile("mail.json")

	if err != nil {

		return fmt.Errorf("failed to init api: %s", err.Error())
	}

	mj := MailJSON{}

	err = json.Unmarshal(fileb, &mj)

	if err != nil {

		return fmt.Errorf("failed to init api: unmarshal: %s", err.Error())
	}

	MAIL_JSON = mj

	return nil
}
