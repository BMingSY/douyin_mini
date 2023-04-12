package mysql

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"tiktok/initial"
	"time"
)

type EventCount struct {
	VideoId    int64
	LikeCnt    int
	CreateTime time.Time
	UpdateTime time.Time
	DeleteTime gorm.DeletedAt
}

func (EventCount) TableName() string {
	return "event_count"
}

// UpdateEventCountLikeCnt 更新点赞计数信息
func UpdateEventCountLikeCnt(videoId int64, cnt int) error {
	db := initial.Database.Model(&EventCount{})

	if err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "video_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"like_cnt"}),
	}).Create(&EventCount{
		VideoId:    videoId,
		LikeCnt:    cnt,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}).Error; err != nil {
		log.Errorf("upsert event count error:[%v]", err)
		return err
	}
	
	return nil
}

// QueryEventCountLikeCntOld 查询旧点赞数
func QueryEventCountLikeCntOld(videoId int64) (int, error) {
	db := initial.Database.Model(&EventCount{})
	var ec EventCount
	result := db.Where("video_id = ?", videoId).Find(&ec)
	if result.Error != nil {
		log.Errorf("find event count like cnt err:[%v]", result.Error)
		return 0, result.Error
	}
	// 有则返回点赞数
	if result.RowsAffected > 0 {
		return ec.LikeCnt, nil
	}
	// 无则创建
	if err := db.Create(&EventCount{
		VideoId:    videoId,
		LikeCnt:    0,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}).Error; err != nil {
		log.Errorf("create event count like cnt err:[%v]", err)
		return 0, err
	}
	return 0, nil
}

// QueryEventCountTime 获取归档时间
func QueryEventCountTime(videoId int64) (time.Time, error) {
	db := initial.Database.Model(&EventCount{})
	var ec EventCount
	result := db.Where("video_id = ?", videoId).Find(&ec)
	if result.Error != nil {
		log.Errorf("find event count like cnt err:[%v]", result.Error)
		return time.Now(), result.Error
	}
	// 有
	if result.RowsAffected > 0 {
		return ec.UpdateTime, nil
	}
	// 无则
	return time.Date(2006, time.January, 2, 15, 4, 5, 0, time.Local), nil
}
