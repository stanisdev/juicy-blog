/*
 * Mutual Blog
 * Copyright(c) 2016 Stanislav Zavalishin <javascript.nodejs.developer@gmail.com>
 * MIT Licensed
 */

package main

import (
  "github.com/stanisdev/juicy-blog/services"
  "github.com/stanisdev/juicy-blog/handlers"
  "github.com/stanisdev/juicy-blog/middlewares"
)

func main()  {
  router := services.Router {
    Handlers: make(map[string] map[string] services.RouterData),
  }
  router.Init()

  router.GET("/", handlers.Index)
  router.GET("/login", handlers.Login)
  router.POST("/login", handlers.LoginPost)
  router.GET("/logout", handlers.Logout, middlewares.Auth)
  router.GET("/articles/:page?", handlers.Articles)
  router.GET("/articles/new", handlers.ArticleNew, middlewares.Auth)
  router.POST("/articles/new", handlers.ArticleNewPost, middlewares.Auth)
  router.GET("/article/:id", handlers.ArticleView)
  router.GET("/article/:id/edit", handlers.ArticleEdit)
  router.POST("/article/:id/edit", handlers.ArticleEditPost)
  router.POST("/article/:id/remove", handlers.ArticleRemovePost)
  router.GET("/user/settings", handlers.UserSettings, middlewares.Auth)
  router.POST("/user/settings/save", handlers.UserSettingsSave, middlewares.Auth)
  router.POST("/user/password/change", handlers.UserPasswordChange, middlewares.Auth)
  router.GET("/user/:id", handlers.UserView)
  router.GET("/user/:id/subscribing", handlers.UserSubscribing)
  router.Run()
}
