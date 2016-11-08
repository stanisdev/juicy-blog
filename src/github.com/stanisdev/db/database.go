package db

import (
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/mysql"
)

func Connect(user string, password string, dbName string) *gorm.DB  {
  con, err := gorm.Open("mysql", user + ":" + password + "@/" + dbName + "?charset=utf8&parseTime=True&loc=Local")
  if err != nil {
    panic("failed to connect database")
  }
  con.LogMode(true)
  return con
}