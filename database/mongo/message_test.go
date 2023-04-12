package mongo

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)
import "tiktok/initial"

func TestMessage(t *testing.T) {
	err := initial.InitConfig("../../config/config.yaml")
	err = initial.InitMongo()
	assert.NoError(t, err)
	err = InsertOneMsg(context.Background(), &Messages{
		FromUserID: 1111,
		ToUserID:   2322,
		Content:    "test",
		CreateTime: time.Now().Unix(),
	})
	assert.NoError(t, err)
}
