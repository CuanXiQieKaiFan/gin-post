package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"nonoSheep/middleware"
	"nonoSheep/model"
	"nonoSheep/util/errormsg"
	"strconv"
)

//点赞
func Praise(c *gin.Context) {
	mod, _ := strconv.Atoi(c.PostForm("model"))
	targetId, _ := strconv.Atoi(c.PostForm("target_id"))
	token := c.Request.Header.Get("Authorization")
	u, _ := middleware.CheckToknen(token)
	code := model.Praise(mod, targetId, u.Username)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errormsg.GetErrMsg(code),
	})

}

//收藏
func Collect(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	u, _ := middleware.CheckToknen(token)
	id, _ := strconv.Atoi(c.PostForm("post_id"))
	code := model.Collect(id, u.Username)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errormsg.GetErrMsg(code),
	})

}

//关注
func Focus(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	u, _ := middleware.CheckToknen(token)
	id, _ := strconv.Atoi(c.PostForm("user_id"))
	code := model.Focus(id, u.Username)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errormsg.GetErrMsg(code),
	})
}

//获取用户收藏列表
func GetCollection(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	u, _ := middleware.CheckToknen(token)
	posts, code := model.GetCollection(u.Username)
	c.JSON(http.StatusOK, gin.H{
		"status":      code,
		"message":     errormsg.GetErrMsg(code),
		"collections": posts,
	})
}

//获取关注列表
func GetFocus(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	u, _ := middleware.CheckToknen(token)
	users, code := model.GetFocus(u.Username)
	c.JSON(http.StatusOK, gin.H{
		"status":      code,
		"message":     errormsg.GetErrMsg(code),
		"collections": users,
	})
}

//获取给你点赞的用户列表
func PraForU(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	u, _ := middleware.CheckToknen(token)
	postform := c.PostForm("model")
	mod, _ := strconv.Atoi(postform)
	users, code := model.PraForU(u.Username, mod)
	c.JSON(http.StatusOK, gin.H{
		"status":     code,
		"message":    errormsg.GetErrMsg(code),
		"praiseList": users,
	})
}

//获取浏览记录
func GetViewList(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	u, _ := middleware.CheckToknen(token)
	posts, code := model.GetViewList(u.Username)
	c.JSON(http.StatusOK, gin.H{
		"status":   code,
		"message":  errormsg.GetErrMsg(code),
		"viewList": posts,
	})
}

//获取关注你的用户列表
func FocusForU(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	u, _ := middleware.CheckToknen(token)
	userList, code := model.FocusForU(u.Username)
	c.JSON(http.StatusOK, gin.H{
		"status":    code,
		"message":   errormsg.GetErrMsg(code),
		"focusList": userList,
	})
}
