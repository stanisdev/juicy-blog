package models

import (
  "github.com/jinzhu/gorm"
)

type Article struct {
  gorm.Model
  Title string `gorm:"not null"valid:"required,length(5|255)"`
  Content string `gorm:"size:3000;not null"valid:"required,length(1|3000)"`
  UserID uint
}

func (sm *StaticMethods) GetArticles(limit int, offset int) interface{} {
  var articles []struct{Id int; Title string; Content string; Userid int; Username string}
  sm.DB.Table("articles a").
    Select("a.id, a.title, a.content, u.name username, u.id userid").
    Joins("LEFT JOIN users u on a.user_id = u.id").
    Limit(5).
    Offset(2).
    Scan(&articles)
  return &articles
}