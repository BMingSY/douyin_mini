package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
	"strconv"
	"tiktok/service"
	"tiktok/utils"
	"time"
)

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

// Feed 返回视频流列表
func Feed(c *gin.Context) {
	latestTime := c.Query("latest_time")
	claims, exist := c.Get("claims")
	if !exist {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "user doesn't exist",
		})
		return
	}
	lastTime, err := strconv.ParseInt(latestTime, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "Get Latest Time Error:" + err.Error(),
		})
		return
	}
	if lastTime == 0 {
		lastTime = time.Now().UnixMilli()
	}

	userId := claims.(*utils.MyClaims).UserID
	feed, err := service.GetFeedList(lastTime)
	hasMore := false
	if len(feed) == 3 {
		hasMore = true
		feed = feed[:len(feed)-1]
	}

	videoList := make([]Video, len(feed))

	for i := 0; i < len(feed); i++ {
		user, err := service.UserServiceInstance().GetUserInfo(userId, int(feed[i].UserId))
		if err != nil {
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  "Get Author Info Error:" + err.Error(),
			})
			return
		}
		author := User{
			Id:              user.User.Id,
			Name:            user.User.Name,
			FollowCount:     user.User.FollowCount,
			FollowerCount:   user.User.FollowerCount,
			IsFollow:        false,
			Avatar:          user.User.Avatar,
			BackgroundImage: user.User.BackgroundImage,
			Signature:       user.User.Signature,
			TotalFavorited:  user.User.TotalFavorited,
			WorkCount:       user.User.WorkCount,
			FavoriteCount:   user.User.FavoriteCount}

		favoriteCount, err := service.Favorite().QueryFavoriteNumberByVideoId(feed[i].VideoId)
		if err != nil {
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  "Get Favorite Count Error:" + err.Error(),
			})
			return
		}

		commentCount, err := service.Comment().QueryCommentNumberByVideoId(feed[i].VideoId)
		if err != nil {
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  "Get Comment Count Error:" + err.Error(),
			})
			return
		}

		isFavorite, err := service.Favorite().IsFavor(int64(userId), feed[i].VideoId)
		if err != nil {
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  "Get Favorite Situation Error:" + err.Error(),
			})
			return
		}

		videoList[i] = Video{
			Id:            feed[i].VideoId,
			Author:        author,
			PlayUrl:       feed[i].PlayURL,
			CoverUrl:      feed[i].CoverURL,
			FavoriteCount: int64(favoriteCount),
			CommentCount:  int64(commentCount),
			IsFavorite:    isFavorite,
			Title:         feed[i].VideoTitle,
		}
	}
	nextTime := feed[len(feed)-1].CreateTime.UnixMilli()

	sort.Slice(videoList, func(i, j int) bool {
		return videoList[i].FavoriteCount > videoList[j].FavoriteCount
	})

	if hasMore == false {
		c.JSON(http.StatusOK, FeedResponse{
			Response:  Response{StatusCode: 0},
			VideoList: videoList,
			NextTime:  0,
		})
	} else {
		c.JSON(http.StatusOK, FeedResponse{
			Response:  Response{StatusCode: 0},
			VideoList: videoList,
			NextTime:  nextTime,
		})
	}
}
