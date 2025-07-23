package db

import (
	"fmt"
	"mime/multipart"
	"os"

	pkgglob "lgtm-rsvp/pkg/glob"

	"github.com/gin-gonic/gin"
)

var mediaPath = pkgglob.G_MEDIA_PATH

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
