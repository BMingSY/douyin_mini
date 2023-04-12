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

func TestUpdateEventCountAttr(t *testing.T) {
	err := UpdateEventCountLikeCnt(5, 2)
	if err != nil {
		panic(err)
	}
}
func TestQueryEventCountTime(t *testing.T) {
	res, err := QueryEventCountTime(22)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
func TestUpdateEventCountLikeCnt(t *testing.T) {
	if err := UpdateEventCountLikeCnt(10, 0); err != nil {
		panic(err)
	}
}
