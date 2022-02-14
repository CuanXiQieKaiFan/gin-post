package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"nonoSheep/middleware"
	"nonoSheep/model"
	"nonoSheep/util/errormsg"
	"strconv"
)

//获取所有帖子列表
func GetPosts(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("size"))
	pageNum, _ := strconv.Atoi(c.Query("page"))
	token := c.Request.Header.Get("Authorization")
	u, _ := middleware.CheckToknen(token)

	data := model.GetPosts(pageSize, pageNum, u.Username)
	code := errormsg.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"message": errormsg.GetErrMsg(code),
	})
}

//获取某一个帖子
func GetPost(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	token := c.Request.Header.Get("Authorization")
	u, _ := middleware.CheckToknen(token)
	data, code := model.GetPost(id, u.Username)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errormsg.GetErrMsg(code),
		"data":    data})
}

//发布帖子
func IssuePost(c *gin.Context) {
	var data model.Post
	tokenHeader := c.Request.Header.Get("Authorization")
	u, _ := middleware.CheckToknen(tokenHeader)
	_ = c.ShouldBind(&data) //从客户端获取信息填入data里面
	postId, code := model.IssuePost(&data, u)
	if data.Title == "" {
		c.JSON(http.StatusOK, gin.H{"message": "参数绑定出错！！"})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"Post_id": postId,
			"message": errormsg.GetErrMsg(code),
		})
	}

}

//更新帖子
func UpdatePost(c *gin.Context) {
	tokenHeader := c.Request.Header.Get("Authorization")
	u, _ := middleware.CheckToknen(tokenHeader) //只能是作者本人才更新帖子
	id, _ := strconv.Atoi(c.Param("id"))
	post := model.OutPutPost(id)
	_ = c.ShouldBind(&post)
	code = model.UpdatePost(id, &post, u)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errormsg.GetErrMsg(code),
	})
}

//删除帖子
func DeletePost(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id")) //从前端返回参数id
	tokenHeader := c.Request.Header.Get("Authorization")
	u, _ := middleware.CheckToknen(tokenHeader) //只能是作者本人才删除帖子
	code := model.DeletePost(id, u)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errormsg.GetErrMsg(code),
	})
}

//搜索帖子
func SearchPost(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("size")) //一页多少个
	pageNum, _ := strconv.Atoi(c.Query("page"))  //第几页
	key := c.Query("key")
	posts, code := model.SearchPost(key, pageSize, pageNum)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errormsg.GetErrMsg(code),
		"posts":   posts,
	})
}
