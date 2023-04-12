package minio

import (
	"github.com/minio/minio-go/v6"
	"mime/multipart"
	"os"
	cfg "tiktok/config"
	"tiktok/initial"
)

type Video struct {
	bucketName    string
	objectName    string
	File          *os.File
	ContentType   string
}

// CreateVideo 上传视频
func CreateVideo(videoName string, file multipart.File, fileSize int64) (videoURL string, err error) {
	bucketName := "video"
	objectName := videoName
	config := cfg.Config.MinIO
	videoURL = "http://" + config.Url + ":" + config.APIPort + "/" + bucketName + "/" + objectName

	_, err = initial.MinioClient.PutObject(bucketName, objectName, file, fileSize, minio.PutObjectOptions{})
	if err != nil {
		return
	}
	return videoURL, nil
}

//CreateImage 上传视频封面
func CreateImage(coverName string, file *os.File) (coverURL string, err error) {
	bucketName := "photo"
	objectName := coverName
	fileStat, err := file.Stat()
	if err != nil {
		return
	}
	config := cfg.Config.MinIO
	coverURL = "http://" + config.Url + ":" + config.APIPort + "/" + bucketName + "/" + objectName

	_, err = initial.MinioClient.PutObject(bucketName, objectName, file, fileStat.Size(), minio.PutObjectOptions{})
	if err != nil {
		return
	}
	return coverURL, nil
}