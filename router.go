package main

import (
	"tiktok/controller"
	"tiktok/middleware"

	"github.com/gin-gonic/gin"
)

func initRouter(r *gin.Engine) {
	// public directory is used to serve static resources
	r.Static("/static", "./public")

	apiRouter := r.Group("/douyin")

	// basic apis
	apiRouter.GET("/feed/",  middleware.JwtAuth(), controller.Feed)
	apiRouter.GET("/user/", middleware.JwtAuth(), controller.UserInfo)
	apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/login/", controller.Login)
	// 视频投稿
	publish := apiRouter.Group("/publish")
	{
		publish.Use(middleware.JwtAuth())
		publish.POST("/action/", controller.Publish)
		publish.GET("/list/", controller.PublishList)
	}

	// extra apis - I
	favorite := apiRouter.Group("/favorite")
	{
		favorite.Use(middleware.JwtAuth())
		favorite.POST("/action/", controller.FavoriteAction)
		favorite.GET("/list/", controller.FavoriteList)
	}
	apiRouter.POST("/comment/action/", middleware.JwtAuth(), controller.CommentAction)
	apiRouter.GET("/comment/list/", controller.CommentList)

	// extra apis - II
	relation := apiRouter.Group("/relation")
	{
		relation.Use(middleware.JwtAuth())
		relation.POST("/action/", controller.RelationAction)
		relation.GET("/follow/list/", controller.FollowList)
		relation.GET("/follower/list/", controller.FollowerList)
		relation.GET("/friend/list/", controller.FriendList)
	}

	message := apiRouter.Group("/message")
	{
		message.Use(middleware.JwtAuth())
		message.GET("/chat/", controller.MessageChat)
		message.POST("/action/", controller.MessageAction)
	}
}
