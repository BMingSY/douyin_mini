package mysql

import (
	"gorm.io/gorm"
	"tiktok/initial"
)

type Relation struct {
	ID           int64
	UserID       int64
	FollowUserId int64
	DeletedAt    gorm.DeletedAt
}

func (Relation) TableName() string {
	return "follow"
}

// create follow 添加关注
func CreateFollow(userId int64, followUserId int64) error {
	return initial.Database.Create(&Relation{UserID: userId, FollowUserId: followUserId}).Error
}

// delete follow 删除关注
func DeleteFollow(userId int64, followUserId int64) error {
	return initial.Database.Where("user_id = ? and follow_user_id = ?", userId, followUserId).Delete(&Relation{}).Error
}

// QueryFollowList 登录用户关注的所有用户列表
func QueryFollowList(userId int64) ([]Relation, error) {
	var followList []Relation
	err := initial.Database.Where("user_id = ? and deleted_at is null", userId).Order("created_at DESC").Find(&followList).Error
	return followList, err
}

// QueryFollowerList 所有关注登录用户的粉丝列表
func QueryFollowerList(userId int64) ([]Relation, error) {
	var followerList []Relation
	err := initial.Database.Where("follow_user_id = ? and deleted_at is null", userId).Order("created_at DESC").Find(&followerList).Error
	return followerList, err
}

// 检查 userid 是否 关注了 followUserId
func CheckIsFollow(userId int64, followUserId int64) (bool, error) {
	res := initial.Database.Where("user_id = ? and follow_user_id = ? and deleted_at is null", userId, followUserId).Find(&Relation{})
	if res.Error != nil {
		return false, res.Error
	}
	if res.RowsAffected > 0 {
		return true, nil
	}
	return false, nil
}
