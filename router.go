package main

import (
	"github.com/gin-gonic/gin"
	"github.com/herrluk/douyin/controller"
)

func initRouter(r *gin.Engine) {
	// public directory is used to serve static resources
	// 配置静态web目录，第一个参数表示路由，第二个参数表示映射目录
	r.Static("/static", "./public")

	// 配置路由分组 /douyin
	apiRouter := r.Group("/douyin")

	// basic apis
	// 分组下的路由
	apiRouter.GET("/feed/", controller.Feed)
	apiRouter.GET("/user/", controller.UserInfo)
	apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/login/", controller.Login)
	apiRouter.POST("/publish/action/", controller.Publish)
	apiRouter.GET("/publish/list/", controller.PublishList)

	// extra apis - I
	apiRouter.POST("/favorite/action/", controller.FavoriteAction)
	apiRouter.GET("/favorite/list/", controller.FavoriteList)
	apiRouter.POST("/comment/action/", controller.CommentAction)
	apiRouter.GET("/comment/list/", controller.CommentList)

	// extra apis - II
	apiRouter.POST("/relation/action/", controller.RelationAction)
	apiRouter.GET("/relation/follow/list/", controller.FollowList)
	apiRouter.GET("/relation/follower/list/", controller.FollowerList)
}
