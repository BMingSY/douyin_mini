package controller

import (
	"net/http"
	"tiktok/service"
	"tiktok/utils"

	"github.com/gin-gonic/gin"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
var usersLoginInfo = map[string]User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

// Register 新用户注册
func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	userService := service.UserServiceInstance()

	response, err := userService.RegisterUser(username, password)
	if err != nil { // error
		c.JSON(http.StatusOK, &response)
		return
	}
	c.JSON(http.StatusOK, &response)
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	userService := service.UserServiceInstance()

	response, err := userService.Login(username, password)
	if err != nil { // error
		c.JSON(http.StatusOK, &response)
		return
	}
	c.JSON(http.StatusOK, &response)
}

func UserInfo(c *gin.Context) {
	claims, st := c.Get("claims")
	if !st {
		return
	}
	userId := claims.(*utils.MyClaims).UserID
	userService := service.UserServiceInstance()
	response, err := userService.GetUserInfo(userId, userId)
	if err != nil { // error
		c.JSON(http.StatusOK, &response)
		return
	}
	c.JSON(http.StatusOK, &response)
}
