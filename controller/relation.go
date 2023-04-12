package controller

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"tiktok/initial"
	"tiktok/service"
	"tiktok/utils"
)

type UserListResponse struct {
	Response
	UserList []User `json:"user_list"`
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "user doesn't exist",
		})
		return
	}
	userId := int64(claims.(*utils.MyClaims).UserID)

	followUserId, err := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	if err != nil {
		log.Errorf("Relation action followUserId error :[%v}]", err)
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "followUserID error",
		})
		return
	}
	op, err := strconv.ParseInt(c.Query("action_type"), 10, 32)
	if err != nil {
		log.Errorf("Relation action action_type error :[%v]", err)
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "action_type error",
		})
		return
	}

	relationReq := service.RelationRequest{
		UserID:       userId,
		FollowUserId: followUserId,
	}

	if op == initial.CREATE_FOLLOW {
		err := service.Relation().CreateFollow(relationReq)
		if err != nil {
			log.Errorf("  create follow action error: [%v]", err)
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  "create follow error",
			})
			return
		}
	} else if op == initial.DELETE_FOLLOW {
		err := service.Relation().DeleteFollow(relationReq)
		if err != nil {
			log.Errorf("  delete follow action error: [%v]", err)
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  "delete follow error",
			})
			return
		}
	}

	c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "success"})
}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		log.Errorf("FollwList get user id error:[%v]", err)
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "follow list error",
			},
			UserList: []User{},
		})
		return
	}

	res, err := service.Relation().QueryFollowList(userId)
	if err != nil {
		log.Errorf("FollwList get user id error:[%v]", err)
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "follow list error",
			},
			UserList: []User{},
		})
	}

	ans := make([]User, len(res))
	for i := 0; i < len(res); i++ {
		ans[i] = User{
			Id:            res[i].Id,
			Name:          res[i].Name,
			FollowCount:   res[i].FollowCount,
			FollowerCount: res[i].FollowerCount,
			IsFollow:      res[i].IsFollow,
		}
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		UserList: ans,
	})
}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		log.Errorf("FollwerList get user id error:[%v]", err)
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "follower list error",
			},
			UserList: []User{},
		})
		return
	}

	res, err := service.Relation().QueryFollowerList(userId)
	if err != nil {
		log.Errorf("FollwerList get user id error:[%v]", err)
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "follower list error",
			},
			UserList: []User{},
		})
	}

	followerList := make([]User, len(res))
	for i := 0; i < len(res); i++ {
		followerList[i] = User{
			Id:            res[i].Id,
			Name:          res[i].Name,
			FollowCount:   res[i].FollowCount,
			FollowerCount: res[i].FollowerCount,
			IsFollow:      res[i].IsFollow,
		}
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		UserList: followerList,
	})
}

// FriendList all users have same friend list
func FriendList(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		log.Errorf("FollwList get user id error:[%v]", err)
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "follow list error",
			},
			UserList: []User{},
		})
		return
	}

	res, err := service.Relation().QueryFriendList(userId)
	if err != nil {
		log.Errorf("FollwList get user id error:[%v]", err)
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "follow list error",
			},
			UserList: []User{},
		})
	}

	ans := make([]User, len(res))
	for i := 0; i < len(res); i++ {
		ans[i] = User{
			Id:            res[i].Id,
			Name:          res[i].Name,
			FollowCount:   res[i].FollowCount,
			FollowerCount: res[i].FollowerCount,
			IsFollow:      res[i].IsFollow,
		}
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		UserList: ans,
	})
}
