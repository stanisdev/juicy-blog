package db

import (
  "github.com/stanisdev/models"
  "github.com/jinzhu/gorm"
)

func DatabaseMigrate(config ...string) {
  var con *gorm.DB = Connect(config[0], config[1], config[2]);
  con.AutoMigrate(&models.User{}, &models.Article{})
  con.Model(&models.Article{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
}