package model

import (
	"fmt"
	"nonoSheep/util/errormsg"
)

type Comment struct {
	ID                int    `json:"id"`
	IsPraised         bool   `json:"is_praised"`
	IsFocus           bool   `json:"is_focus"`
	PostId            int    `json:"post_id"`
	Content           string `json:"content"`
	Pictures          string `json:"pictures"`
	UserId            int    `json:"user_id,omitempty"`
	Avatar            string `json:"avatar"`
	NickName          string `json:"nick_name,omitempty"`
	ReplyUserId       int    `json:"reply_user_id,omitempty"`
	ReplyUserNickname string `json:"reply_user_nickname,omitempty"`
	ReplyComId        int    `json:"reply_com_id,omitempty"`
	PraiseCount       int    `json:"praise_count"`
	PublishTime       string `json:"publish_time"`
}

//获取帖子下的评论
func GetComment(model int, targId int, size int, page int) ([]Comment, int) {
	var com = make([]Comment, size)
	num := (page - 1) * size
	if model == 1 { //获取帖子下的评论
		rows, err := Db.Query("select id,postId,publishTime,content,userId,avatar,nickName,praiseCount,isFocus,isPraised from comments where postId = ? limit ?,?", targId, num, size)
		if err != nil {
			fmt.Printf("err is:%v", err)
			return nil, errormsg.ERROR
		} else {
			var i int
			for rows.Next() {
				_ = rows.Scan(&com[i].ID, &com[i].PostId, &com[i].PublishTime, &com[i].Content, &com[i].UserId, &com[i].Avatar, &com[i].NickName, &com[i].PraiseCount, &com[i].IsFocus, &com[i].IsPraised)
				i++
			}
			return com, errormsg.SUCCESS
		}
	}
	if model == 2 { //2级评论的id
		rows, err := Db.Query("select id,publishTime,content,replyUserId,avatar,replyUserNickName,praiseCount,isFocus,isPraised from comments where id = ?", targId)
		if err != nil {
			return nil, errormsg.ERROR
		} else {
			var i int
			for rows.Next() {
				_ = rows.Scan(&com[i].ID, &com[i].PublishTime, &com[i].Content, &com[i].ReplyUserId, &com[i].Avatar, &com[i].ReplyUserNickname, &com[i].PraiseCount, &com[i].IsFocus, &com[i].IsPraised)
				i++
			}
			return com, errormsg.SUCCESS
		}
	}
	return com, errormsg.ERROR_COMMENT_NOTEXIST
}

//发布评论
func IssueCom(mod int, targId int, com *Comment, username string) (int, int) {
	_ = Db.QueryRow("select id,nickName,avatar from users where username =?", username).Scan(&com.UserId, &com.NickName, &com.Avatar)
	if mod == 1 { //发布帖子下的评论
		result, err := Db.Exec("insert into comments (content,pictures,postId,userId,nickName,avatar) values (?,?,?,?,?,?)", com.Content, com.Pictures, targId, com.UserId, com.NickName, com.Avatar)
		if err != nil {
			return 0, errormsg.ERROR
		} else {
			insertId, _ := result.LastInsertId()
			var comCount int
			_ = Db.QueryRow("select commentCount from posts where id=?", targId).Scan(&comCount)
			comCount += 1
			_, _ = Db.Exec("update posts set commentCount=? where id=?", comCount, targId)
			return int(insertId), errormsg.SUCCESS
		}
	}
	if mod == 2 { //1级评论下评论
		res, err := Db.Exec("insert into comments (content,pictures,replyComId,replyUserId,replyUserNickName,avatar) values (?,?,?,?,?,?)", com.Content, com.Pictures, targId, com.UserId, com.NickName, com.Avatar)
		if err != nil {
			return 0, errormsg.ERROR
		} else {
			insertId, _ := res.LastInsertId()
			return int(insertId), errormsg.SUCCESS
		}
	}
	return -1, errormsg.ERROR_CREATE_COM
}

//更新评论
func UpdateComment(id int, newCom *Comment) int {
	_, err := Db.Exec("update comments set content=?,pictures=? where id=?", newCom.Content, newCom.Pictures, id)
	if err != nil {
		return errormsg.ERROR
	}
	return errormsg.SUCCESS
}

//删除评论
func DeleteComment(id int) int {
	_, err := Db.Exec("DELETE from comments where id = ?", id)
	if err != nil {
		return errormsg.ERROR
	}
	return errormsg.SUCCESS
}
