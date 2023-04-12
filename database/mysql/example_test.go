package mysql

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"tiktok/initial"
)

func TestExample(t *testing.T) {
	_ = initial.InitConfig("../../config/config.yaml")
	_ = initial.InitMySQL()
	name, err := GetNameByID(66)
	assert.NoError(t, err)
	assert.Equal(t, name, "lcl fan club")

}
