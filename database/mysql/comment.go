package mysql

import (
	"gorm.io/gorm"
	"tiktok/initial"
	"time"
)

type Comment struct {
	ID         int64
	UserId     int64
	VideoId    int64
	Text       string
	CreateTime time.Time
	DeleteTime gorm.DeletedAt
}

func (Comment) TableName() string {
	return "comment"
}

// CreateComment 创建评论
func CreateComment(comment Comment) (int64, time.Time, error) {
	comment.CreateTime = time.Now()
	if err := initial.Database.Create(&comment).Error; err != nil {
		return 0, time.Now(), err
	}
	return comment.ID, comment.CreateTime, nil
}

// DeleteComment 删除评论
func DeleteComment(ID int64) error {
	if err := initial.Database.Where("ID = ?", ID).
		Delete(&Comment{}).Error; err != nil {
		return err
	}
	return nil
}

// QueryAllCommentNumberByVideoID 获取video_id的评论数
func QueryAllCommentNumberByVideoID(videoId int64) (int64, error) {
	var ans int64
	if err := initial.Database.Model(&Comment{}).
		Where("video_id = ?", videoId).
		Count(&ans).Error; err != nil {
		return ans, err
	}
	return ans, nil
}

// QueryAllCommentListByVideoID 获取video_id的评论列表
func QueryAllCommentListByVideoID(videoId int64) ([]Comment, error) {
	var commentList []Comment
	if err := initial.Database.Where("video_id = ?", videoId).
		Order("create_time DESC").
		Find(&commentList).Error; err != nil {
		return commentList, err
	}
	return commentList, nil
}

// GetAllCommentVideoId 获取点赞库中所有的videoId
func GetAllCommentVideoId() ([]int64, error) {
	var commentList []Comment
	res := initial.Database.Find(&commentList)
	if res.Error != nil {
		return nil, res.Error
	}
	videoIdList := make([]int64, len(commentList))
	for i := 0; i < len(commentList); i++ {
		videoIdList[i] = commentList[i].VideoId
	}
	return videoIdList, nil
}

// CommentSetDeleteAtNil 恢复删除的数据
func CommentSetDeleteAtNil(ID int64) error {
	if err := initial.Database.Model(&Comment{}).Unscoped().
		Where("ID = ?", ID).
		Update("delete_time", nil).
		Error; err != nil {
		return err
	}
	return nil
}
