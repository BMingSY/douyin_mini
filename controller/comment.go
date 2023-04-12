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

type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	Response
	Comment Comment `json:"comment,omitempty"`
}

// CommentAction 评论或者删除评论
func CommentAction(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "user don't exists",
		})
		return
	}
	userId := int64(claims.(*utils.MyClaims).UserID)

	log.Infof("user id :[%v]", userId)
	actionType, err := strconv.ParseInt(c.Query("action_type"), 10, 32)
	if err != nil {
		log.Errorf("comment action get action type err:[%v]", err)
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "action type error",
		})
		return
	}

	videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		log.Errorf("comment action get video id error:[%v]", err)
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "video id error",
		})
	}
	res, err := service.UserServiceInstance().GetUserInfo(int(userId), int(userId))
	if err != nil {
		log.Errorf("comment list get user info error:[%v]", err)
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "user doesn't exist",
		})
		return
	}

	user := User{
		Id:            res.User.Id,
		Name:          res.User.Name,
		FollowCount:   res.User.FollowerCount,
		FollowerCount: res.User.FollowerCount,
		IsFollow:      res.User.IsFollow,
	}
	if actionType == initial.CREATE_COMMENT {
		text := c.Query("comment_text")
		req := service.CommentCreateRequest{
			UserId:  userId,
			VideoId: videoId,
			Text:    text,
		}
		ID, CreateTime, err := service.Comment().CreateComment(req)
		if err != nil {
			log.Errorf("comment action service comment return error:[%v]", err)
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  "unknown error",
			})
			return
		}

		c.JSON(http.StatusOK, CommentActionResponse{Response: Response{StatusCode: 0},
			Comment: Comment{
				Id:         ID,
				User:       user,
				Content:    text,
				CreateDate: CreateTime.Format("01-02"),
			}})
	} else if actionType == initial.DELETE_COMMENT {
		commentId, err := strconv.ParseInt(c.Query("comment_id"), 10, 64)
		if err != nil {
			log.Errorf("comment action comment delete error:[%v]", err)
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  "unknown error",
			})
			return
		}

		if err := service.Comment().DeleteComment(commentId); err != nil {
			log.Errorf("comment action service comment return error:[%v]", err)
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  "unknown error",
			})
			return
		}
		c.JSON(http.StatusOK, CommentActionResponse{Response: Response{StatusCode: 0},
			Comment: Comment{
				Id:      commentId,
				User:    user,
				Content: "",
			}})

	}
}

// CommentList 获取视频评论列表
func CommentList(c *gin.Context) {

	videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		log.Errorf("comment list get video id error:[%v]", err)
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "video id error",
		})
		return
	}
	res, err := service.Comment().QueryCommentList(videoId)
	if err != nil {
		log.Errorf("comment list service comment return error:[%v]", err)
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "unknown error",
		})
		return
	}
	commentList := make([]Comment, len(res))
	for i, r := range res {
		res, err := service.UserServiceInstance().GetUserInfo(int(r.UserId), int(r.UserId))
		if err != nil {
			log.Warning("comment list get user info warning:[%v]", err)
		}

		commentList[i] = Comment{
			Id: r.ID,
			User: User{
				Id:            res.User.Id,
				Name:          res.User.Name,
				FollowCount:   res.User.FollowCount,
				FollowerCount: res.User.FollowerCount,
				IsFollow:      res.User.IsFollow,
			},
			Content:    r.Text,
			CreateDate: r.CreateTime.Format("01-02"),
		}
	}
	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0},
		CommentList: commentList,
	})

}
