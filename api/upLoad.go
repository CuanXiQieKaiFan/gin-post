package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"nonoSheep/middleware"
	"path"
)

func UploadFile(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	_, code := middleware.CheckToknen(token)
	if code == 200 {
		f, _ := c.FormFile("file")
		dst := path.Join("./img", f.Filename)
		_ = c.SaveUploadedFile(f, dst)
		url := c.Request.Host + "/img/" + f.Filename
		c.JSON(http.StatusOK, gin.H{
			"message": "files uploaded!",
			"url":     url,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "上传失败！",
		})
	}

}
