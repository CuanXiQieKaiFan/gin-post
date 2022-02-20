package model

import (
	"database/sql"
	"fmt"
	"nonoSheep/util/errormsg"
	"strconv"
	"strings"
)

//点赞
func Praise(mod int, tag int, username string) int {
	if mod == 1 { //为帖子点赞
		var pid int
		_ = Db.QueryRow("select id from posts where id = ? ", tag).Scan(&pid)
		if pid == tag {
			var praPost string                                                                        //获取你点赞的人
			var postUserId int                                                                        //帖子作者id
			var uid int                                                                               //点赞者id
			var praPostFU string                                                                      //获取为你点赞的人(该帖子的作者)
			_ = Db.QueryRow("select praisePost from users where username=?", username).Scan(&praPost) //取出原来的点赞记录列表
			_, err := Db.Exec("update users set praisePost=? where username=?", praPost+strconv.Itoa(tag)+" ", username)
			_ = Db.QueryRow("select userId from posts where id=?", tag).Scan(&postUserId)               //查出该帖子的作者
			_ = Db.QueryRow("select id from users where username=?", username).Scan(&uid)               //查出点赞者id
			_ = Db.QueryRow("select praisePostForU from users where id=?", postUserId).Scan(&praPostFU) //取出该作者被点赞列表
			_, _ = Db.Exec("update users set praisePostForU =? where id = ?", praPostFU+strconv.Itoa(uid)+" ", postUserId)
			if err != nil {
				fmt.Printf("err is:%v", err)
				return errormsg.ERROR_PRAISE
			} else {
				//帖子点赞数+1
				var praiseCount int
				_ = Db.QueryRow("select praiseCount from posts where id=?	", tag).Scan(&praiseCount)
				Db.Exec("update posts set praiseCount=? where id=?", praiseCount+1, tag)
				return errormsg.SUCCESS
			}
		} else {
			return errormsg.ERROR_POST_NOTEXIST
		}

	}
	if mod == 2 { //为评论点赞
		var cid int
		_ = Db.QueryRow("select id from posts where id = ? ", tag).Scan(&cid)
		if cid == tag {
			var praCom string
			var comUserId int
			var uid int
			var praComFU string
			_ = Db.QueryRow("select praiseCom from users where username=?", username).Scan(&praCom)
			_, err := Db.Exec("update users set praiseCom=? where username=?", praCom+strconv.Itoa(tag)+" ", username)
			_ = Db.QueryRow("select userId from comments where id=?", tag).Scan(&comUserId)          //查出该评论的作者
			_ = Db.QueryRow("select id from users where username=?", username).Scan(&uid)            //查出点赞者id
			_ = Db.QueryRow("select praiseComForU from users where id=?", comUserId).Scan(&praComFU) //取出该作者被点赞列表
			_, _ = Db.Exec("update users set praisePostForU =? where id = ?", praComFU+strconv.Itoa(uid)+" ", comUserId)
			if err != nil {
				return errormsg.ERROR_PRAISE
			} else {
				//评论点赞数+1
				var praiseCount int
				_ = Db.QueryRow("select praiseCount from comments where id=?	", tag).Scan(&praiseCount)
				Db.Exec("update comments set praiseCount=? where id=?", praiseCount+1, tag)
				return errormsg.SUCCESS
			}
		} else {
			return errormsg.ERROR_COMMENT_NOTEXIST
		}

	}

	return errormsg.ERROR
}

//收藏
func Collect(targId int, username string) int {
	var collection string
	_ = Db.QueryRow("select collection from users where username=?", username).Scan(&collection)
	_, err = Db.Exec("update users set collection=? where username = ?", collection+strconv.Itoa(targId)+" ", username)
	if err != nil {
		return errormsg.ERROR_COLLECT
	} else {
		return errormsg.SUCCESS
	}
}

//关注
func Focus(targId int, username string) int {
	var focus string
	var uid int          //当前用户id
	var focusList string //被关注用户的被关注列表
	_ = Db.QueryRow("select focus from users where username=?", username).Scan(&focus)
	_, err = Db.Exec("update users set focus=? where username = ?", focus+strconv.Itoa(targId)+" ", username)
	//将自己id写入关注的用户中
	_ = Db.QueryRow("select id from users where username =?", username).Scan(&uid)
	_ = Db.QueryRow("select focusForU from users where id = ?	", targId).Scan(&focusList)
	_, _ = Db.Exec("update users set focusForU=? where id =?", focusList+strconv.Itoa(uid)+" ", targId)
	if err != nil {
		return errormsg.ERROR_FOCUS
	}
	return errormsg.SUCCESS
}

//获取收藏列表哦
func GetCollection(username string) ([]map[string]interface{}, int) {
	var str string
	row := Db.QueryRow("select collection from users where username=?", username)
	_ = row.Scan(&str)
	postIdSlice := strings.SplitN(str, " ", -1)
	if postIdSlice == nil {
		return nil, errormsg.ERROR
	}
	var pid int
	var posts = make([]map[string]interface{}, len(postIdSlice)-1) //创建一个map类型的切片   len-1去除字段最后一位空格
	fmt.Printf("len of postSlic is:%v\n", len(postIdSlice))
	var post Post //创建一个User类型接受数据
	for i, id := range postIdSlice {
		if id == "" {
			break
		}
		fmt.Printf("id%vis:%v\n", i, id)
		posts[i] = make(map[string]interface{}) //为每一个切片创建map
		pid, _ = strconv.Atoi(id)
		_ = Db.QueryRow("select id,title,publishTime,userId,avatar,nickName from posts where id=?", pid).Scan(&post.Id, &post.Title, &post.PublishTime, &post.UserId, &post.Avatar, &post.NickName)
		posts[i] = map[string]interface{}{
			"post_id":      post.Id,
			"title":        post.Title,
			"publish_time": post.PublishTime,
			"user_id":      post.UserId,
			"avatar":       post.Avatar,
			"nickName":     post.NickName,
		}
	}
	return posts, errormsg.SUCCESS
}

