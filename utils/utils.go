package utils

import (
	"github.com/gin-gonic/gin"
)

func GetUserID(c *gin.Context) int64 {
	claims, _ := c.Get("claims")
	return int64(claims.(*MyClaims).UserID)
}
