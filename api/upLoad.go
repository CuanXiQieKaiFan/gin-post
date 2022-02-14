package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
)

func UploadFile(c *gin.Context) {
	f, _ := c.FormFile("file")
	dst:=filepath.Base(f.Filename)
	_=c.SaveUploadedFile(f, dst)
	url:=c.Request.Host+"/img/"+f.Filename
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("%d ","files uploaded!"),
		"url":url,
	})
}
