package mysql

import (
	"fmt"
	"testing"
	"tiktok/initial"
)

func init() {
	_ = initial.InitConfig("../../config/config.yaml")
	_ = initial.InitMySQL()
}

func TestCreateComment(t *testing.T) {

	comment := Comment{VideoId: 1, UserId: 2, Text: "你牌的太好了啦~"}
	ID, CreateTime, err := CreateComment(comment)
	if err != nil {
		panic(err)
	}
	println(ID, CreateTime.Format("01-02"))
}
func TestDeleteComment(t *testing.T) {

	if err := DeleteComment(10); err != nil {
		panic(err)
	}
}

func TestQueryAllCommentListByVideoID(t *testing.T) {
	commentList, err := QueryAllCommentListByVideoID(1)
	if err != nil {
		panic(err)
	}
	for _, i := range commentList {
		println(i.UserId, i.VideoId)
		println(i.Text)
	}
}

func TestQueryAllCommentNumberByVideoID(t *testing.T) {
	ans, err := QueryAllCommentNumberByVideoID(1)
	if err != nil {
		panic(err)
	}
	println(ans)
}
func TestGetAllCommentVideoId(t *testing.T) {
	res, err := GetAllCommentVideoId()
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(res); i++ {
		fmt.Println(res[i])
	}
}

func TestCommentSetDeleteAtNil(t *testing.T) {
	if err := CommentSetDeleteAtNil(1); err != nil {
		panic(err)
	}
}
