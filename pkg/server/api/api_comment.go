package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"time"

	pkgauth "our-wedding-rsvp/pkg/auth"
	pkgdb "our-wedding-rsvp/pkg/db"
	pkgglob "our-wedding-rsvp/pkg/glob"
	pkgutils "our-wedding-rsvp/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
)

type CommentInfo struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func GetApprovedComments(c *gin.Context) {

	comments, err := pkgdb.ListApprovedComments()

	if err != nil {

		log.Printf("comments: list: %s\n", err.Error())

		c.JSON(http.StatusBadRequest, SERVER_RESP{Status: "error", Reply: "invalid format"})

		return

	}

	c_info := make([]CommentInfo, 0)

	clen := len(comments)

	for i := 0; i < clen; i++ {

		c_info = append(c_info, CommentInfo{
			Title:   comments[i].Title,
			Content: comments[i].Content,
		})
	}

	cbytes, err := json.Marshal(c_info)

	if err != nil {

		log.Printf("comments: marshal: %s\n", err.Error())

		c.JSON(http.StatusBadRequest, SERVER_RESP{Status: "error", Reply: "invalid format"})

		return
	}

	c.JSON(http.StatusOK, SERVER_RESP{Status: "success", Reply: string(cbytes)})

}

func RegisterComment(c *gin.Context) {

	var req CLIENT_REQ

	var c_info CommentInfo

	if err := c.BindJSON(&req); err != nil {

		log.Printf("register comment: failed to bind: %s\n", err.Error())

		c.JSON(http.StatusBadRequest, SERVER_RESP{Status: "error", Reply: "invalid format"})

		return
	}

	now := time.Now().UTC()

	timeRegistered := now.Format("2006-01-02-15-04-05")

	err := json.Unmarshal([]byte(req.Data), &c_info)

	if err != nil {

		log.Printf("c info comment: failed to unmarshal: %s\n", err.Error())

		c.JSON(http.StatusBadRequest, SERVER_RESP{Status: "error", Reply: "invalid format"})

		return

	}

	p := bluemonday.UGCPolicy()

	title_san := p.Sanitize(c_info.Title)

	content_san := p.Sanitize(c_info.Content)

	if title_san == "" {

		log.Printf("c info comment: title san failed\n")

		c.JSON(http.StatusBadRequest, SERVER_RESP{Status: "error", Reply: "empty title"})

		return

	}

	if content_san == "" {

		log.Printf("c info comment: content san failed\n")

		c.JSON(http.StatusBadRequest, SERVER_RESP{Status: "error", Reply: "empty content"})

		return

	}

	comment_id, _ := pkgutils.GetRandomHex(32)

	err = sendMail(comment_id, title_san, content_san)

	if err != nil {

		log.Printf("c info comment: mail failed: %s\n", err.Error())

		c.JSON(http.StatusInternalServerError, SERVER_RESP{Status: "error", Reply: "server error"})

		return

	}

	err = pkgdb.RegisterComment(comment_id, title_san, content_san, timeRegistered)

	if err != nil {

		log.Printf("c info comment: db register failed: %s\n", err.Error())

		c.JSON(http.StatusInternalServerError, SERVER_RESP{Status: "error", Reply: "server error"})

		return

	}

	reply := fmt.Sprintf("%s", title_san)

	c.JSON(http.StatusOK, SERVER_RESP{Status: "success", Reply: reply})

}

func ApproveComment(c *gin.Context) {

	watchId := c.Param("commentId")

	if !pkgauth.VerifyDefaultValue(watchId) {

		log.Printf("approve comment: illegal: %s\n", watchId)

		c.JSON(http.StatusBadRequest, SERVER_RESP{Status: "error", Reply: "invalid format"})

		return

	}

	now := time.Now().UTC()

	timeApproved := now.Format("2006-01-02-15-04-05")

	err := pkgdb.ApproveComment(watchId, timeApproved)

	if err != nil {

		log.Printf("approve comment: %s\n", err.Error())

		c.JSON(http.StatusInternalServerError, SERVER_RESP{Status: "error", Reply: "error approve"})

		return
	}

	c.JSON(http.StatusOK, SERVER_RESP{Status: "success", Reply: "approved"})

}

func sendMail(commentId string, title string, content string) error {

	var pass = MAIL_JSON.Pass
	var from = "seantywork@gmail.com"
	var to = "seantywork@gmail.com"
	var smtpHost = "smtp.gmail.com"

	callback := pkgglob.G_CONF.Url + "/api/comment/approve/" + commentId

	message := "Subject: our-wedding-rsvp comment approval from [" + title + "]"

	message += "\r\n\r\n"

	message += "title: " + title + "\n"

	message += "content: \n" + content + "\n"

	message += "approval link: \n" + callback + "\n"

	message += "\r\n"

	messageBytes := []byte(message)

	auth := smtp.PlainAuth("", from, pass, smtpHost)
	err := smtp.SendMail(smtpHost+":587", auth, from, []string{to}, messageBytes)
	if err != nil {
		return fmt.Errorf("failed to send mail: %s", err.Error())
	}

	return nil
}
