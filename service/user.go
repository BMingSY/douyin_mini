package service

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"sync"
	cfg "tiktok/config"
	"tiktok/database/mysql"
	"tiktok/utils"

	log "github.com/sirupsen/logrus"
)

type UserServiceImpl struct {
}

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserRegisterResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
	UserId     int64  `json:"user_id"`
	Token      string `json:"token"`
}

type UserResponse struct {
	Response
	User UserInfo `json:"user"`
}

type UserInfo struct {
	Id              int64  `json:"id,omitempty"`
	Name            string `json:"name,omitempty"`
	FollowCount     int64  `json:"follow_count,omitempty"`
	FollowerCount   int64  `json:"follower_count,omitempty"`
	IsFollow        bool   `json:"is_follow,omitempty"`
	Avatar          string `json:"avatar,omitempty"`
	BackgroundImage string `json:"background_image,omitempty"`
	Signature       string `json:"signature,omitempty"`
	TotalFavorited  int64  `json:"total_favorited,omitempty"`
	WorkCount       int64  `json:"work_count,omitempty"`
	FavoriteCount   int64  `json:"favorite_count,omitempty"`
}

var (
	userService     *UserServiceImpl
	userServiceOnce sync.Once
)

func UserServiceInstance() *UserServiceImpl {
	userServiceOnce.Do(
		func() {
			userService = &UserServiceImpl{}
		})
	return userService
}

func (*UserServiceImpl) RegisterUser(username string, password string) (*UserRegisterResponse, error) {
	if len(username) > 32 || len(password) > 32 {
		log.Info("用户名或密码过长")
		return &UserRegisterResponse{StatusCode: 1, StatusMsg: "用户名或密码过长"}, errors.New("用户名或密码过长")
	}
	// 检查用户名是否存在
	var usr mysql.User
	result, err := mysql.SearchUserByUsername(username, &usr)
	if err != nil {
		log.Error(err)
		return &UserRegisterResponse{StatusCode: 1, StatusMsg: "查找信息失败"}, errors.New("查找信息失败")
	}
	if result { // 用户名已注册
		log.Info("用户名已存在")
		return &UserRegisterResponse{StatusCode: 1, StatusMsg: "用户名已存在"}, errors.New("用户名已存在")
	} else { // 创建用户
		// 密码加密
		h := md5.New()
		h.Write([]byte(password))
		passwordPlus := hex.EncodeToString(h.Sum(nil))
		// 数据填入数据库
		usr = mysql.User{
			UserName:      username,
			PassWord:      passwordPlus,
			FollowCount:   0,
			FollowerCount: 0,
		}
		err := mysql.CreateNewUser(&usr)
		if err != nil {
			log.Error(err)
			return &UserRegisterResponse{StatusCode: 1, StatusMsg: "用户注册失败"}, errors.New("用户注册失败")
		}

		result, err = mysql.SearchUserByUsername(username, &usr)
		if err != nil {
			log.Error(err)
			return &UserRegisterResponse{StatusCode: 1, StatusMsg: "查找信息失败"}, errors.New("查找信息失败")
		}
		if !result {
			log.Errorf("用户注册失败")
			return &UserRegisterResponse{
				StatusCode: 1,
				StatusMsg:  "用户创建失败",
			}, errors.New("用户创建失败")
		}
		tokenString, err := utils.GenerateToken(username, int(usr.UserId), 0)
		if err != nil {
			e := mysql.DeleteUser(usr)
			if e != nil {
				log.Error(e)
			}
			log.Error(err)
			return &UserRegisterResponse{
				StatusCode: 1,
				StatusMsg:  "token获取失败",
			}, err
		}
		return &UserRegisterResponse{
			StatusCode: 0,
			StatusMsg:  "用户创建成功",
			UserId:     int64(usr.UserId),
			Token:      tokenString,
		}, nil
	}
}

