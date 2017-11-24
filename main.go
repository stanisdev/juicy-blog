/*
 * Mutual Blog
 * Copyright(c) 2016 Stanislav Zavalishin <javascript.nodejs.developer@gmail.com>
 * MIT Licensed
 */

package main

import (
  "github.com/stanisdev/juicy-blog/services"
  "github.com/stanisdev/juicy-blog/handlers"
  mw "github.com/stanisdev/juicy-blog/middlewares"
)

func main()  {
  router := services.Router {
    Handlers: make(map[string] map[string] services.RouterData),
  }
  router.Init()
  IdIntMiddleware := mw.ValidateUrl(map[string]string{ "id": "int" })

  router.GET("/", handlers.Index)
  router.GET("/login", handlers.Login)
  router.POST("/login", handlers.LoginPost)
  router.GET("/logout", handlers.Logout, mw.Auth)
  router.GET("/articles/:page?", handlers.Articles)
  router.GET("/articles/new", handlers.ArticleNew, mw.Auth)
  router.POST("/articles/new", handlers.ArticleNewPost, mw.Auth)
  router.GET("/article/:id", handlers.ArticleView, IdIntMiddleware, mw.GetAritcle)

  router.GET("/article/:id/edit",
    handlers.ArticleEdit,
    mw.Auth,
    IdIntMiddleware,
    mw.GetAritcle,
  )
  router.POST("/article/:id/edit", handlers.ArticleEditPost, mw.Auth, IdIntMiddleware, mw.GetAritcle)
  router.POST("/article/:id/remove", handlers.ArticleRemovePost, mw.Auth, IdIntMiddleware, mw.GetAritcle)

  router.GET("/user/settings", handlers.UserSettings, mw.Auth)
  router.POST("/user/settings/save", handlers.UserSettingsSave, mw.Auth)
  router.POST("/user/password/change", handlers.UserPasswordChange, mw.Auth)
  router.GET("/user/:id", handlers.UserView, IdIntMiddleware)
  router.GET("/user/:id/subscribing", handlers.UserSubscribing, mw.Auth, IdIntMiddleware)
  router.POST("/article/:id/comment", handlers.ArticleAddCommentPost, mw.Auth, IdIntMiddleware, mw.GetAritcle)
  router.Run()
}
