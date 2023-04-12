package service

import (
	"log"
	"testing"
	cfg "tiktok/config"
	"tiktok/initial"
)

func Init() {
	_ = initial.InitConfig("../config/config.yaml")
	_ = initial.InitMySQL()
	_ = initial.InitMinIO()
}

// TestGetSnapshot 测试GetSnapshot
func TestGetSnapshot(t *testing.T) {
	config := cfg.Config.MinIO
	videoPath := "http://" + config.Url + ":" + config.APIPort + "/video/bear.mp4"

	if err := GetSnapshot(videoPath, 1); err != nil {
		log.Fatal(err)
	} else {
		println("Success")
	}
}

// TestCreateVideo 测试CreateVideo
func TestCreateVideo(t *testing.T) {
	Init()
	var userId int64 = 25
	fileName := "../public/test.mp4"

	req := PublishRequest{
		UserId:   userId,
		FilePath: fileName,
		Title:    "test",
	}

	if err := CreateVideo(req); err != nil {
		log.Fatal(err)
	}
}

func TestGetVideoList(t *testing.T) {
	Init()
	video, err := GetVideoList(1)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range video {
		println(v.VideoId, v.VideoTitle, v.PlayURL, v.CoverURL)
	}
}
