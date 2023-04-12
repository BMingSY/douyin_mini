package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"strconv"
	"tiktok/initial"
	"time"
)

// IncrVideoLikeCnt 增加点赞计数
func IncrVideoLikeCnt(videoId, incr int64) (int64, string, error) {
	// 时间戳
	timeNow := strconv.FormatInt(time.Now().Unix()/3600/24, 10)
	key := fmt.Sprint(strconv.FormatInt(videoId, 10), "_", "like")
	res, err := initial.RedisClient.HIncrBy(key, timeNow, incr).Result()

	if err != nil {
		log.Errorf("redis incr video like cnt error:[%v]", err)
		return res, timeNow, err
	}
	return res, timeNow, nil
}

// IncrVideoLikeCntPre 增加前一天点赞计数
func IncrVideoLikeCntPre(videoId, incr int64) error {
	// 时间戳
	timeNow := strconv.FormatInt(time.Now().Unix()/3600/24-1, 10)
	key := fmt.Sprint(strconv.FormatInt(videoId, 10), "_", "like")
	_, err := initial.RedisClient.HIncrBy(key, timeNow, incr).Result()
	if err != nil {
		log.Errorf("redis incr video like cnt pre error:[%v]", err)
		return err
	}
	return nil
}

// IncrVideoCommentCnt 增加评论计数
func IncrVideoCommentCnt(videoId, incr int64) error {
	key := fmt.Sprint(strconv.FormatInt(videoId, 10), "_", "comment")
	_, err := initial.RedisClient.HIncrBy(key, "cnt", incr).Result()
	if err != nil {
		log.Errorf("redis incr video comment cnt error:[%v]", err)
		return err
	}
	return nil
}

// GetVideoLikeCnt 获取视频点赞数
func GetVideoLikeCnt(videoId int64) (int, int, error) {
	key := fmt.Sprint(strconv.FormatInt(videoId, 10), "_", "like")
	res, err := initial.RedisClient.HGetAll(key).Result()
	if err != nil {
		// 没有这个key
		if err == redis.Nil {
			return 0, 0, nil
		}
		log.Errorf("redis get video like cnt error:[%v]", err)
		return 0, 0, err
	}
	ans := 0

	for _, v := range res {
		t, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			log.Warnf("redis get video like cnt warning:[%v]", err)
		}
		ans += int(t)
	}
	return ans, len(res), nil
}

// GetVideoLikeCntOld 获取前几天的视频点赞数
func GetVideoLikeCntOld(videoId int64, timeNow string) (int, error) {
	key := fmt.Sprint(strconv.FormatInt(videoId, 10), "_", "like")
	res, err := initial.RedisClient.HGetAll(key).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		log.Errorf("redis get video like cnt old error:[%v]", err)
		return 0, err
	}

	var (
		timeOld = "0"
		ans     = 0
		sum     = 0
	)

	for k, v := range res {
		if k == timeNow {
			continue
		}
		if timeOld == "0" {
			timeOld = k
		}
		t, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			log.Warnf("redis get video like cnt old warning:[%v]", err)
		}
		ans += int(t)
		// 删除旧key
		if k != timeNow && k != timeOld {
			sum += int(t)
			initial.RedisClient.HDel(key, k)
		}
	}
	// 将剩余的归档给昨天
	initial.RedisClient.HIncrBy(key, timeOld, int64(sum))

	return ans, nil
}

// GetVideoCommentCnt 获取视频评论数量
func GetVideoCommentCnt(videoId int64) (int, error) {
	key := fmt.Sprint(strconv.FormatInt(videoId, 10), "_", "comment")
	result, err := initial.RedisClient.HGetAll(key).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		log.Errorf("redis get video comment cnt error:[%v]", err)
		return 0, err
	}
	if result["cnt"] == "" {
		return 0, nil
	}
	res, err := strconv.ParseInt(result["cnt"], 10, 64)
	if err != nil {
		log.Errorf("redis parse video comment cnt error:[%v]", err)
		return 0, err
	}
	return int(res), nil
}
