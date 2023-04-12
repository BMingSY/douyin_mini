package controller

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"strconv"
	"tiktok/service"
	"tiktok/utils"
)

//var tempChat = map[string][]Message{}
//
//var messageIdSequence = int64(1)

type ChatResponse struct {
	//Response
	StatusCode  int32     `json:"status_code"`
	StatusMsg   string    `json:"status_msg,omitempty"`
	MessageList []Message `json:"message_list"`
}

// MessageAction no practical effect, just check if token is valid
func MessageAction(c *gin.Context) {
	userID := utils.GetUserID(c)
	toUserID, _ := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	content := c.Query("content")
	err := service.Message().AddMessage(userID, toUserID, content)
	if err != nil {
		//TODO 加密msg
		log.Errorf("add message err:[%v],userid[%v],touserID[%v],msg[%v]",
			err, userID, toUserID, content)
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "service internal error!"})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	}
}

// MessageChat all users have same follow list
func MessageChat(c *gin.Context) {
	userID := utils.GetUserID(c)
	toUserID, _ := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	preMsgTime, _ := strconv.ParseInt(c.Query("pre_msg_time"), 10, 64)

	res, err := service.Message().GetMessage(userID, toUserID, preMsgTime)
	if err != nil {
		log.Errorf("get message err:[%v]", err)
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		msgs := make([]Message, 0, len(res))
		log.Info(res)
		for _, val := range res {
			msgs = append(msgs, Message{
				Id:         int64(rand.Int()),
				ToUserID:   val.ToUserID,
				FromUserID: val.FromUserID,
				Content:    val.Content,
				CreateTime: int(val.CreateTime),
			})
		}
		c.JSON(http.StatusOK, ChatResponse{
			StatusCode:  0,
			StatusMsg:   "success",
			MessageList: msgs,
		})
	}
}
