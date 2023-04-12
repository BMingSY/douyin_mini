package service

import (
	log "github.com/sirupsen/logrus"
	"sync"
	"tiktok/database/mysql"
)

type RelationRequest struct {
	UserID       int64
	FollowUserId int64
}

type RelationServiceInterface interface {
	CreateFollow(req RelationRequest) error
	DeleteFollow(req RelationRequest) error
	QueryFollowList(userId int64) ([]UserInfo, error)
	QueryFollowerList(userId int64) ([]UserInfo, error)
	QueryFriendList(userId int64) ([]UserInfo, error)
}

type RelationService struct {
}

var (
	relationService *RelationService
	relationOnce    sync.Once
)

func Relation() *RelationService {
	relationOnce.Do(
		func() {
			relationService = &RelationService{}
		})
	return relationService
}

// CreateFollow
func (*RelationService) CreateFollow(req RelationRequest) error {
	return mysql.CreateFollow(req.UserID, req.FollowUserId)
}

// DeleteFollow
func (*RelationService) DeleteFollow(req RelationRequest) error {
	return mysql.DeleteFollow(req.UserID, req.FollowUserId)
}

// QueryFollowList 登录用户关注的所有用户列表
func (*RelationService) QueryFollowList(userId int64) ([]UserInfo, error) {
	res, err := mysql.QueryFollowList(userId)
	if err != nil {
		return []UserInfo{}, err
	}

	ans := make([]UserInfo, len(res))
	for i := 0; i < len(res); i++ {
		userService := UserServiceInstance()
		response, err := userService.GetUserInfo(int(userId), int(res[i].FollowUserId))
		if err != nil { // error
			log.Errorf(" query follow List error [%#v]", err)
			continue
		}
		ans[i] = response.User
	}
	return ans, err
}

// QueryFollowerList 所有关注登录用户的粉丝列表
func (*RelationService) QueryFollowerList(userId int64) ([]UserInfo, error) {
	res, err := mysql.QueryFollowerList(userId)
	if err != nil {
		return []UserInfo{}, err
	}

	followerList := make([]UserInfo, len(res))
	for i := 0; i < len(res); i++ {
		userService := UserServiceInstance()
		response, err := userService.GetUserInfo(int(userId), int(res[i].UserID))
		if err != nil { // error
			log.Errorf(" query follower List error [%#v]", err)
			continue
		}
		followerList[i] = response.User

		log.Infof("%v %v\n", response.User.Id, response.User.Name)
	}
	return followerList, err
}

// 检查 userid 是否 关注了 followUserId
func (*RelationService) CheckIsFollow(userId int64, followUserId int64) (bool, error) {
	return mysql.CheckIsFollow(userId, followUserId)
}

// 进行查询所有的 firend
func (*RelationService) QueryFriendList(userId int64) ([]UserInfo, error) {
	res, err := mysql.QueryFollowerList(userId)
	if err != nil {
		return []UserInfo{}, err
	}

	followerList := make([]UserInfo, len(res))
	ct := 0
	for i := 0; i < len(res); i++ {
		userService := UserServiceInstance()
		response, err := userService.GetUserInfo(int(userId), int(res[i].UserID))
		if err != nil { // error
			log.Errorf(" query follower List error [%#v]", err)
			continue
		}
		followerList[i] = response.User
		if followerList[i].IsFollow {
			ct++
		}
	}

	firendList := make([]UserInfo, ct)
	j := 0
	for i := 0; i < len(res); i++ {
		if followerList[i].IsFollow {
			firendList[j] = followerList[i]
			j++
		}
	}
	return firendList, err
}
