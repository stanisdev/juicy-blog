package models

import (
  "github.com/jinzhu/gorm"
)

type SessionCookieKey struct {
  gorm.Model
  CookieName string `gorm:"size:40;not null;index:idx_cookie_name"`
  SessionDatas []SessionData `gorm:"ForeignKey:SessionCookieKeyID"`
}
