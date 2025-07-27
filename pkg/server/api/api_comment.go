package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"time"

	pkgauth "lgtm-rsvp/pkg/auth"
	pkgdb "lgtm-rsvp/pkg/db"
	pkgglob "lgtm-rsvp/pkg/glob"
	pkgutils "lgtm-rsvp/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
)

var _comment chan CommentData

type CommentInfo struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type CommntDataList []CommentData

type CommentData struct {
	CommentId string `json:"commentid"`
	Title     string `json:"title"`
	Content   string `json:"content"`
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

	err = pkgdb.RegisterComment(comment_id, title_san, content_san, timeRegistered)

	if err != nil {

		log.Printf("c info comment: db register failed: %s\n", err.Error())

		c.JSON(http.StatusInternalServerError, SERVER_RESP{Status: "error", Reply: "server error"})

		return

	}

	cdata := CommentData{
		CommentId: comment_id,
		Title:     title_san,
		Content:   content_san,
	}

	_comment <- cdata

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

func StartMailer(reterr chan error) {

	_comment = make(chan CommentData)

	reterr <- nil

	for {

		c := <-_comment

		err := sendMail(c.CommentId, c.Title, c.Content)

		if err != nil {

			err = writeMailErr(c)

			if err != nil {
				log.Printf("CRIT: failed to write mail err\n")
				fmt.Println(c)
				log.Printf("------------------------------\n")
			}

			log.Printf("send mail failed: %s\n", c.CommentId)
		} else {
			log.Printf("send mail success\n")
		}

	}

}

func sendMail(commentId string, title string, content string) error {

	var pass = *pkgglob.G_CONF.Api.GoogleComment
	var from = pkgglob.G_CONF.Admin.Id
	var to = pkgglob.G_CONF.Admin.Id
	var smtpHost = "smtp.gmail.com"

	callback := pkgglob.G_CONF.Url + "/api/comment/approve/" + commentId

	message := "Subject: lgtm-rsvp comment approval from [" + title + "]"

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

func writeMailErr(c CommentData) error {

	errCommData := make(CommntDataList, 0)

	if _, err := os.Stat(pkgglob.G_MAIL_ERR_PATH); err == nil {
		f, err := os.ReadFile(pkgglob.G_MAIL_ERR_PATH)

		if err != nil {
			return fmt.Errorf("failed to read mail err path: %v", err)
		}

		err = json.Unmarshal(f, &errCommData)

		if err != nil {
			return fmt.Errorf("failed to unmarshal mail err: %v", err)
		}
	}

	errCommData = append(errCommData, c)

	eb, err := json.Marshal(errCommData)

	if err != nil {
		return fmt.Errorf("failed to marshal mail err: %v", err)
	}

	err = os.WriteFile(pkgglob.G_MAIL_ERR_PATH, eb, 0644)

	if err != nil {
		fmt.Println(string(eb))
		return fmt.Errorf("failed to write mail err: %v", err)
	}

	return nil

}

func CommentSudoCmd(c *gin.Context) {

	if !pkgauth.Is0(c, nil, nil) {

		log.Printf("comment sudo cmd: not admin\n")

		c.JSON(http.StatusForbidden, SERVER_RESP{Status: "error", Reply: "you're not admin"})

		return

	}

	cmd := c.Param("cmd")

	file, err := c.FormFile("file")

	if err != nil {

		log.Printf("comment sudo cmd: file not found: %v\n", err)
		c.JSON(http.StatusForbidden, SERVER_RESP{Status: "error", Reply: "invalid"})
		return
	}

	f, err := file.Open()

	if err != nil {
		log.Printf("comment sudo cmd: file open failed: %v\n", err)
		c.JSON(http.StatusForbidden, SERVER_RESP{Status: "error", Reply: "invalid"})
		return
	}

	var coli CommntDataList

	buf := bytes.NewBuffer(nil)

	if _, err := io.Copy(buf, f); err != nil {
		log.Printf("comment sudo cmd: file read failed: %v\n", err)
		c.JSON(http.StatusForbidden, SERVER_RESP{Status: "error", Reply: "invalid"})
		return
	}

	err = json.Unmarshal(buf.Bytes(), &coli)

	if err != nil {
		log.Printf("comment sudo cmd: unmarshal failed: %v\n", err)
		c.JSON(http.StatusForbidden, SERVER_RESP{Status: "error", Reply: "invalid"})
		return
	}

	reply := ""

	if cmd == "allow" {

		reply = allowLogic(coli)

	} else if cmd == "block" {

		reply = blockLogic(coli)

	} else {
		log.Printf("comment sudo cmd: invalid cmd: %s\n", cmd)
		c.JSON(http.StatusForbidden, SERVER_RESP{Status: "error", Reply: "invalid"})
		return
	}

	log.Printf("comment sudo cmd: success\n")
	c.JSON(http.StatusOK, SERVER_RESP{Status: "success", Reply: reply})
	return

}

func allowLogic(coli CommntDataList) string {

	var retstr = ""

	clen := len(coli)

	log.Printf("allowing all data in comment list...\n")

	for i := 0; i < clen; i++ {

		c, err := pkgdb.GetCommentById(coli[i].CommentId)

		id := ""
		title := ""

		if err != nil {
			retstr += fmt.Sprintf("  - comment by id doesn't exit: %v\n", err)
			comment_id, _ := pkgutils.GetRandomHex(32)
			now := time.Now().UTC()
			p := bluemonday.UGCPolicy()
			title_san := p.Sanitize(coli[i].Title)
			content_san := p.Sanitize(coli[i].Content)
			if title_san == "" {
				retstr += fmt.Sprintf("  - register comment: invalid title\n")
				continue
			}
			if content_san == "" {
				retstr += fmt.Sprintf("  - register comment: invalid content\n")
				continue
			}
			timeRegistered := now.Format("2006-01-02-15-04-05")
			err = pkgdb.RegisterComment(comment_id, title_san, content_san, timeRegistered)
			if err != nil {
				retstr += fmt.Sprintf("  - register comment failed: %v\n", err)
			} else {
				retstr += fmt.Sprintf("  - register comment success: %s\n", title_san)
			}
			id = comment_id
			title = title_san
		} else {
			id = c.Id
			title = c.Title
		}

		now := time.Now().UTC()
		timeApproved := now.Format("2006-01-02-15-04-05")

		err = pkgdb.ApproveComment(id, timeApproved)

		if err != nil {
			retstr += fmt.Sprintf("  - approve comment failed: %v\n", err)
		} else {
			retstr += fmt.Sprintf("  - approve comment success: %s\n", title)
		}
	}
	log.Printf("allow done\n")
	fmt.Println(retstr)
	log.Printf("==========\n")

	return retstr
}

func blockLogic(coli CommntDataList) string {

	var retstr = ""

	clen := len(coli)

	log.Printf("blocking all data in comment list...\n")

	for i := 0; i < clen; i++ {
		err := pkgdb.DisapproveCommentByTitle(coli[i].Title)
		if err != nil {
			retstr += fmt.Sprintf("  - disapprove by title failed: %v\n", err)
		} else {
			retstr += fmt.Sprintf("  - disapprove by title success: %s\n", coli[i].Title)
		}

	}
	log.Printf("blocking done\n")
	fmt.Println(retstr)
	log.Printf("==========\n")
	return retstr
}
