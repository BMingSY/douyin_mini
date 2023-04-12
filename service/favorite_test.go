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

func TestFavoriteService_CreateLike(t *testing.T) {
	if err := Favorite().CreateLike(FavoriteRequest{
		UserId:  1,
		VideoId: 1,
	}); err != nil {
		panic(err)
	}
}
func TestFavoriteService_DeleteLike(t *testing.T) {
	if err := Favorite().DeleteLike(FavoriteRequest{
		UserId:  1,
		VideoId: 1,
	}); err != nil {
		panic(err)
	}
}
func TestFavoriteService_IsFavor(t *testing.T) {
	res, err := Favorite().IsFavor(1, 1)
	if err != nil {
		panic(err)
	}
	println(res)
}

func TestFavoriteService_QueryFavoriteList(t *testing.T) {
	res, err := Favorite().QueryFavoriteList(25)
	if err != nil {
		panic(err)
	}
	for _, i := range res {
		fmt.Printf("%v ", i)
	}
	fmt.Println()
}

func TestFavoriteService_QueryFavoriteNumberByVideoId(t *testing.T) {
	res, err := Favorite().QueryFavoriteNumberByVideoId(22)
	if err != nil {
		panic(err)
	}
	println(res)
}