func (*UserServiceImpl) Login(username string, password string) (*UserLoginResponse, error) {
	var usr mysql.User
	result, err := mysql.SearchUserByUsername(username, &usr)
	if err != nil {
		log.Error(err)
		return &UserLoginResponse{Response{StatusCode: 1, StatusMsg: "查找信息失败"}, 0, ""}, errors.New("查找信息失败")
	}
	if !result {
		log.Errorf("用户名不存在")
		return &UserLoginResponse{Response{StatusCode: 1, StatusMsg: "用户名不存在"}, 0, ""}, errors.New("用户名不存在")
	}
	h := md5.New()
	h.Write([]byte(password))
	passwordPlus := hex.EncodeToString(h.Sum(nil))
	if usr.PassWord != passwordPlus {
		log.Info("密码错误")
		return &UserLoginResponse{Response{StatusCode: 1, StatusMsg: "密码错误"}, 0, ""}, errors.New("密码错误")
	}

	tokenString, err := utils.GenerateToken(username, int(usr.UserId), 0)
	if err != nil {
		log.Error(err)
		return &UserLoginResponse{Response{StatusCode: 1, StatusMsg: "token获取失败"}, 0, ""}, err
	}
	return &UserLoginResponse{Response{StatusCode: 0, StatusMsg: "登录成功"}, int64(usr.UserId), tokenString}, nil
}

func (*UserServiceImpl) GetUserInfo(myselfId int, userId int) (*UserResponse, error) {
	var usr mysql.User
	result, err := mysql.SearchUserByUserId(userId, &usr)
	if err != nil {
		log.Error(err)
		return &UserResponse{Response{StatusCode: 1, StatusMsg: "获取信息失败"}, UserInfo{}}, nil
	}
	if !result {
		log.Errorf("获取信息失败")
		return &UserResponse{Response{StatusCode: 1, StatusMsg: "获取信息失败"}, UserInfo{}}, nil
	}
	// 需增加查询关注关系信息
	var flag bool
	if myselfId == userId {
		flag = false
	}
	// 查询 myself 是否关注 user 若是 flag 为 true
	flag, err = Relation().CheckIsFollow(int64(myselfId), int64(userId))
	if err != nil {
		log.Error(err)
		return &UserResponse{Response{StatusCode: 1, StatusMsg: "获取信息失败"}, UserInfo{}}, nil
	}

	// 头像
	bucketName := "photo"
	objectName := "avatar.jpg"
	config := cfg.Config.MinIO
	avatar := "http://" + config.Url + ":" + config.APIPort + "/" + bucketName + "/" + objectName

	// 背景图片
	bucketName = "photo"
	objectName = "background.gif"
	config = cfg.Config.MinIO
	background := "http://" + config.Url + ":" + config.APIPort + "/" + bucketName + "/" + objectName

	// 获赞数量
	totalFavorited, err := GetTotalFavorited(int64(userId))
	if err != nil {
		log.Error(err)
		return &UserResponse{Response{StatusCode: 1, StatusMsg: "获取信息失败"}, UserInfo{}}, nil
	}

	// 作品数量
	workList, err := GetVideoList(int64(userId))
	if err != nil {
		log.Error(err)
		return &UserResponse{Response{StatusCode: 1, StatusMsg: "获取信息失败"}, UserInfo{}}, nil
	}
	workCount := int64(len(workList))

	// 点赞数量
	favoriteList, err := favoriteService.QueryFavoriteList(int64(userId))
	if err != nil {
		log.Error(err)
		return &UserResponse{Response{StatusCode: 1, StatusMsg: "获取信息失败"}, UserInfo{}}, nil
	}
	favoriteCount := int64(len(favoriteList))

	user := UserInfo{
		Id:              int64(userId),
		Name:            usr.UserName,
		FollowCount:     usr.FollowCount,
		FollowerCount:   usr.FollowerCount,
		IsFollow:        flag,
		Avatar:          avatar,
		BackgroundImage: background,
		Signature:       "很有个性的个性签名",
		TotalFavorited:  int64(totalFavorited),
		WorkCount:       workCount,
		FavoriteCount:   favoriteCount,
	}
	return &UserResponse{Response{StatusCode: 0, StatusMsg: "get user info success"}, user}, nil
}
