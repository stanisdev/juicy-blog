package models

import (
  "github.com/jinzhu/gorm"
)

type SessionData struct {
  gorm.Model
  Key string
  Value string
  SessionCookieKeyID uint
}
