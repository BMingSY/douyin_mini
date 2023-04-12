package redis

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"tiktok/initial"
	"time"
)

// AddMessage 向消息队列中加一条消息
func AddMessage(videoId int64, count int) error {
	_, err := initial.RedisClient.LPush(initial.REDIS_LIST_KEY, fmt.Sprintf("%v:%v", videoId, count)).Result()
	if err != nil {
		log.Errorf("redis add message error:[%v]", err)
		return err
	}
	return nil
}

// GetMessage 获取消息
func GetMessage() ([]string, error) {
	// 阻塞5秒
	value, err := initial.RedisClient.BRPop(5*time.Second, initial.REDIS_LIST_KEY).Result()
	if err != nil {
		return value, err
	}
	return value, nil
}
