package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var mainImage ImageInfo

var imageList []ImageInfo

type ImageInfo struct {
	Name string `json:"name"`
}

func AddImageList(paths []string) {

	imageList = make([]ImageInfo, 0)

	for i := 0; i < len(paths); i++ {

		if i == 0 {
			mainImage = ImageInfo{
				Name: paths[i],
			}
		}

		imageList = append(imageList, ImageInfo{
			Name: paths[i],
		})

	}

}

func GetMainImage() ImageInfo {
	return mainImage
}

func GetImageList(c *gin.Context) {

	ibytes, err := json.Marshal(imageList)

	if err != nil {
		log.Printf("image list: marshal: %s\n", err.Error())

		c.JSON(http.StatusBadRequest, SERVER_RESP{Status: "error", Reply: "invalid format"})

		return
	}

	c.JSON(http.StatusOK, SERVER_RESP{Status: "success", Reply: string(ibytes)})

}
