package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func Cors() gin.HandlerFunc {

	return func(c *gin.Context) {
		cors.New(cors.Config{
			AllowAllOrigins: true, //允许所有跨域
			AllowMethods:    []string{"*"},
			AllowHeaders:    []string{"origin"},
			ExposeHeaders:   []string{"Content-Length", "Authorization"},
			//AllowCredentials: true,  //是否为cockie
			MaxAge: 12 * time.Hour,
		})
	}
}
