package service

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/disintegration/imaging"
	log "github.com/sirupsen/logrus"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"io"
	"mime/multipart"
	"os"
	"strings"
	"tiktok/database/minio"
	"tiktok/database/mysql"
)

type PublishRequest struct {
	UserId     int64
	FileHeader *multipart.FileHeader
	FilePath   string
	Title      string
}

// HashSHAFile 对视频进行SHA256哈希
func HashSHAFile(file multipart.File) (string, error) {
	var hashValue string
	defer file.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return hashValue, err
	}
	hashInBytes := hash.Sum(nil)
	hashValue = hex.EncodeToString(hashInBytes)
	return hashValue, nil
}

// GetSnapshot 获取视频第一帧作为视频封面
func GetSnapshot(videoPath string, frameNum int) error {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(videoPath).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		return err
	}

	img, err := imaging.Decode(buf)
	if err != nil {
		return err
	}

	file, err := os.OpenFile("./tiktok_photo.png", os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	//将img转换成os.File类型
	if err := imaging.Encode(file, img, imaging.PNG); err != nil {
		return err
	}
	return nil
}

// CreateVideo 上传视频
func CreateVideo(req PublishRequest) error {
	file, err := req.FileHeader.Open()
	if err != nil {
		log.Errorf("Open MultipartFile error: %v", err)
		return err
	}

	hashValue, err := HashSHAFile(file)
	if err != nil {
		return err
	}

	fileArrIndex := strings.LastIndex(req.FileHeader.Filename, ".")
	fileSuf := req.FileHeader.Filename[fileArrIndex:]
	videoName := hashValue + fileSuf
	fileSize := req.FileHeader.Size

	file, err = req.FileHeader.Open()
	if err != nil {
		log.Errorf("Open MultipartFile error: %v", err)
		return err
	}

	videoURL, err := minio.CreateVideo(videoName, file, fileSize)
	if err != nil {
		log.Errorf("Save video error: %v", err)
		return err
	}

	if err := file.Close(); err != nil {
		log.Errorf("Close file error: %v", err)
		return err
	}

	err = GetSnapshot(videoURL, 1)
	if err != nil {
		log.Errorf("Get coverimage error: %v", err)
		return err
	}

	img, err := os.Open("./tiktok_photo.png")
	if err != nil {
		log.Errorf("Open coverimage error: %v", err)
		return err
	}

	coverName := hashValue + ".png"
	coverURL, err := minio.CreateImage(coverName, img)
	if err != nil {
		log.Errorf("Save coverimage error: %v", err)
		return err
	}

	video := mysql.Video{
		UserId:     req.UserId,
		VideoTitle: req.Title,
		PlayURL:    videoURL,
		CoverURL:   coverURL,
	}
	if err := mysql.CreateVideo(video); err != nil {
		log.Errorf("Mysql create data error: %v", err)
		return err
	}
	return nil
}

// GetVideoList 获取用户的视频列表
func GetVideoList(userId int64) ([]mysql.Video, error) {
	return mysql.GetVideoList(userId)
}

// GetVideoInfoByVideoId 使用videoId查询视频信息
func GetVideoInfoByVideoId(videoId int64) (mysql.Video, error) {
	return mysql.GetVideoInfoByVideoId(videoId)
}

// GetFeedList 获取feed流视频
func GetFeedList(latestTime int64) ([]mysql.Video, error) {
	return mysql.GetFeedList(latestTime)
}
