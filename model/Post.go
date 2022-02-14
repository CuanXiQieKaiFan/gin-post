package model

import (
	"fmt"
	"nonoSheep/middleware"
	"nonoSheep/util/errormsg"
	"strconv"
	"strings"
)

type Post struct {
	Id           int    `json:"id"`
	IsPraised    bool   `json:"isPraised"`
	IsCollect    bool   `json:"isCollect"`
	CommentCount int    `json:"comment_count"`
	PraiseCount  int    `json:"praiseCount"`
	PublishTime  string `json:"publishTime"`
	Title        string `json:"title"`
	TopicId      int    `json:"topicId"`
	Content      string `json:"content"`
	Pictures     string `json:"pictures"`
	UserId       int    `json:"userId"`
	NickName     string `json:"nickName"`
	Avatar       string `json:"avatar"`
}

//获取帖子列表
func GetPosts(pageSize int, pageNum int, username string) []Post {
	var posts = make([]Post, pageSize) //不可以直接 var posts []Post 这样posts是nil
	num := (pageNum - 1) * pageSize
	rows, err := Db.Query("select id,publishTime,content,pictures,topicId,userId,nickName,praiseCount,commentCount,title from posts limit ?,?", num, pageSize)
	//如 pageNum=1，则会返回第一页的数据，=2返回第二页的数据
	if err != nil {
		return nil
	}
	//取出点赞,收藏列表
	var postId, focus string
	row := Db.QueryRow("select praisePost,focus from users where username=?", username)
	_ = row.Scan(&postId, &focus)
	praisePosts := strings.SplitN(postId, " ", -1) //post的id
	focusId := strings.SplitN(focus, " ", -1)      //post的id

	var i = 0
	for rows.Next() {
		_ = rows.Scan(&posts[i].Id, &posts[i].PublishTime, &posts[i].Content, &posts[i].Pictures, &posts[i].TopicId, &posts[i].UserId, &posts[i].NickName, &posts[i].PraiseCount, &posts[i].CommentCount, &posts[i].Title)
		for j := 0; j < len(praisePosts)-1; j++ {
			pid, _ := strconv.Atoi(praisePosts[j])
			if posts[i].Id == pid {
				posts[i].IsPraised = true
			}
		}
		for j := 0; j < len(focusId)-1; j++ {
			fid, _ := strconv.Atoi(focusId[j])
			if posts[i].Id == fid {
				posts[i].IsCollect = true
			}
		}
		i++
	}
	return posts
}

//获取某一个帖子
func GetPost(id int, username string) (Post, int) {
	var post Post
	row := Db.QueryRow("select id,publishTime,content,pictures,topicId,userId,nickName,praiseCount,commentCount,title from posts where id =?", id)
	row.Scan(&post.Id, &post.PublishTime, &post.Content, &post.Pictures, &post.TopicId, &post.UserId, &post.NickName, &post.PraiseCount, &post.CommentCount, &post.Title)
	if err != nil {
		return post, errormsg.ERROR_POST_NOTEXIST
	}

	//取出点赞,收藏列表
	var postId, focus string
	r := Db.QueryRow("select praisePost,focus from users where username=?", username)
	_ = r.Scan(&postId, &focus)
	praisePosts := strings.SplitN(postId, " ", -1) //post的id
	focusId := strings.SplitN(focus, " ", -1)      //post的id
	for j := 0; j < len(praisePosts)-1; j++ {
		pid, _ := strconv.Atoi(praisePosts[j])
		if post.Id == pid {
			post.IsPraised = true
		}
	}
	for j := 0; j < len(focusId)-1; j++ {
		fid, _ := strconv.Atoi(focusId[j])
		if post.Id == fid {
			post.IsCollect = true
		}
	}
	return post, errormsg.SUCCESS
}

//发布帖子
func IssuePost(data *Post, u *middleware.MyClaim) (int, int) {
	err1 := Db.QueryRow("select id,nickName from users where username=?", u.Username).Scan(&data.UserId, &data.NickName)
	if err1 != nil {
		fmt.Printf("err1 is:%v\n", err1)
	}
	result, err := Db.Exec("insert into posts (title,content,topicId,pictures,userId,nickName) values (?,?,?,?,?,?)", data.Title, data.Content, data.TopicId, data.Pictures, data.UserId, data.NickName)
	if err != nil {
		fmt.Printf("err is:", err1)
		return 0, errormsg.ERROR
	}
	postId, _ := result.LastInsertId()
	return int(postId), errormsg.SUCCESS
}

//更新帖子
func UpdatePost(id int, data *Post, u *middleware.MyClaim) int {
	_ = Db.QueryRow("select id from users where username =?", u.Username).Scan(&data.UserId) //取到作者id
	_, err := Db.Exec("update posts set title=?,content=?,topicId=?,pictures=? where id = ?", data.Title, data.Content, data.TopicId, data.Pictures, id)
	if err != nil {
		fmt.Printf("err is:%v\n", err)
		return errormsg.ERROR
	}
	return errormsg.SUCCESS
}

//删除帖子
func DeletePost(id int, u *middleware.MyClaim) int {
	var userId int
	_ = Db.QueryRow("select id from users where username=?", u.Username).Scan(&userId)
	_, err := Db.Exec("delete from posts where id=? and userId =?", id, userId)
	if err != nil {
		fmt.Printf("err is:%v\n", err)
		return errormsg.ERROR
	}
	return errormsg.SUCCESS
}

//搜索帖子
func SearchPost(key string, pageSize int, pageNum int) ([]Post, int) {
	var posts = make([]Post, pageSize) //模糊查询使用concat 连接%和字段 concat('%',?,'%')
	num := (pageNum - 1) * pageSize
	rows, err := Db.Query("select id,publishTime,content,pictures,topicId,userId,nickName,praiseCount,commentCount,title from posts where content like concat('%',?,'%') limit ?,?", key, num, pageSize)
	if err != nil || rows.Err() != nil {
		fmt.Printf("err is:%v\n", err)
		return posts, errormsg.ERROR_POST_NOTEXIST
	}
	var i int = 0
	for rows.Next() {
		_ = rows.Scan(&posts[i].Id, &posts[i].PublishTime, &posts[i].Content, &posts[i].Pictures, &posts[i].TopicId, &posts[i].UserId, &posts[i].NickName, &posts[i].PraiseCount, &posts[i].CommentCount, &posts[i].Title)
		i++
	}
	return posts, errormsg.SUCCESS
}

//取出帖子
func OutPutPost(id int) Post {
	var post Post
	row := Db.QueryRow("select id,publishTime,content,pictures,topicId,userId,nickName,praiseCount,commentCount,title from posts where id =?", id)
	_ = row.Scan(&post.Id, &post.PublishTime, &post.Content, &post.Pictures, &post.TopicId, &post.UserId, &post.NickName, &post.PraiseCount, &post.CommentCount, &post.Title)
	return post
}
