package mysql

import (
	"tiktok/initial"
	"time"
)

// type DeletedAt sql.NullTime

type User struct {
	// gorm.Model
	UserId        int64     `gorm:"column:user_id;primarykey"`
	UserName      string    `gorm:"column:username;type:varchar(32)not null;comment:用户名"`
	PassWord      string    `gorm:"column:password;type:varchar(255)not null;comment:密码"`
	FollowCount   int64     `gorm:"column:follow_count;default:0;comment:关注总数"`
	FollowerCount int64     `gorm:"column:follower_count;default:0;comment:粉丝总数"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
	// DeletedAt     DeletedAt `gorm:"column:delete_at"`
}

func SearchUserByUsername(username string, usr *User) (bool, error) {
	result := initial.Database.Where("username = ?", username).Find(&usr)
	return result.RowsAffected != 0, result.Error
}

func SearchUserByUserId(userId int, usr *User) (bool, error) {
	result := initial.Database.Where("user_id = ?", userId).Find(&usr)
	return result.RowsAffected != 0, result.Error
}

func CreateNewUser(usr *User) error {
	err := initial.Database.Create(&usr).Error
	return err
}

func DeleteUser(usr User) error {
	err := initial.Database.Delete(&usr).Error
	return err
}
