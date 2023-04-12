package redis

import (
	"fmt"
	"strconv"
	"testing"
	"tiktok/initial"
	"time"
)

func init() {
	_ = initial.InitConfig("../../config/config.yaml")
	_ = initial.InitRedis()
}

func TestIncrVideoLikeCnt(t *testing.T) {
	res, timeNow, err := IncrVideoLikeCnt(20, -3)
	if err != nil {
		panic(err)
	}
	fmt.Println(res, timeNow)
}

func TestIncrVideoCommentCnt(t *testing.T) {
	err := IncrVideoCommentCnt(22, 1)
	if err != nil {
		panic(err)
	}
}

func TestGetVideoLikeCnt(t *testing.T) {
	res, _, err := GetVideoLikeCnt(17)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
func TestGetVideoLikeCntPre(t *testing.T) {
	timeNow := strconv.FormatInt(time.Now().Unix()/3600/24, 10)
	res, err := GetVideoLikeCntOld(21, timeNow)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func TestGetVideoCommentCnt(t *testing.T) {
	res, err := GetVideoCommentCnt(22)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func Test(t *testing.T) {
	for i := 14; i <= 22; i++ {
		res := initial.RedisClient.Del(fmt.Sprint(strconv.FormatInt(int64(i), 10), "_", "like"))
		fmt.Println(res)
	}
}
