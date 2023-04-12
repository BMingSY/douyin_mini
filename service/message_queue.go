package service

import (
	redis2 "github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"tiktok/database/mysql"
	"tiktok/database/redis"
	"tiktok/utils"
)

func Consumer() {
	for {
		if !utils.TB.Allow() {
			// 如果没有容量，就暂缓处理MySQL
			continue
		}
		//获取消息
		val, err := redis.GetMessage()
		if err != nil {
			if err == redis2.Nil {
				continue
			}
			log.Errorf("Consumer get message error:[%v]", err)
			continue
		}
		arr := strings.Split(val[1], ":")
		videoId, err := strconv.ParseInt(arr[0], 10, 64)
		if err != nil {
			log.Errorf("Consumer parse videoId int error:[%v]", err)
			continue
		}
		count, err := strconv.ParseInt(arr[1], 10, 64)
		if err != nil {
			log.Errorf("Consumer parse count int error:[%v]", err)
			continue
		}
		if err := mysql.UpdateEventCountLikeCnt(videoId, int(count)); err != nil {
			log.Warnf("Consumer update event count like cnt warning:[%v]", err)
		}
	}
}
