package db

import (
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/mysql"
)

func Connect() *gorm.DB  {
  con, err := gorm.Open("mysql", "root:root@/gorm?charset=utf8&parseTime=True&loc=Local")
  if err != nil {
    panic("failed to connect database")
  }
  con.LogMode(true)
  return con
}
