package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"tiktok/initial"
)

type Messages struct {
	ID         int64  `json:"id"`
	FromUserID int64  `json:"from_user_id"`
	ToUserID   int64  `json:"to_user_id"`
	Content    string `json:"content"`
	CreateTime int64  `json:"create_time"`
}

func InsertOneMsg(ctx context.Context, msg *Messages) error {
	_, err := initial.MessageCol.InsertOne(ctx, msg)
	return err
}

func QueryMsg(ctx context.Context, userID, toUserID, preMsgTime int64) ([]Messages, error) {
	var res []Messages
	//log.Info("userID ", userID, " toUserID ", toUserID, " premsg time ", preMsgTime)
	err := initial.MessageCol.Find(ctx, bson.M{
		"fromuserid": userID,
		"touserid":   toUserID,
		"createtime": bson.M{
			"$gt": preMsgTime,
		},
	}).All(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
