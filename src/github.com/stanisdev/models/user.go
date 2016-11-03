package models

import (
  "crypto/sha1"
  "encoding/hex"
  "github.com/jinzhu/gorm"
)

type User struct {
  gorm.Model
  Email string `gorm:"size:50"`
  Password string `gorm:"size:40"`
}

func (u *User) ComparePassword(password string) bool {
  h := sha1.New()
  h.Write([]byte(password))
  return hex.EncodeToString(h.Sum(nil)) == u.Password
}