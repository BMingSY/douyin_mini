package mysql

import (
	"gorm.io/gorm"
	"tiktok/initial"
	"time"
)

type Favorite struct {
	ID         int64
	UserId     int64
	VideoId    int64
	DeleteTime gorm.DeletedAt
}

func (Favorite) TableName() string {
	return "favorite"
}

// QueryLikeNumberByVideoId 查询视频id点赞数量
func QueryLikeNumberByVideoId(videoId int64) (int64, error) {
	var ans int64
	if err := initial.Database.Model(&Favorite{}).
		Where("video_id = ?", videoId).
		Count(&ans).Error; err != nil {
		return ans, err
	}
	return ans, nil
}

// CreateLike 点赞
func CreateLike(videoId int64, userId int64) (int64, error) {

	res := initial.Database.Create(&Favorite{
		UserId:  userId,
		VideoId: videoId})
	if res.Error != nil {
		return res.RowsAffected, res.Error
	}
	return res.RowsAffected, nil
}

// DeleteLike 取消点赞
func DeleteLike(videoId int64, userId int64) (int64, error) {
	res := initial.Database.
		Where("user_id = ? AND video_id = ?", userId, videoId).
		Delete(&Favorite{})
	if res.Error != nil {
		return res.RowsAffected, res.Error
	}
	return res.RowsAffected, nil
}

// QueryListFavoriteByUserId 根据用户id获取点赞列表
func QueryListFavoriteByUserId(userId int64) ([]Favorite, error) {
	var favoriteList []Favorite
	if err := initial.Database.Where("user_id = ?", userId).
		Order("create_time DESC").
		Find(&favoriteList).Error; err != nil {
		return favoriteList, err
	}
	return favoriteList, nil
}

// IsFavor 查看userId是否已经对videoId点赞
func IsFavor(userId int64, videoId int64) (bool, error) {
	res := initial.Database.
		Where("user_id = ? AND video_id = ?", userId, videoId).
		Find(&Favorite{})
	if res.Error != nil {
		return false, res.Error
	}
	if res.RowsAffected > 0 {
		return true, nil
	}
	return false, nil
}

// GetAllLikeVideoId 获取点赞库中所有得videoId
func GetAllLikeVideoId() ([]int64, error) {
	var favoriteList []Favorite
	res := initial.Database.Select("DISTINCT video_id").Find(&favoriteList)
	if res.Error != nil {
		return nil, res.Error
	}
	videoIdList := make([]int64, len(favoriteList))
	for i := 0; i < len(favoriteList); i++ {
		videoIdList[i] = favoriteList[i].VideoId
	}
	return videoIdList, nil
}

// QueryFromToCntByVideoId 查询某个时间段内的点赞数量
func QueryFromToCntByVideoId(videoId int64, fromTime time.Time, toTime time.Time) (int, error) {
	var ans int64
	if err := initial.Database.Model(&Favorite{}).
		Where("video_id = ? AND create_time >= ? AND create_time < ?", videoId, fromTime, toTime).
		Count(&ans).Error; err != nil {
		return int(ans), err
	}
	return int(ans), nil
}
