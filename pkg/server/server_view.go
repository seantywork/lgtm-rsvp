package server

import (
	"github.com/gin-gonic/gin"
)

func getIndex(c *gin.Context) {

	c.HTML(200, "index.html", gin.H{})

}
