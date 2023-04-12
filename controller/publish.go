package controller

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"tiktok/service"
	"tiktok/utils"
)

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

// Publish 上传视频
func Publish(c *gin.Context) {
	title := c.PostForm("title")
	claims, exist := c.Get("claims")
	if !exist {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "user doesn't exist",
		})
		return
	}

	userId := claims.(*utils.MyClaims).UserID
	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "Get Video Error:" + err.Error(),
		})
		return
	}

	fileName := data.Filename
	if err != nil {
		log.Errorf("Get VideoName Error: %v", err)
	}

	req := service.PublishRequest{
		UserId:     int64(userId),
		FileHeader: data,
		Title:      title,
	}

	if err := service.CreateVideo(req); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "Update Video Error:" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  fileName + "Upload Video Successful",
	})
}

// PublishList 获取视频列表
func PublishList(c *gin.Context) {
	_, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "user doesn't exist",
		})
		return
	}
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "Get UserId Error:" + err.Error(),
		})
		return
	}

	user, err := service.UserServiceInstance().GetUserInfo(int(userId), int(userId))

	author := User{Id: user.User.Id, Name: user.User.Name, FollowCount: user.User.FollowCount,
		FollowerCount: user.User.FollowerCount, IsFollow: false, Avatar: user.User.Avatar,
		BackgroundImage: user.User.BackgroundImage, Signature: user.User.Signature,
		TotalFavorited: user.User.TotalFavorited, WorkCount: user.User.WorkCount,
		FavoriteCount: user.User.FavoriteCount}

	list, err := service.GetVideoList(userId)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "Get VideoList Error:" + err.Error(),
		})
		return
	}
	videoList := make([]Video, len(list))
	for i := 0; i < len(list); i++ {
		favoriteCount, err := service.Favorite().QueryFavoriteNumberByVideoId(list[i].VideoId)
		if err != nil {
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  "Get Favorite Number Error:" + err.Error(),
			})
			return
		}
		commentCount, err := service.Comment().QueryCommentNumberByVideoId(list[i].VideoId)
		if err != nil {
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  "Get Comment Number Error:" + err.Error(),
			})
			return
		}
		isFavorite, err := service.Favorite().IsFavor(userId, list[i].VideoId)
		if err != nil {
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  "Get IsFavorite Number Error:" + err.Error(),
			})
			return
		}
		videoList[i] = Video{
			Id: list[i].VideoId,
			Author: author,
			PlayUrl: list[i].PlayURL,
			CoverUrl: list[i].CoverURL,
			FavoriteCount: int64(favoriteCount),
			CommentCount: int64(commentCount),
			IsFavorite: isFavorite,
			Title: list[i].VideoTitle,
		}
	}

	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg: "success",
		},
		VideoList: videoList,
	})
}
