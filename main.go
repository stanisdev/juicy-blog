/*
 * Mutual Blog
 * Copyright(c) 2016 Stanislav Zavalishin <javascript.nodejs.developer@gmail.com>
 * MIT Licensed
 */
 
package main

import (
  "github.com/stanisdev/core"
  "github.com/stanisdev/db"
) 

func main() {
  router := core.Router{Handlers: make(map[string]map[string]core.RouterHandler)}
  router.Config = core.GetConfig()
  db.DatabaseMigrate(router.Config.DbUser, router.Config.DbPass, router.Config.DbName)
  //db.ImportFixtures(router.Config.DbUser, router.Config.DbPass, router.Config.DbName)
  
  router.GET("/", core.Index)
  router.GET("/login", core.Login)
  router.POST("/login", core.LoginPost)
  router.GET("/logout", core.Logout)
  router.GET("/articles/:page?", core.Articles)
  router.GET("/articles/new", core.ArticleNew)
  router.POST("/articles/new", core.ArticleNewPost)
  router.GET("/article/:id", core.ArticleView)
  router.GET("/article/:id/edit", core.ArticleEdit)
  router.POST("/article/:id/edit", core.ArticleEditPost)
  router.POST("/article/:id/remove", core.ArticleRemovePost)
  router.GET("/profile/:id", core.ProfileView)
  router.Start()
}
