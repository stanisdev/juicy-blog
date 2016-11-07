package models

import (
  "crypto/sha1"
  "encoding/hex"
  "github.com/jinzhu/gorm"
)

type User struct {
  gorm.Model
  Name string `gorm:"size:20;unique;not null"`
  Email string `gorm:"size:50;not null"`
  Password string `gorm:"size:40;not null"`
  Articles []Article
}

func (u *User) ComparePassword(password string) bool {
  h := sha1.New()
  h.Write([]byte(password))
  return hex.EncodeToString(h.Sum(nil)) == u.Password
}