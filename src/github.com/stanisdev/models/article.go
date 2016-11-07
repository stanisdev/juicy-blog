package models

import (
  "github.com/jinzhu/gorm"
)

type Article struct {
  gorm.Model
  Title string `gorm:"not null"`
  Content string `gorm:"size:3000;not null"`
  UserID uint `gorm:"not null"`
}
