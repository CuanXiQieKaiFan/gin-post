package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"nonoSheep/util/errormsg"
	"strings"
	"time"
)

var jwtKey = "984af51f4fhy9jt"
var JwtKey = []byte(jwtKey) //密钥
var code int

type MyClaim struct { //要跟User模型保持一致
	Username string `json:"username"`
	jwt.StandardClaims
}

type RClaim struct {
	Ref_Token string `json:"ref_token"`
	jwt.StandardClaims
}

//生成token
func SetToken(username string) (string, int) {
	expireTime := time.Now().Add(1 * time.Hour) //设置token的有效时间 1h
	SetClaim := MyClaim{
		username,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), //unix()给一个时间戳
			Issuer:    "ginblog",
		},
	}

	reqClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, SetClaim) //使用指定的签名方法HS256和声明类型创建新令牌
	token, err := reqClaim.SignedString(JwtKey)                     //返回完整的已签名令牌  JwtKey 加盐
	if err != nil {
		return "", errormsg.ERROR
	}
	return token, errormsg.SUCCESS
}

//生成refresh_token
func RefreshToken(username string) string {
	expireTime := time.Now().Add(10 * time.Hour) //设置RefreshToken的有效时间 10h
	SetClaim := MyClaim{
		username,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), //unix()给一个时间戳
			Issuer:    "nonosheep",
		},
	}

	reqClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, SetClaim) //使用指定的签名方法HS256和声明类型创建新令牌
	Rtoken, err := reqClaim.SignedString(JwtKey)                    //返回完整的已签名令牌  JwtKey 加盐
	if err != nil {
		return ""
	}
	return Rtoken
}

//验证token
func CheckToknen(token string) (*MyClaim, int) {
	settoken, _ := jwt.ParseWithClaims(token, &MyClaim{}, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})

	if key := settoken.Claims.(*MyClaim); settoken.Valid {
		return key, errormsg.SUCCESS
	} else {
		return nil, errormsg.ERROR
	}
}

//jwt中间件
func JwtToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenHeader := c.Request.Header.Get("Authorization")

		//检查token是否存在
		if tokenHeader == "" {
			code = errormsg.ERROR_TOKEN_NOT_EXIST
			c.JSON(http.StatusOK, gin.H{
				"status":  code,
				"message": errormsg.GetErrMsg(code),
			})
			c.Abort()
			return
		}
		checkToken := strings.SplitN(tokenHeader, " ", 2) //对于strings.SplitN如下
		checkRToken := strings.SplitN(tokenHeader, " ", 2)
		//检查token类型是否正确
		if len(checkToken) != 2 && checkToken[0] != "Bearer" {
			code = errormsg.ERROR_TOKEN_TYPE_WRONG
			c.JSON(http.StatusOK, gin.H{
				"status":  code,
				"message": errormsg.GetErrMsg(code),
			})
			c.Abort()
			return
		}
		//检查token是否错误
		key, Tcode := CheckToknen(checkToken[1])
		if Tcode == errormsg.ERROR {
			code = errormsg.ERROR_TOKEN_WRONG
			c.JSON(http.StatusOK, gin.H{
				"status":  code,
				"message": errormsg.GetErrMsg(code),
			})
			c.Abort()
			return
		}
		//检查token是否超时 超时则用refreshToken返回新token
		if time.Now().Unix() > key.ExpiresAt {
			newKey, _ := CheckToknen(checkRToken[1])
			if time.Now().Unix() > newKey.ExpiresAt {
				code = errormsg.ERROR_TOKEN_RUNTIME
				c.JSON(http.StatusOK, gin.H{
					"status":  code,
					"message": errormsg.GetErrMsg(code),
				})
				c.Abort()
				return
			} else {
				newToken, code := SetToken(newKey.Username)
				c.JSON(http.StatusOK, gin.H{
					"status":   code,
					"message":  errormsg.GetErrMsg(code),
					"NewToken": newToken,
				})
			}

		}

		c.Set("username", key.Username) //绑定username,存到当前上下文中 c.context
		c.Next()
	}
}

func CheckHeader(tokenHeader string) int {
	//检查token是否存在
	if tokenHeader == "" {
		code = errormsg.ERROR_TOKEN_NOT_EXIST
		return code
	}
	checkToken := strings.SplitN(tokenHeader, " ", 2)
	//检查token类型是否正确
	//if len(checkToken) != 2 && checkToken[0] != "Bearer" {
	//	code = errormsg.ERROR_TOKEN_TYPE_WRONG
	//	return code
	//}
	//检查token是否错误
	key, Tcode := CheckToknen(checkToken[1])
	if Tcode == errormsg.ERROR {
		code = errormsg.ERROR_TOKEN_WRONG
		return code
	}
	//检查token是否超时
	if time.Now().Unix() > key.ExpiresAt {
		return errormsg.ERROR_TOKEN_RUNTIME
	}
	return errormsg.SUCCESS
}
