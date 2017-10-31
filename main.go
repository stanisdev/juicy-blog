/*
 * Mutual Blog
 * Copyright(c) 2016 Stanislav Zavalishin <javascript.nodejs.developer@gmail.com>
 * MIT Licensed
 */

package main

import (
  "github.com/stanisdev/juicy-blog/services"
  "github.com/stanisdev/juicy-blog/handlers"
)

func main()  {
  router := services.Router { Handlers: make(map[string] map[string] services.RouterHandler) }
  router.Init()

  router.GET("/", handlers.Index)
  router.GET("/login", handlers.Login)
  router.POST("/login", handlers.LoginPost)
  router.GET("/logout", handlers.Logout)
  router.GET("/articles/:page?", handlers.Articles)
  router.GET("/articles/new", handlers.ArticleNew)
  router.POST("/articles/new", handlers.ArticleNewPost)
  router.GET("/article/:id", handlers.ArticleView)
  router.GET("/article/:id/edit", handlers.ArticleEdit)
  router.POST("/article/:id/edit", handlers.ArticleEditPost)
  router.POST("/article/:id/remove", handlers.ArticleRemovePost)
  router.GET("/profile/:id", handlers.ProfileView)
  router.Run()
}
