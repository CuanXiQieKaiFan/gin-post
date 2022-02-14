package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"nonoSheep/middleware"
	"nonoSheep/model"
	"nonoSheep/util/errormsg"
	"strconv"
)

//获取帖子下的评论
func GetComment(c *gin.Context) {
	mod, _ := strconv.Atoi(c.Query("model"))
	targId, _ := strconv.Atoi(c.Query("target_id"))
	size, _ := strconv.Atoi(c.Query("size"))
	page, _ := strconv.Atoi(c.Query("page"))
	comments, code := model.GetComment(mod, targId, size, page)
	c.JSON(http.StatusOK, gin.H{
		"status":   code,
		"message":  errormsg.GetErrMsg(code),
		"comments": comments,
	})
}

//发布评论
func IssueCommnet(c *gin.Context) {
	mod, _ := strconv.Atoi(c.PostForm("model"))
	targId, _ := strconv.Atoi(c.PostForm("target_id"))
	content := c.PostForm("content")
	pictures := c.PostForm("photo")
	tokenHeader := c.Request.Header.Get("Authorization")
	u, _ := middleware.CheckToknen(tokenHeader)

	var comment model.Comment
	comment.Content = content
	comment.Pictures = pictures
	comId, code := model.IssueCom(mod, targId, &comment, u.Username)
	c.JSON(http.StatusOK, gin.H{
		"status":    code,
		"message":   errormsg.GetErrMsg(code),
		"commentId": comId,
	})
}

//更新评论
func UpdateCommnet(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	content := c.PostForm("content")
	photo := c.PostForm("photo")
	tokenHeader := c.Request.Header.Get("Authorization")
	_, code1 := middleware.CheckToknen(tokenHeader)
	if code1 != 200 {
		c.JSON(http.StatusOK, gin.H{
			"status":  code1,
			"message": errormsg.GetErrMsg(code1),
		})
		c.Abort()
	}

	var newCom model.Comment
	newCom.Content = content
	newCom.Pictures = photo

	code := model.UpdateComment(id, &newCom)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errormsg.GetErrMsg(code),
	})
}

//删除评论
func DeleteComment(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	tokenHeader := c.Request.Header.Get("Authorization") //检验token
	_, code1 := middleware.CheckToknen(tokenHeader)
	if code1 != 200 {
		c.JSON(http.StatusOK, gin.H{
			"status":  code1,
			"message": errormsg.GetErrMsg(code1),
		})
		c.Abort()
	}

	code := model.DeleteComment(id)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errormsg.GetErrMsg(code),
	})
}
