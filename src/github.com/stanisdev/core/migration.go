package core

import (
  m "github.com/stanisdev/models"
  "github.com/jinzhu/gorm"
  "github.com/stanisdev/db"
)

func DatabaseMigrate() {
  var con *gorm.DB = db.Connect();
  con.AutoMigrate(&m.SessionCookieKey{}, &m.SessionData{}, &m.User{})
  con.Model(&m.SessionData{}).AddForeignKey("session_cookie_key_id", "session_cookie_keys(id)", "CASCADE", "CASCADE")
}