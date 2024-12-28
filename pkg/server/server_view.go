package server

import (
	"github.com/gin-gonic/gin"
)

func getIndex(c *gin.Context) {

	c.HTML(200, "index.html", gin.H{})

}

func getSignin(c *gin.Context) {

	c.HTML(200, "signin.html", gin.H{})
}

func getRead(c *gin.Context) {

	c.HTML(200, "read.html", gin.H{})
}

func getWrite(c *gin.Context) {

	c.HTML(200, "write.html", gin.H{})
}
