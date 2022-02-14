package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"nonoSheep/middleware"
	"nonoSheep/model"
	"nonoSheep/util/errormsg"
)

//刷新token
func RefreshToken(c *gin.Context) {
	var rToken string
	c.ShouldBind(&rToken)
	token, refToken, code := model.RToken(rToken)
	c.JSON(http.StatusOK, gin.H{
		"status":       code,
		"message":      errormsg.GetErrMsg(code),
		"token":        token,
		"refreshToken": refToken,
	})
}

//登录
func Login(c *gin.Context) {
	uname := c.Query("username")
	upwd := c.Query("password")
	var token string
	var refreshToken string
	var code int
	code = model.CheckLogin(uname, upwd)

	if code == errormsg.SUCCESS {
		token, code = middleware.SetToken(uname)
		refreshToken = middleware.RefreshToken(uname)
	}
	c.JSON(http.StatusOK, gin.H{
		"status":       code,
		"message":      errormsg.GetErrMsg(code),
		"token":        token,
		"refreshToken": refreshToken,
	})
}
