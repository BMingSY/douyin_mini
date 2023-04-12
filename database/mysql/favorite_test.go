package mysql

import (
	"fmt"
	"testing"
	"tiktok/initial"
	"time"
)

func init() {
	_ = initial.InitConfig("../../config/config.yaml")
	_ = initial.InitMySQL()
}

func TestQueryLikeNumberByVideoId(t *testing.T) {
	ans, err := QueryLikeNumberByVideoId(1)
	if err != nil {
		panic(err)
	}
	println(ans)
}
func TestCreateLike(t *testing.T) {
	if err := CreateLike(1, 1); err != nil {
		panic(err)
	}
}

func TestDeleteLike(t *testing.T) {
	if err := DeleteLike(1, 1); err != nil {
		panic(err)
	}
}

func TestQueryListFavoriteByUserId(t *testing.T) {

	res, err := QueryListFavoriteByUserId(111)
	if err != nil {
		panic(err)
	}
	for _, i := range res {
		println(i.VideoId)
	}
}

func TestIsFavor(t *testing.T) {
	res, err := IsFavor(111, 1)
	if err != nil {
		panic(err)
	}
	println(res)
}
func TestGetAllLikeVideoId(t *testing.T) {
	res, err := GetAllLikeVideoId()
	if err != nil {
		panic(err)
	}
	for _, i := range res {
		fmt.Println(i)
	}
}

func TestQueryFromToCntByVideoId(t *testing.T) {
	res, err := QueryFromToCntByVideoId(21,
		time.Date(2022, time.February, 7, 20, 28, 46, 0, time.Local),
		time.Now())
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
