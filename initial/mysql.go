package initial

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	cfg "tiktok/config"
)

var Database *gorm.DB

func InitMySQL() error {
	config := cfg.Config.MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%s",
		config.Username, config.Password, config.Host, config.Port, config.DBName, config.Timeout)
	var err error
	Database, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return err
}
