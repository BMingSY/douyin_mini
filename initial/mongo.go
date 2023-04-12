package initial

import (
	"context"
	"fmt"
	"github.com/qiniu/qmgo"
	log "github.com/sirupsen/logrus"
	cfg "tiktok/config"
)

var MessageCol *qmgo.Collection

func InitMongo() error {
	ctx := context.Background()
	log.Info(*cfg.Config.Mongo)
	url := fmt.Sprintf("mongodb://%s/mydb", cfg.Config.Mongo.Host)
	client, err := qmgo.NewClient(ctx, &qmgo.Config{Uri: url})
	if err != nil {
		panic(err)
	}
	db := client.Database(cfg.Config.Mongo.Database)
	MessageCol = db.Collection(cfg.Config.Mongo.Collection)
	return nil
}
