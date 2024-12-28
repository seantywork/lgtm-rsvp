package db

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

var mediaPath = "./data/media/"

func UploadMedia(c *gin.Context, file *multipart.FileHeader, new_filename string) error {

	this_file_path := mediaPath + new_filename

	err := c.SaveUploadedFile(file, this_file_path)

	if err != nil {

		return fmt.Errorf("failed to upload: %s", err.Error())
	}

	return nil
}

func DownloadMedia(c *gin.Context, watchId string) error {

	this_media_path := mediaPath + watchId

	if _, err := os.Stat(this_media_path); err != nil {

		return err

	}

	c.File(this_media_path)

	return nil
}

func DeleteMedia(media_key string) error {

	this_media_path := mediaPath + media_key

	err := os.Remove(this_media_path)

	if err != nil {

		return fmt.Errorf("failed to delete media: rm: %s", err.Error())
	}

	return nil
}

func GetAssociateMediaKeysForEditorjsSrc(rawArticle []byte) ([]string, error) {

	var retlist []string

	var editorjsSrc map[string]interface{}

	err := json.Unmarshal(rawArticle, &editorjsSrc)

	if err != nil {

		return nil, fmt.Errorf("failed to unmarshal: %s", err.Error())

	}

	blocks, okay := editorjsSrc["blocks"]

	if !okay {

		return nil, fmt.Errorf("invalid format: %s", "no blocks")
	}

	blocksList := blocks.([]interface{})

	blocksLen := len(blocksList)

	for i := 0; i < blocksLen; i++ {

		blockObj := blocksList[i].(map[string]interface{})

		objType, okay := blockObj["type"]

		if !okay {
			continue
		}

		if objType != "image" {
			continue
		}

		objData, okay := blockObj["data"]

		if !okay {
			continue
		}

		objFields := objData.(map[string]interface{})

		fileField, okay := objFields["file"]

		if !okay {
			continue
		}

		targetProps := fileField.(map[string]interface{})

		urlTarget, okay := targetProps["url"]

		if !okay {
			continue
		}

		target := urlTarget.(string)

		pathList := strings.Split(target, "/")

		keyExt := pathList[len(pathList)-1]

		keyExtList := strings.Split(keyExt, ".")

		key := keyExtList[0]

		retlist = append(retlist, key)
	}

	return retlist, nil
}
