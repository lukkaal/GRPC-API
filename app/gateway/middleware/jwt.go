package middleware

import (
	"time"

	"github.com/lukkaal/GRPC-API/pkg/jsonres"
	"github.com/lukkaal/GRPC-API/pkg/utils/jwt"

	Err "github.com/lukkaal/GRPC-API/pkg/errcode"

	"github.com/gin-gonic/gin"
)

func JWT() gin.HandlerFunc {
	handler := func(c *gin.Context) {
		var code int
		var data interface{}
		code = 200

		// task 验证用户 token
		token := c.GetHeader("Authorization")

		if token == "" {
			code = 404
			c.JSON(200, gin.H{
				"status": code,
				"msg":    Err.GetMsg(code),
				"data":   data,
			})
			c.Abort()
		}

		// validate the token
		claims, err := jwt.ParseToken(token)

		// no timeout no err
		if err != nil {
			code = Err.ErrorAuthCheckTokenFail
		} else if claims.ExpiresAt == nil ||
			time.Now().After(claims.ExpiresAt.Time) {
			code = Err.ErrorAuthCheckTokenTimeout
		}

		if code != Err.SUCCESS {
			c.JSON(200, gin.H{
				"status": code,
				"msg":    Err.GetMsg(code),
				"data":   data,
			})
			c.Abort()
			return
		}

		// set token_info into gin.context
		// c.Request: *http.Request
		c.Request = c.Request.WithContext(jsonres.NewContext(
			c.Request.Context(), &jsonres.UserInfo{Id: claims.UserId}))

		jsonres.InitUserInfo(c.Request.Context())

		c.Next()
	}

	return handler
}
