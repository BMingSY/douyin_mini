package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"tiktok/initial"
	"tiktok/service"
)

func main() {
	//go service.RunMessageServer()
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "01-02 15:04:05",
	})
	if err := initial.InitConfig("./config/config.yaml"); err != nil {
		panic(err)
	}
	if err := initial.InitMySQL(); err != nil {
		panic(err)
	}
	if err := initial.InitMinIO(); err != nil {
		panic(err)
	}
	if err := initial.InitRedis(); err != nil {
		panic(err)
	}
	if err := initial.InitMongo(); err != nil {
		panic(err)
	}

	initial.InitTokenBucket()

	go service.Consumer()

	r := gin.Default()

	initRouter(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
