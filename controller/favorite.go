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

// FavoriteAction 点赞或者取消点赞
func FavoriteAction(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "user doesn't exist",
		})
		return
	}
	userId := int64(claims.(*utils.MyClaims).UserID)

	log.Infof("user id :[%v]", userId)

	videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		log.Errorf("favorite action get video id error:[%v]", err)
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "video id error",
		})
		return
	}
	op, err := strconv.ParseInt(c.Query("action_type"), 10, 32)
	if err != nil {
		log.Errorf("favorite action get action type error:[%v]", err)
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "action type error",
		})
		return
	}
	favoriteReq := service.FavoriteRequest{
		UserId:  userId,
		VideoId: videoId,
	}
	if op == initial.CREATE_LIKE {
		if err := service.Favorite().CreateLike(favoriteReq); err != nil {
			log.Errorf("favorite action service favorite create like return error:[%v]", err)
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  "unknown error",
			})
			return
		}
	} else if op == initial.DELETE_LIKE {
		if err := service.Favorite().DeleteLike(favoriteReq); err != nil {
			log.Errorf("favorite action service favorite delete like return error:[%v]", err)
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  "unknown error",
			})
			return
		}
	}

	c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "success"})

}

// FavoriteList 获取user_id的点赞列表
func FavoriteList(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		log.Errorf("favorite list controller get user_id error:[%v]", err)
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "user_id error",
		})
		return
	}
	log.Infof("user id :[%v]", userId)

	res, err := service.Favorite().QueryFavoriteList(userId)
	if err != nil {
		log.Errorf("favorite list service favorite error:[%v]", err)
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "unknown error",
		})
		return
	}

	videoList := make([]Video, len(res))
	for i := 0; i < len(res); i++ {
		video, err := service.GetVideoInfoByVideoId(res[i])
		if err != nil {
			log.Errorf("Favorite list controller get video info error:[%v],video_id:[%v]", err, res[i])
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  "get video information failed",
			})
			return
		}
		user, err := service.UserServiceInstance().GetUserInfo(int(userId), int(video.UserId))
		if err != nil {
			log.Errorf("Favorite list controller get user info error:[%v],user_id:[%v]", err, video.UserId)
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  "get user information failed",
			})
			return
		}
		favoriteCnt, err := service.Favorite().QueryFavoriteNumberByVideoId(video.VideoId)
		if err != nil {
			log.Errorf("Favorite list controller get favorite cnt error:[%v],video_id:[%v]", err, video.VideoId)
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  "get favorite count failed",
			})
			return
		}
		commentCnt, err := service.Comment().QueryCommentNumberByVideoId(video.VideoId)
		if err != nil {
			log.Errorf("Favorite list controller get comment cnt error:[%v],video_id:[%v]", err, video.VideoId)
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  "get comment count failed",
			})
			return
		}
		videoList[i] = Video{
			Id: video.VideoId,
			Author: User{
				Id:            user.User.Id,
				Name:          user.User.Name,
				FollowCount:   user.User.FollowCount,
				FollowerCount: user.User.FollowerCount,
				IsFollow:      user.User.IsFollow,
			},
			PlayUrl:       video.PlayURL,
			CoverUrl:      video.CoverURL,
			FavoriteCount: int64(favoriteCnt),
			CommentCount:  int64(commentCnt),
			IsFavorite:    true,
		}
	}
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		VideoList: videoList,
	})
}
