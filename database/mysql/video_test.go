package mysql

import (
	log "github.com/sirupsen/logrus"
	"testing"
	"tiktok/initial"
	"time"
)


//TestCreateVideo 测试CreateVideo
func TestCreateVideo(t *testing.T) {
	_ = initial.InitConfig("../../config/config.yaml")
	_ = initial.InitMySQL()

	req := Video{VideoId: 1, UserId: -1,VideoTitle: "bear",PlayURL: "./1.mp4",CoverURL: "./1.png"}
	err := CreateVideo(req)
	if err != nil {
		println(err)
	}
}

//TestGetVideoList 测试GetVideoList
func TestGetVideoList(t *testing.T) {
	_ = initial.InitConfig("../../config/config.yaml")
	_ = initial.InitMySQL()
	video, err := GetVideoList(1)
	if err != nil {
		log.Fatal(err)
	} else {
		for _, v := range video {
			println(v.VideoId, v.PlayURL, v.CoverURL)
		}
	}
}

//TestGetVideoInfoByVideoId 测试GetVideoInfoByVideoId
func TestGetVideoInfoByVideoId(t *testing.T) {
	_ = initial.InitConfig("../../config/config.yaml")
	_ = initial.InitMySQL()
	video, err := GetVideoInfoByVideoId(3)
	if err != nil {
		log.Fatal(err)
	} else {
		println(video.PlayURL, video.CoverURL)
	}
}

//TestGetFeedList 测试GetFeedList
func TestGetFeedList(t *testing.T) {
	_ = initial.InitConfig("../../config/config.yaml")
	_ = initial.InitMySQL()

	latestTime := time.Now()
	timeLayout := "2006-01-02 15:04:05"

	println(latestTime.Format(timeLayout))
	feed, err := GetFeedList(latestTime.Unix())
	if err != nil {
		log.Fatal(err)
	}
	for _, video := range feed {
		println(video.VideoId, video.CreateTime.Format(timeLayout))
	}
}