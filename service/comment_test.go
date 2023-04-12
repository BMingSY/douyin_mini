package service

import (
	"fmt"
	"testing"
	"tiktok/initial"
)

func init() {
	_ = initial.InitConfig("../config/config.yaml")
	_ = initial.InitMySQL()
	_ = initial.InitRedis()
}

func TestCommentService_CreateComment(t *testing.T) {
	req := CommentCreateRequest{
		UserId:  1,
		VideoId: 1,
		Text:    "太好啦~~~~~",
	}
	ID, CreateTime, err := Comment().CreateComment(req)
	if err != nil {
		panic(err)
	}
	fmt.Println(ID, CreateTime.Format("01-02"))
}

func TestCommentService_DeleteComment(t *testing.T) {
	err := Comment().DeleteComment(42)
	if err != nil {
		panic(err)
	}
}

func TestCommentService_QueryCommentList(t *testing.T) {
	comList, err := Comment().QueryCommentList(111)
	if err != nil {
		panic(err)
	}
	for _, i := range comList {
		println(i.CreateTime.Format("01-02"), i.ID, i.Text)
	}
}

func TestCommentService_QueryCommentNumberByVideoId(t *testing.T) {
	res, err := Comment().QueryCommentNumberByVideoId(22)
	if err != nil {
		panic(err)
	}
	println(res)
}
