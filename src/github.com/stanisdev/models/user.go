package models

import (
  "github.com/jinzhu/gorm"
)

type User struct {
  gorm.Model
  Email string `gorm:"size:50"`
  Password string `gorm:"size:40"`
}

func (u *User) GetPassword() string {
  return u.Password
}