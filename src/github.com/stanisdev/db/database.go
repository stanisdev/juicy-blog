package db

import (
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/mysql"
)

type manager struct {
  Connection *gorm.DB
}

var instance *manager

func Connect()  {
  if instance == nil {
    con, err := gorm.Open("mysql", "root:root@/gorm?charset=utf8&parseTime=True&loc=Local")
    if err != nil {
      panic("failed to connect database")
    }
    instance = &manager{Connection: con}
  }
}

func GetConnection() *gorm.DB {
  return instance.Connection
}
