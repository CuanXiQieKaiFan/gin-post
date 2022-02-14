package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"nonoSheep/middleware"
	"nonoSheep/model"
	"nonoSheep/util/errormsg"
	"strconv"
)

var code int

//注册用户
func Register(c *gin.Context) {
	var data model.User
	var token string
	var refreshToken string
	_ = c.ShouldBindJSON(&data)           //从客户端获取信息填入data里面
	code = model.CheckUser(data.Username) //检查是否存在
	if code == errormsg.SUCCESS {
		model.Register(&data)
		token, _ = middleware.SetToken(data.Password)
		refreshToken = middleware.RefreshToken(data.Password)
	}
	if code == errormsg.ERROR_USERNAME_USED {
		code = errormsg.ERROR_USERNAME_USED
	}
	c.JSON(http.StatusOK, gin.H{
		"status":       code,
		"message":      errormsg.GetErrMsg(code),
		"token":        token,
		"refreshToken": refreshToken,
	})
}

//刷新token
func RToken(c *gin.Context) {
	refToken := c.Query("refresh_token")
	newToken, newRtoken, code := model.RToken(refToken)
	c.JSON(http.StatusOK, gin.H{
		"status":        code,
		"message":       errormsg.GetErrMsg(code),
		"token":         newToken,
		"refresh_token": newRtoken,
	})
}

//修改密码
func UpdatePwd(c *gin.Context) {
	var data model.User
	tokenHeader := c.Request.Header.Get("Authorization")
	u, _ := middleware.CheckToknen(tokenHeader) //检验token 并通过token获取用户信息

	oldPwd := c.Query("old_password")
	newPwd := model.ScryptPw(c.Query("new_password"))
	data.Username = u.Username
	data.Password = model.ScryptPw(oldPwd)
	code = model.UpdatePwd(&data, newPwd)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errormsg.GetErrMsg(code),
	})
}

//获取单个用户信息
func GetUser(c *gin.Context) {
	var data model.User
	id, _ := strconv.Atoi(c.Param("user_id"))
	code, data = model.GetUser(id)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"message": errormsg.GetErrMsg(code),
	})

}

//更改用户信息
func UpdateUserInfo(c *gin.Context) {
	tokenHeader := c.Request.Header.Get("Authorization")
	u, _ := middleware.CheckToknen(tokenHeader) //检验token 并通过token获取用户信息
	user := model.OutPutUser(u.Username)        //从库中取出的数据，进行数据绑定实现按需更改

	err := c.ShouldBind(&user) //ShouldBind并不会将空值的字段也绑定下来
	fmt.Printf("user:%v\n", user)
	if err != nil {
		fmt.Printf("err is:%v", err)
	}
	code = model.EditUser(&user)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errormsg.GetErrMsg(code),
	})
}
