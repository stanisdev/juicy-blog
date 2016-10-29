package models

import (
  "github.com/jinzhu/gorm"
)

type SessionCookieKey struct {
  gorm.Model
  CookieName string
  SessionDatas []SessionData
}
