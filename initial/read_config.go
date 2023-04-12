package initial

import (
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	cfg "tiktok/config"
)

func InitConfig(path string) error {
	viper.SetConfigType("yaml")
	log.Info(os.Getwd())
	viper.SetConfigFile(path)
	err := viper.ReadInConfig()
	if err != nil {
		log.Infof("read config err:[%v]", err)
		return err
	}
	err = viper.Unmarshal(&cfg.Config)
	log.Infof("%+v", cfg.Config.MySQL)
	if err != nil {
		log.Infof("unmarshal error : [%v]", err)
		return err
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Infof("config has been updated, msg:[%v]", e.Name)
	})
	return nil
}
