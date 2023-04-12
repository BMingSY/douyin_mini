package redis

import (
	"fmt"
	"strings"
	"testing"
	"tiktok/initial"
)

func init() {
	_ = initial.InitConfig("../../config/config.yaml")
	_ = initial.InitRedis()
}

func TestAddMessage(t *testing.T) {
	err := AddMessage(10, 5)
	if err != nil {
		panic(err)
	}
}

func TestGetMessage(t *testing.T) {
	val, err := GetMessage()
	if err != nil {
		panic(err)
	}

	fmt.Println(val)

	arr := strings.Split(val[1], ":")
	for kk, vv := range arr {
		fmt.Printf("kk->%v,vv->%v\n", kk, vv)
	}

}
