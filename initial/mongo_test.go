package initial

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitMongo(t *testing.T) {
	err := InitConfig("../config/config.yaml")
	assert.NoError(t, err)
	err = InitMongo()
	assert.NoError(t, err)
}

func TestInitMySQL(t *testing.T) {
	err := InitConfig("../config/config.yaml")
	assert.NoError(t, err)
	err = InitMySQL()
	assert.NoError(t, err)
}
