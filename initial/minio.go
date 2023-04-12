package initial

import (
	"github.com/minio/minio-go/v6"
	cfg "tiktok/config"
)

var MinioClient *minio.Client

func InitMinIO() error {
	config := cfg.Config.MinIO

	var err error
	MinioClient, err = minio.New(config.Url + ":" + config.APIPort, config.AccessKeyID, config.SecretAccessKey, false)
	return err
}