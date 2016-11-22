package models

import (
  "github.com/jinzhu/gorm"
  "strings"
  "time"
)

type Article struct {
  gorm.Model
  Title string `gorm:"not null"valid:"required,length(5|255)"`
  Content string `gorm:"size:3000;not null"valid:"required,length(1|3000)"`
  UserID uint `gorm:"not null"valid:"required"`
}

/**
 * Get list of articles
 */
func (sm *StaticMethods) GetArticles(limit int, offset int) interface{} {
  var articles []struct{Id int; Title string; Content string; CreatedAt time.Time; Userid int; Username string}
  sm.DB.Table("articles a").
    Select("a.id, a.title, SUBSTR(a.content, 1, 190) AS content, a.created_at, u.name username, u.id userid").
    Joins("LEFT JOIN users u ON a.user_id = u.id").
    Order("a.created_at desc").
    Limit(limit).
    Offset(offset).
    Scan(&articles)
  for key, article := range articles {
    if len(article.Content) >= 190 {
      article.Content = article.Content[:strings.LastIndex(article.Content, " ")] + "..."
    }
    articles[key] = article
  }
  return &articles
}

/**
 * Find article by id
 */
func (sm *StaticMethods) FindArticleById(id int) (interface{}, int, *string) {
  var article struct{ID int; Title string; Content string; CreatedAt time.Time; Userid int; Username string}
  sm.DB.Table("articles a").
    Select("a.id, a.title, a.content, a.created_at, u.name username, u.id userid").
    Joins("LEFT JOIN users u ON a.user_id = u.id").
    Where("a.id = ?", id).
    Limit(1).
    Scan(&article)
  return &article, article.ID, &article.Title
}