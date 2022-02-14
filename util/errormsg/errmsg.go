package errormsg

const (
	SUCCESS = 200
	ERROR   = 500

	// code =1000...用户模块错误
	ERROR_USERNAME_USED    = 1001
	ERROR_PASSWORD_WRONG   = 1002
	ERROR_USER_NOT_EXIST   = 1003
	ERROR_TOKEN_NOT_EXIST  = 1004
	ERROR_TOKEN_RUNTIME    = 1005
	ERROR_TOKEN_WRONG      = 1006
	ERROR_TOKEN_TYPE_WRONG = 1007
	ERROR_USER_NO_RIGHT    = 1008
	// code =2000...帖子模块错误
	ERROR_POST_NOTEXIST = 2001

	// code =3000...评论模块的错误
	ERROR_CREATE_COM       = 3001
	ERROR_COMMENT_NOTEXIST = 3002

	ERROR_PRAISE  = 4001
	ERROR_COLLECT = 4002
	ERROR_FOCUS   = 4003

	ERROP_DATABASE_SCAN_FALL = 5001
)

var codeMsg = map[int]string{
	SUCCESS:                "OK",
	ERROR:                  "FALL",
	ERROR_USERNAME_USED:    "用户名已存在",
	ERROR_PASSWORD_WRONG:   "密码错误",
	ERROR_USER_NOT_EXIST:   "用户不存在",
	ERROR_TOKEN_NOT_EXIST:  "TOKEN不存在",
	ERROR_TOKEN_RUNTIME:    "TOKEN已过期",
	ERROR_TOKEN_WRONG:      "TOKEN不正确",
	ERROR_TOKEN_TYPE_WRONG: "TOKEN格式不正确",
	ERROR_USER_NO_RIGHT:    "该用户无管理权限",

	ERROR_POST_NOTEXIST: "帖子不存在",

	ERROR_CREATE_COM:       "发布评论失败",
	ERROR_COMMENT_NOTEXIST: "评论不存在",

	ERROR_PRAISE:  "点赞失败",
	ERROR_COLLECT: "收藏失败",
	ERROR_FOCUS:   "关注失败",

	ERROP_DATABASE_SCAN_FALL: "绑定数据库数据失败",
}

func GetErrMsg(code int) string {
	return codeMsg[code]
}