//获取用户关注列表
func GetFocus(username string) ([]map[string]interface{}, int) {
	var str string
	row := Db.QueryRow("select focus from users where username=?", username)
	_ = row.Scan(&str)
	userIdSlice := strings.SplitN(str, " ", -1) //以空格为分割符将focus字符串中的值取出
	if userIdSlice == nil {
		return nil, errormsg.ERROR
	}
	var uid int
	var users = make([]map[string]interface{}, len(userIdSlice)-1) //创建一个map类型的切片   len-1去除focus字段最后一位空格
	var user User                                                  //创建一个User类型接受数据
	for i, id := range userIdSlice {
		if id == "" {
			break
		}
		fmt.Printf("id%vis:%v\n", i, id)
		users[i] = make(map[string]interface{}) //为每一个切片创建map
		uid, _ = strconv.Atoi(id)
		_ = Db.QueryRow("select id,avatar,nickName,introduction from users where id=?", uid).Scan(&user.ID, &user.Avatar, &user.NickName, &user.Introduction)
		users[i] = map[string]interface{}{
			"user_id":      user.ID,
			"avatar":       user.Avatar,
			"nickName":     user.NickName,
			"introduction": user.Introduction,
		}
	}
	return users, errormsg.SUCCESS
}

//获取给你点赞的用户列表
func PraForU(username string, model int) ([]map[string]interface{}, int) {
	//获取帖子点赞用户
	var str string
	var row *sql.Row
	if model == 1 {
		row = Db.QueryRow("select praisePostForU from users where username=?", username)
	}
	//评论...
	if model == 2 {
		row = Db.QueryRow("select praiseComForU from users where username=?", username)
	}
	if row != nil {
		_ = row.Scan(&str)
	}
	userIdSlice := strings.SplitN(str, " ", -1) //以空格为分割符将字符串中的值取出
	if userIdSlice == nil {
		return nil, errormsg.ERROR
	}
	var uid int
	var users = make([]map[string]interface{}, len(userIdSlice)-1) //创建一个map类型的切片   len-1去除字段最后一位空格
	var user User                                                  //创建一个User类型接受数据
	for i, id := range userIdSlice {
		if id == "" {
			break
		}
		fmt.Printf("id%vis:%v\n", i, id)
		users[i] = make(map[string]interface{}) //为每一个切片创建map
		uid, _ = strconv.Atoi(id)
		_ = Db.QueryRow("select id,avatar,nickName,introduction from users where id=?", uid).Scan(&user.ID, &user.Avatar, &user.NickName, &user.Introduction)
		users[i] = map[string]interface{}{
			"user_id":      user.ID,
			"avatar":       user.Avatar,
			"nickName":     user.NickName,
			"introduction": user.Introduction}
	}
	return users, errormsg.SUCCESS
}

//获取浏览记录
func GetViewList(username string) ([]map[string]interface{}, int) {
	var str string
	row := Db.QueryRow("select viewList from users where username=?", username)
	_ = row.Scan(&str)
	postIdSlice := strings.SplitN(str, " ", -1) //以空格为分割符将字符串中的值取出
	if postIdSlice == nil {
		return nil, errormsg.ERROR
	}
	var pid int
	var posts = make([]map[string]interface{}, len(postIdSlice)-1) //创建一个map类型的切片   len-1去除字段最后一位空格
	var post Post                                                  //创建一个post类型接受数据
	for i, id := range postIdSlice {
		if id == "" {
			break
		}
		posts[i] = make(map[string]interface{}) //为每一个切片创建map
		pid, _ = strconv.Atoi(id)
		_ = Db.QueryRow("select id,title,topicId from posts where id=?", pid).Scan(&post.Id, &post.Title, &post.TopicId)
		posts[i] = map[string]interface{}{
			"post_id": post.Id,
			"title":   post.Title,
			"topicId": post.TopicId,
		}
	}
	return posts, errormsg.SUCCESS
}

//获取关注你的用户
func FocusForU(username string) ([]map[string]interface{}, int) {
	var str string
	row := Db.QueryRow("select focusForU from users where username=?", username)
	_ = row.Scan(&str)
	userIdSlice := strings.SplitN(str, " ", -1) //以空格为分割符将字符串中的值取出
	if userIdSlice == nil {
		return nil, errormsg.ERROR
	}
	var uid int
	var users = make([]map[string]interface{}, len(userIdSlice)-1) //创建一个map类型的切片   len-1去除focus字段最后一位空格
	var user User                                                  //创建一个User类型接受数据
	for i, id := range userIdSlice {
		if id == "" {
			break
		}
		fmt.Printf("id%vis:%v\n", i, id)
		users[i] = make(map[string]interface{}) //为每一个切片创建map
		uid, _ = strconv.Atoi(id)
		_ = Db.QueryRow("select id,avatar,nickName,introduction from users where id=?", uid).Scan(&user.ID, &user.Avatar, &user.NickName, &user.Introduction)
		users[i] = map[string]interface{}{
			"user_id":  user.ID,
			"avatar":   user.Avatar,
			"nickName": user.NickName,
		}
	}
	return users, errormsg.SUCCESS
}
