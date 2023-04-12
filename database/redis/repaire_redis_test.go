package redis

import (
	"testing"
	"tiktok/initial"
)

func init() {
	_ = initial.InitConfig("../../config/config.yaml")
	_ = initial.InitRedis()
	_ = initial.InitMySQL()
}
func TestRepair(t *testing.T) {
	if err := Repair(); err != nil {
		panic(err)
	}
}
