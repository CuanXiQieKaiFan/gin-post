package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"nonoSheep/model"
	"nonoSheep/util/errormsg"
)

//获取所有主题
func GetTopics(c *gin.Context) {
	topics, code := model.GetTopics()
	c.JSON(http.StatusOK,gin.H{
		"status":code,
		"message":errormsg.GetErrMsg(code),
		"topics":topics,
	})
}
//创建主题
func CreateTopic(c *gin.Context) {
	var data model.Topic
	_ = c.ShouldBindJSON(&data)
	code:=model.CreatTopic(&data)
	c.JSON(http.StatusOK,gin.H{
		"status":code,
		"message":errormsg.GetErrMsg(code),
	})
}
