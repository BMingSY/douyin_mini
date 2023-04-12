package middleware

import (
	"net/http"
	"tiktok/utils"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		if token == "" {
			token = c.PostForm("token")
			if token == "" {
				c.JSON(http.StatusOK, Response{
					1,
					"未登录或非法访问",
				})
				log.Info("未登录或非法访问")
				c.Abort()
				return
			}
		}

		// parseToken 解析token包含的信息
		claims, err := utils.ParseToken(token)
		if err != nil {
			if err == utils.ErrTokenExpired {
				c.JSON(http.StatusOK, Response{
					1,
					"授权已过期",
				})
				log.Info("授权已过期")
				c.Abort()
				return
			}
			c.JSON(http.StatusOK, Response{
				1,
				err.Error(),
			})
			c.Abort()
			return
		}
		log.Info("验证通过")
		c.Set("claims", claims)
		c.Next()
	}
}
