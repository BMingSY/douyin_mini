package mysql

import (
	"testing"
)

func TestSearchUser(t *testing.T) {
	var usr User
	SearchUserByUsername("test", &usr)
	println("test")
	// println(usr.UserName, "test")
}
