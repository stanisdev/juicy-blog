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

type ArticleValidator struct {
  Title string `valid:"required,length(5|255)"`
  Content string `valid:"required,length(1|3000)"`
  UserID uint
}