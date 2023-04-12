package initial

import (
	"fmt"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	cfg "tiktok/config"
)

var RedisClient *redis.Client

func InitRedis() error {
	config := cfg.Config.Redis
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Host, config.Port),
		Password: fmt.Sprintf("%s", config.Password),
	})
	if _, err := RedisClient.Ping().Result(); err != nil {
		log.Errorf("init redis error:[%v]", err)
		return err
	}
	return nil
}
