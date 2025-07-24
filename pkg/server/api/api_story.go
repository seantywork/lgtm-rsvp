package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	pkgauth "lgtm-rsvp/pkg/auth"
	pkgdb "lgtm-rsvp/pkg/db"
	pkgutils "lgtm-rsvp/pkg/utils"
)

type ArticleInfo struct {
	Id               string `json:"id,omitempty"`
	Title            string `json:"title"`
	Content          string `json:"content,omitempty"`
	Intro            string `json:"intro"`
	DateMarked       string `json:"dateMarked"`
	PrimaryMediaName string `json:"primaryMediaName"`
}

func UploadStory(c *gin.Context) {

	if !pkgauth.Is0(c, nil, nil) {

		log.Printf("article upload: not admin\n")

		c.JSON(http.StatusForbidden, SERVER_RESP{Status: "error", Reply: "you're not admin"})

		return

	}

	var req CLIENT_REQ

	var a_info ArticleInfo

	if err := c.BindJSON(&req); err != nil {

		log.Printf("article upload: failed to bind: %s\n", err.Error())

		c.JSON(http.StatusBadRequest, SERVER_RESP{Status: "error", Reply: "invalid format"})

		return
	}

	err := json.Unmarshal([]byte(req.Data), &a_info)

	if err != nil {

		log.Printf("article upload: failed to unmarshal: %s\n", err.Error())

		c.JSON(http.StatusBadRequest, SERVER_RESP{Status: "error", Reply: "invalid format"})

		return

	}

	if a_info.Title == "" {

		log.Printf("article upload: needs title\n")

		c.JSON(http.StatusBadRequest, SERVER_RESP{Status: "error", Reply: "invalid format"})

		return

	}

	if a_info.Content == "" {

		log.Printf("article upload: needs content\n")

		c.JSON(http.StatusBadRequest, SERVER_RESP{Status: "error", Reply: "invalid format"})

		return

	}

	if a_info.Intro == "" {

		log.Printf("article: upload: needs intro\n")

		c.JSON(http.StatusBadRequest, SERVER_RESP{Status: "error", Reply: "invalid format"})

		return
	}

	if a_info.DateMarked == "" {

		log.Printf("article upload: needs dateMarked\n")

		c.JSON(http.StatusBadRequest, SERVER_RESP{Status: "error", Reply: "invalid format"})

		return

	}

	if a_info.PrimaryMediaName == "" {

		log.Printf("article upload: needs primaryMediaName\n")

		c.JSON(http.StatusBadRequest, SERVER_RESP{Status: "error", Reply: "invalid format"})

		return

	}

	new_file_name, _ := pkgutils.GetRandomHex(32)

	err = pkgdb.SaveStory(new_file_name, a_info.Title, a_info.Intro, a_info.DateMarked, a_info.PrimaryMediaName, a_info.Content)

	if err != nil {

		log.Printf("article upload: failed to upload: %s", err.Error())

		c.JSON(http.StatusInternalServerError, SERVER_RESP{Status: "error", Reply: "failed to upload"})

		return
	}

	c.JSON(http.StatusOK, SERVER_RESP{Status: "success", Reply: "uploaded"})

}

func DownloadStoryById(c *gin.Context) {

	watchId := c.Param("storyId")

	if !pkgauth.VerifyDefaultValue(watchId) {

		log.Printf("get article: illegal: %s\n", watchId)

		c.JSON(http.StatusBadRequest, SERVER_RESP{Status: "error", Reply: "invalid format"})

		return

	}

	story, err := pkgdb.GetStoryById(watchId)

	if err != nil {

		log.Printf("failed to get article: %s\n", err.Error())

		c.JSON(http.StatusBadRequest, SERVER_RESP{Status: "error", Reply: "invalid format"})

		return
	}

	content := story.Content

	c.JSON(http.StatusOK, SERVER_RESP{Status: "success", Reply: content})

}

func UploadStoryMedia(c *gin.Context) {

	if !pkgauth.Is0(c, nil, nil) {

		log.Printf("media upload: not admin\n")

		c.JSON(http.StatusForbidden, SERVER_RESP{Status: "error", Reply: "you're not admin"})

		return

	}

	file, err := c.FormFile("file")

	if err != nil {

		log.Printf("media upload: file not found: %v\n", err)

		c.JSON(http.StatusForbidden, SERVER_RESP{Status: "error", Reply: "invalid"})

		return

	}

	rawMediaType := file.Header.Get("Content-Type")

	mediaProplist := strings.Split(rawMediaType, "/")

	//mediaType := mediaProplist[0]
	mediaExt := mediaProplist[1]

	f_name := file.Filename

	f_name_list := strings.Split(f_name, ".")

	f_name_len := len(f_name_list)

	if f_name_len < 1 {

		log.Println("no extension specified")

		c.JSON(http.StatusBadRequest, SERVER_RESP{Status: "error", Reply: "invalid format"})

		return
	}

	extension := f_name_list[f_name_len-1]

	log.Printf("received: %s, size: %d, type: %s, ext: %s\n", file.Filename, file.Size, rawMediaType, extension)

	file_name, _ := pkgutils.GetRandomHex(32)

	file_name = file_name + "." + mediaExt

	err = pkgdb.UploadMedia(c, file, file_name)

	if err != nil {

		log.Println(err.Error())

		c.JSON(http.StatusInternalServerError, SERVER_RESP{Status: "error", Reply: "failed to save"})

		return

	}

	client_file_name := file_name

	c.JSON(http.StatusOK, SERVER_RESP{Status: "success", Reply: client_file_name})

}

func DownloadStoryMediaById(c *gin.Context) {

	watchId := c.Param("mediaId")

	if !pkgauth.VerifyMediaKey(watchId) {

		log.Printf("download media: illegal: %s\n", watchId)

		c.JSON(http.StatusBadRequest, SERVER_RESP{Status: "error", Reply: "invalid format"})

		return

	}

	err := pkgdb.DownloadMedia(c, watchId)

	if err != nil {

		log.Printf("download media: %s\n", err.Error())

		c.JSON(http.StatusBadRequest, SERVER_RESP{Status: "error", Reply: "invalid format"})

		return

	}

	log.Println("media download success")
}

func GetStoryList(c *gin.Context) {

	articles := make([]ArticleInfo, 0)

	stories, err := pkgdb.GetAllStory()

	if err != nil {

		log.Printf("all stories: %s\n", err.Error())

		c.JSON(http.StatusBadRequest, SERVER_RESP{Status: "error", Reply: "invalid format"})

		return

	}

	slen := len(stories)

	for i := 0; i < slen; i++ {

		article := ArticleInfo{
			Id:               stories[i].Id,
			Title:            stories[i].Title,
			Intro:            stories[i].Intro,
			DateMarked:       stories[i].DateMarked,
			PrimaryMediaName: stories[i].PrimaryMediaName,
		}

		articles = append(articles, article)

	}

	abytes, err := json.Marshal(articles)

	if err != nil {
		log.Printf("all stories: marshal: %s\n", err.Error())

		c.JSON(http.StatusBadRequest, SERVER_RESP{Status: "error", Reply: "invalid format"})

		return
	}

	c.JSON(http.StatusOK, SERVER_RESP{Status: "success", Reply: string(abytes)})

}
