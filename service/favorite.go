package service

import (
	"sync"
	"tiktok/database/mysql"
	"tiktok/database/redis"

	log "github.com/sirupsen/logrus"
)

type FavoriteRequest struct {
	UserId  int64
	VideoId int64
}

type FavoriteServiceInterface interface {
	IsFavor(userId, videoId int64) (bool, error)
	CreateLike(req FavoriteRequest) error
	DeleteLike(req FavoriteRequest) error
	QueryFavoriteList(userId int64) ([]int64, error)
	QueryFavoriteNumberByVideoId(req int64) (int, error)
}

type FavoriteService struct {
}

var (
	favoriteService *FavoriteService
	favorOnce       sync.Once
)

func Favorite() *FavoriteService {
	favorOnce.Do(
		func() {
			favoriteService = &FavoriteService{}
		})
	return favoriteService
}

// CreateLike 点赞
func (*FavoriteService) CreateLike(req FavoriteRequest) error {
	rows, err := mysql.CreateLike(req.VideoId, req.UserId)
	if err != nil {
		log.Errorf("favorite service create like error:[%v]", err)
		return err
	}
	if rows != 1 {
		log.Warnf("favorite service create like warning:rowsAffected is not 1")
		return nil
	}
	//redis incr 1
	res, timeNow, err := redis.IncrVideoLikeCnt(req.VideoId, 1)
	if err != nil {
		log.Warnf("favorite service create like warning[%v]", err)
		// mysql delete
		rows, err1 := mysql.DeleteLike(req.VideoId, req.UserId)
		if err1 != nil {
			log.Errorf("favorite service create like error:[%v]", err)
			return err1
		}
		if rows != 1 {
			log.Warnf("favorite service delete like warning:rowsAffected is not 1")
			return err
		}
		return err
	}
	// 如果是新的一天 归档旧数据到mysql
	if res == 1 {
		// 获取所有旧的一天的数据
		ans, err := redis.GetVideoLikeCntOld(req.VideoId, timeNow)
		if err != nil {
			log.Errorf("favorite service get old like cnt error:[%v]", err)
			return err
		}
		// 发送消息给消费者进行归档
		if err := redis.AddMessage(req.VideoId, ans); err != nil {
			// 归档失败，但是有可能第二天能够成功
			log.Warnf("favorite service update mysql warning:[%v]", err)
		}
	}
	return nil
}

// DeleteLike 取消点赞
func (*FavoriteService) DeleteLike(req FavoriteRequest) error {
	rows, err := mysql.DeleteLike(req.VideoId, req.UserId)
	if err != nil {
		log.Errorf("favorite service delete like error:[%v]", err)
		return err
	}
	if rows != 1 {
		log.Warnf("favorite service delete like warning:rowsAffected is not 1")
		return nil
	}
	//redis incr -1
	res, timeNow, err := redis.IncrVideoLikeCnt(req.VideoId, -1)
	if err != nil {
		log.Warnf("favorite service delete like warning[%v]", err)
		// mysql create
		rows, err1 := mysql.CreateLike(req.VideoId, req.UserId)
		if err1 != nil {
			log.Errorf("favorite service delete like error:[%v]", err)
			return err1
		}
		if rows != 1 {
			log.Warnf("favorite service delete like warning:rowsAffected is not 1")
			return err
		}
		return err
	}

	// 如果是新的一天 归档旧数据到mysql
	if res == -1 {
		// 获取所有旧的一天的数据
		ans, err := redis.GetVideoLikeCntOld(req.VideoId, timeNow)
		if err != nil {
			log.Errorf("favorite service get old like cnt error:[%v]", err)
			return err
		}
		// 发送消息给消费者进行归档
		if err := redis.AddMessage(req.VideoId, ans); err != nil {
			// 归档失败，但是有可能第二天能够成功
			log.Warnf("favorite service update mysql warning:[%v]", err)
		}
	}
	return nil
}

// QueryFavoriteList 查询点赞列表
func (*FavoriteService) QueryFavoriteList(userId int64) ([]int64, error) {
	res, err := mysql.QueryListFavoriteByUserId(userId)
	if err != nil {
		return []int64{}, err
	}

	r := make([]int64, len(res))
	for i := 0; i < len(res); i++ {
		r[i] = res[i].VideoId
	}
	return r, nil
}

// QueryFavoriteNumberByVideoId 查询视频点赞数量
func (*FavoriteService) QueryFavoriteNumberByVideoId(videoId int64) (int, error) {
	res, length, err := redis.GetVideoLikeCnt(videoId)
	if err != nil {
		log.Errorf("favorite service query favorite number error:[%v]", err)
		return 0, err
	}
	// 如果长度为1 查询是否丢失了昨天及之前的数据
	if length == 1 {
		ans, err := mysql.QueryEventCountLikeCntOld(videoId)
		if err != nil {
			log.Errorf("favorite service query mysql event count error:[%v]", err)
			return res, err
		}
		// 创建一个昨天及以前的数据给redis
		if err := redis.IncrVideoLikeCntPre(videoId, int64(ans)); err != nil {
			// 如果创建失败，那么可能下一次创建会成功
			log.Warnf("favorite service redis incr pre warning:[%v]", err)
			return res + ans, err
		}
		return res + ans, nil
	}
	return res, nil
}

// IsFavor 查询视频是否点赞
func (*FavoriteService) IsFavor(userId, videoId int64) (bool, error) {
	return mysql.IsFavor(userId, videoId)
}

// GetTotalFavorited 总的点赞数
func GetTotalFavorited(userId int64) (int64, error) {
	var totalFavorited int64 = 0
	video, err := GetVideoList(userId)
	if err != nil {
		return 0, err
	}
	for i := 0; i < len(video); i++ {
		favoriteCount, err := Favorite().QueryFavoriteNumberByVideoId(video[i].VideoId)
		if err != nil {
			return 0, err
		}
		totalFavorited += int64(favoriteCount)
	}
	return totalFavorited, nil
}
