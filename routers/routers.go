package routers

import (
	"github.com/gin-gonic/gin"
	"nonoSheep/api"
	"nonoSheep/middleware"
)

func Start() {
	r := gin.Default()
	r.Static("img", "./img")
	r.Use(middleware.Cors()) //跨域

	user := r.Group("/user") //用户相关操作
	{
		user.GET("/token/refresh", api.RToken)  //刷新token
		user.POST("/register", api.Register)    //注册用户
		user.GET("/token", api.Login)           //登录用户
		user.PUT("/password", api.UpdatePwd)    //修改密码
		user.PUT("/info", api.UpdateUserInfo)   //更改用户信息
		user.GET("/info/:user_id", api.GetUser) //获取用户信息
	}
	post := r.Group("/post") //帖子相关操作
	{
		post.GET("/list", api.GetPosts)            //获取帖子列表
		post.GET("/single/:id", api.GetPost)       //获取某一个帖子
		post.POST("/single", api.IssuePost)        //发布帖子
		post.PUT("/single/:id", api.UpdatePost)    //更新帖子
		post.DELETE("/single/:id", api.DeletePost) //删除帖子
		post.GET("/search", api.SearchPost)        //搜索帖子  分页 模糊
	}
	comment := r.Group("/") //评论相关操作
	{
		comment.POST("/comment", api.IssueCommnet)        //发布评论
		comment.PUT("/comment/:id", api.UpdateCommnet)    //更新评论
		comment.DELETE("/comment/:id", api.DeleteComment) //删除评论
		comment.GET("/comment", api.GetComment)           //获取某个帖子下的评论
	}
	r.POST("topic/create", api.CreateTopic) //创建主题
	r.GET("/topic/list", api.GetTopics)     //获取所有主题

	others := r.Group("/operate") //其他操作 收藏点赞
	{
		others.PUT("/praise", api.Praise)              //点赞
		others.GET("/collect/list", api.GetCollection) //获取用户收藏列表
		others.PUT("/collect", api.Collect)            //收藏
		others.GET("/focus/list", api.GetFocus)        //获取关注列表
		others.PUT("/focus", api.Focus)                //关注
		others.GET("/praise/toyou", api.PraForU)       //获取为我点赞的用户列表
		others.GET("/view/list", api.GetViewList)      //获取浏览记录
		others.GET("/focus/onyou", api.FocusForU)      //获取关注我的用户列表
	}

	_ = r.Run(":9090")
}
