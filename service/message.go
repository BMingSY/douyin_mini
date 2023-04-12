package service

import (
	"context"
	"sync"
	"tiktok/database/mongo"
	"time"
)

type MessageService struct {
}

var (
	messageService *MessageService
	messageOnce    sync.Once
)

func Message() *MessageService {
	messageOnce.Do(
		func() {
			messageService = &MessageService{}
		})
	return messageService
}

func (m *MessageService) AddMessage(userID, toUserID int64, content string) error {
	err := mongo.InsertOneMsg(context.Background(), &mongo.Messages{
		FromUserID: userID,
		ToUserID:   toUserID,
		Content:    content,
		CreateTime: time.Now().Unix(),
		ID:         time.Now().UnixNano(),
	})
	if err != nil {
		return err
	}
	return nil
}
func (m *MessageService) GetMessage(userID, toUserID, preMsgTims int64) ([]mongo.Messages, error) {

	res, err := mongo.QueryMsg(context.Background(), userID, toUserID, preMsgTims)
	if err != nil {
		return nil, err
	}
	res1, err := mongo.QueryMsg(context.Background(), toUserID, userID, preMsgTims)
	if err != nil {
		return nil, err
	}
	res = append(res, res1...)
	return res, nil
}
