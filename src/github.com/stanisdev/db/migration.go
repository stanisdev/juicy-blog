package db

import (
  m "github.com/stanisdev/models"
  "github.com/jinzhu/gorm"
)

func DatabaseMigrate(config ...string) {
  var con *gorm.DB = Connect(config[0], config[1], config[2]);
  con.AutoMigrate(&m.SessionCookieKey{}, &m.SessionData{}, &m.User{}, &m.Article{})
  con.Model(&m.SessionData{}).AddForeignKey("session_cookie_key_id", "session_cookie_keys(id)", "CASCADE", "CASCADE")
  con.Model(&m.Article{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
}