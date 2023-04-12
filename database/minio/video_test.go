package minio

import (
	"os"
	"testing"
	"tiktok/initial"
)

func Init() {
	_ = initial.InitConfig("../../config/config.yaml")
	_ = initial.InitMinIO()
}

//TestCreateVideo 测试CreateVideo
func TestCreateVideo(t *testing.T) {
	Init()
	fileName := "../../public/bear.mp4"

	file, err := os.Open(fileName)
	if err != nil {
		println(err)
	}
	if _, err = CreateVideo("1", file, 54363643); err != nil {
		println(err)
	}
}

//TestCreateImage 测试CreateImage
func TestCreateImage(t *testing.T) {
	Init()

	fileName := "../../public/bear.mp4"

	file, err := os.Open(fileName)
	if err != nil {
		println(err)
	}
	if _, err = CreateImage("1.png", file); err != nil {
		println(err)
	}
}