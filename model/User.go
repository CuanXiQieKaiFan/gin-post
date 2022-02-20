package model

import (
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/scrypt"
	"log"
	"nonoSheep/middleware"
	"nonoSheep/util/errormsg"
)

type User struct {
	ID             int    `json:"id"`
	Gender         string ` json:"gender"`
	NickName       string `json:"nickName"`
	QQ             string `json:"qq"`
	Birthday       string `json:"birthday"`
	Email          string `json:"email"`
	Avatar         string `json:"avatar"`
	Introduction   string `json:"introduction"`
	Phone          string `json:"phone"`
	Username       string `json:"username" validate:"required,min=4,max=12"`
	Password       string `json:"password,omitempty" validate:"required,min=6,max=20"`
	Collection     string `json:"collection,omitempty"`
	Focus          string `json:"focus,omitempty"`
	PraisePost     string `json:"praise_post,omitempty"`       //存储点赞的帖子id
	PraiseCom      string `json:"praise_com,omitempty"`        //存储点赞的评论id
	PraisePostForU string `json:"praise_post_for_u,omitempty"` //存储为你帖子点赞的用户
	PraiseComForU  string `json:"praise_Com_for_u,omitempty"`  //存储为你评论点赞的用户
	FocusForU      string `json:"focus_for_u,omitempty"`       //存储关注你的用户列表
	Viewlist       string `json:"viewlist,omitempty"`          //浏览记录
}

//刷新token
func RToken(refToken string) (string, string, int) {
	//生成refresh_token
	var code int
	//验证token
	u, _ := middleware.CheckToknen(refToken)
	newToken, code := middleware.SetToken(u.Username)
	newRToken := middleware.RefreshToken(u.Username)
	return newToken, newRToken, code

}

//查询用户是否存在
func CheckUser(name string) (code int) {
	var users User
	dst := "	select * from users where username = ?"
	Db.QueryRow(dst, name).Scan(&users)
	if users.ID > 0 {
		return errormsg.ERROR_USERNAME_USED
	}
	return errormsg.SUCCESS
}

//密码加密
func ScryptPw(password string) string {
	const KeyLen = 10
	salt := make([]byte, 8)
	salt = []byte{66, 21, 14, 5, 56, 46, 13, 18}
	HasPw, err := scrypt.Key([]byte(password), salt, 16384, 8, 1, KeyLen)
	if err != nil {
		log.Fatal(err)
	}
	fpw := base64.StdEncoding.EncodeToString(HasPw)
	return fpw
}

//注册用户
func Register(data *User) int {
	data.Password = ScryptPw(data.Password)
	dst := "insert into users (users.username,users.password) values (?,?)"
	_, err = Db.Exec(dst, data.Username, data.Password)
	if err != nil {
		return errormsg.ERROR
	}
	return errormsg.SUCCESS
}

//登录验证
func CheckLogin(username string, password string) int {
	var user User
	row := Db.QueryRow("select username,password from users where username = ?", username)
	err := row.Scan(&user.Username, &user.Password)
	if err != nil {
		fmt.Printf("err is:%v\n", err)
	}
	if user.Username == "" {
		return errormsg.ERROR_USER_NOT_EXIST
	}
	if ScryptPw(password) != user.Password {
		return errormsg.ERROR_PASSWORD_WRONG
	}

	return errormsg.SUCCESS
}

//修改密码
func UpdatePwd(data *User, newPwd string) int {
	var key string
	e := Db.QueryRow("select password from users where username = ?", data.Username).Scan(&key)
	if e != nil {
		return errormsg.ERROP_DATABASE_SCAN_FALL
	}
	if key == data.Password {
		_, err = Db.Exec("update users set password = ? where username=?", newPwd, data.Username)
		if err != nil {
			return errormsg.ERROR
		} else {
			return errormsg.SUCCESS
		}
	}
	return errormsg.ERROR_PASSWORD_WRONG

}

//获取单个用户信息
func GetUser(id int) (int, User) {
	var user User
	err := Db.QueryRow("select id,avatar,nickName,introduction,phone,qq,gender,email,birthday,username from users where id = ?", id).Scan(&user.ID, &user.Avatar, &user.NickName, &user.Introduction, &user.Phone, &user.QQ, &user.Gender, &user.Email, &user.Birthday, &user.Username)
	if err != nil {
		fmt.Printf("err is:%v\n", err)
		code := errormsg.ERROR
		return code, user
	} else {
		code := errormsg.SUCCESS
		return code, user
	}
}

//更改用户信息
func EditUser(data *User) int {
	_, err = Db.Exec("update users set gender=?,nickName=?,qq=?,birthday=?,email=?,avatar=?,introduction=?,phone=? where username = ?", data.Gender, data.NickName, data.QQ, data.Birthday, data.Email, data.Avatar, data.Introduction, data.Phone, data.Username)
	if err != nil {
		fmt.Printf("err is:%v\n", err)
		return errormsg.ERROR
	}
	return errormsg.SUCCESS
}

//取出用户信息
func OutPutUser(username string) User {
	var user User
	row := Db.QueryRow("select id,avatar,nickName,introduction,phone,qq,gender,email,birthday,username from users where username =?", username)
	err := row.Scan(&user.ID, &user.Avatar, &user.NickName, &user.Introduction, &user.Phone, &user.QQ, &user.Gender, &user.Email, &user.Birthday, &user.Username)
	if err != nil {
		return User{}
	}
	return user
}
