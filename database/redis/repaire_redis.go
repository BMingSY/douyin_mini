package redis

import (
	log "github.com/sirupsen/logrus"
	"tiktok/database/mysql"
	"time"
)

// Repair if redis num wrong ,run this function
// it will sync mysql's cnt
func Repair() error {
	// repair like_cnt
	videoIdLikeList, err := mysql.GetAllLikeVideoId()
	// fmt.Println(len(videoIdLikeList))
	if err != nil {
		log.Errorf("repair get video id list error:[%v]", err)
		return err
	}
	for i := 0; i < len(videoIdLikeList); i++ {
		// 获取最后归档时间
		//fromTime, err := mysql.QueryEventCountTime(videoIdLikeList[i])
		//if err != nil {
		//	log.Errorf("repaire get event count time error:[%v]", err)
		//	return err
		//}
		fromTime := time.Date(2022, time.February, 18, 0, 0, 0, 0, time.Local)
		// 指定到redis启动时间
		toTime := time.Date(2023, time.February, 29, 0, 0, 0, 0, time.Local)
		// 获取期间丢失的点赞数
		cnt, err := mysql.QueryFromToCntByVideoId(videoIdLikeList[i], fromTime, toTime)
		if err != nil {
			log.Errorf("repair get event count error:[%v]", err)
			return err
		}
		// 加给昨天及以前的redis
		if err := IncrVideoLikeCntPre(videoIdLikeList[i], int64(cnt)); err != nil {
			log.Errorf("repair redis like cnt pre error:[%v]", err)
			return err
		}
	}
	// repair comment_cnt
	videoIdCommentList, err := mysql.GetAllCommentVideoId()
	if err != nil {
		log.Errorf("repair get video id list error:[%v]", err)
		return err
	}
	for i := 0; i < len(videoIdCommentList); i++ {
		if cnt, err := GetVideoCommentCnt(videoIdCommentList[i]); cnt == 0 && err == nil {
			if likeCnt, err := mysql.QueryAllCommentNumberByVideoID(videoIdCommentList[i]); err == nil {
				if err := IncrVideoCommentCnt(videoIdCommentList[i], likeCnt); err != nil {
					log.Errorf("repair incr video comment cnt error:[%v]", err)
					return err
				}
			} else {
				log.Errorf("repair get video comment cnt error:[%v]", err)
				return err
			}
		}
	}
	return nil
}
