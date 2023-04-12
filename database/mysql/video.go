package mysql

import (
	log "github.com/sirupsen/logrus"
	"tiktok/initial"
	"time"
)

type Video struct {
	VideoId    int64
	UserId     int64
	VideoTitle string
	PlayURL    string
	CoverURL   string
	CreateTime time.Time
}

func (Video) TableName() string {
	return "video"
}

// CreateVideo 上传视频
func CreateVideo(req Video) error {
	req.CreateTime = time.Now()

	err := initial.Database.Create(req).Error
	if err != nil {
		return err
	}
	return nil
}

// GetVideoList 获取用户上传视频的视频名
func GetVideoList(userId int64) ([]Video, error) {
	var videoList []Video
	if err := initial.Database.
		Where("user_id = ?", userId).
		Order("create_time DESC").
		Find(&videoList).Error; err != nil {
		return videoList, err
	}
	return videoList, nil
}

// GetVideoInfoByVideoId 通过视频id获取视频信息
func GetVideoInfoByVideoId(videoId int64) (Video, error) {
	var video Video
	if err := initial.Database.Where("video_id = ?", videoId).Find(&video).Error; err != nil {
		return video, err
	}
	return video, nil
}

// GetFeedList 获取Feed流的视频列表
func GetFeedList(latestTime int64) ([]Video, error) {
	video := make([]Video, 0)
	timeCursor := time.UnixMilli(latestTime)
	err := initial.Database.Table("video").
		Where("create_time < ?", timeCursor).
		Order("create_time DESC").Limit(3).Find(&video).Error
	if err != nil {
		log.Warnf("db get video err :%v", err)
		return nil, err
	}
	return video, nil
}
