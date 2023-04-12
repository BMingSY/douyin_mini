package service

import (
	log "github.com/sirupsen/logrus"
	"sync"
	"tiktok/database/mysql"
	"tiktok/database/redis"
	"time"
)

type CommentCreateRequest struct {
	UserId  int64
	VideoId int64
	Text    string
}
type CommentListResponse struct {
	ID         int64
	UserId     int64
	Text       string
	CreateTime time.Time
}

type CommentServiceInterface interface {
	CreateComment(req CommentCreateRequest) (int64, time.Time, error)
	DeleteComment(ID int64) error
	QueryCommentList(videoId int64) ([]CommentListResponse, error)
	QueryCommentNumberByVideoId(videoId int64) (int, error)
}

type CommentService struct {
}

var (
	commentService *CommentService
	commentOnce    sync.Once
)

func Comment() *CommentService {
	commentOnce.Do(
		func() {
			commentService = &CommentService{}
		})
	return commentService
}

// CreateComment 评论
func (*CommentService) CreateComment(req CommentCreateRequest) (int64, time.Time, error) {
	com := mysql.Comment{
		UserId:  req.UserId,
		VideoId: req.VideoId,
		Text:    req.Text,
	}

	ID, CreateTime, err := mysql.CreateComment(com)
	if err != nil {
		return ID, CreateTime, err
	}

	//redis incr 1
	if err := redis.IncrVideoCommentCnt(req.VideoId, 1); err != nil {
		log.Warnf("comment service create comment warning:[%v]", err)
		// mysql delete
		if err := mysql.DeleteComment(ID); err != nil {
			log.Errorf("comment service create comment error:[%v]", err)
			return ID, CreateTime, err
		}
		return ID, CreateTime, err
	}
	return ID, CreateTime, nil
}

// DeleteComment 删除评论
func (*CommentService) DeleteComment(ID int64) error {
	if err := mysql.DeleteComment(ID); err != nil {
		log.Errorf("comment service delete comment error:[%v]", err)
		return err
	}
	//redis incr -1
	if err := redis.IncrVideoCommentCnt(ID, -1); err != nil {
		log.Warnf("comment service delete comment warning:[%v]", err)
		// mysql create
		if err := mysql.CommentSetDeleteAtNil(ID); err != nil {
			log.Errorf("comment service delete comment error:[%v]", err)
			return err
		}
		return err
	}
	return nil
}

// QueryCommentList 查询评论列表
func (*CommentService) QueryCommentList(videoId int64) ([]CommentListResponse, error) {
	res, err := mysql.QueryAllCommentListByVideoID(videoId)
	if err != nil {
		return []CommentListResponse{}, err
	}
	r := make([]CommentListResponse, len(res))
	for i := 0; i < len(res); i++ {
		r[i].ID = res[i].ID
		r[i].UserId = res[i].UserId
		r[i].Text = res[i].Text
		r[i].CreateTime = res[i].CreateTime
	}
	return r, nil
}

// QueryCommentNumberByVideoId 查询视频评论数量
func (*CommentService) QueryCommentNumberByVideoId(videoId int64) (int, error) {
	return redis.GetVideoCommentCnt(videoId)
}
