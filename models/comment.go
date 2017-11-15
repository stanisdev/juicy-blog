package models

import (
  "github.com/jinzhu/gorm"
  "time"
)

type Comment struct {
  gorm.Model
  Content string `gorm:"size:1000;not null"valid:"required,length(1|1000)"`
  ArticleID uint `gorm:"not null"valid:"required"`
  UserID uint `gorm:"not null"valid:"required"`
}

/**
 * Get comments by articles
 */
func (sm *DatabaseStaticMethods) GetComments(articleId uint, limit int) interface{} {
  var comments []struct {
    Id int;
    Content string;
    CreatedAt time.Time;
    UserId int;
    UserName string;
  }

  sm.DB.Table("comments c").
    Select("c.id, c.content, c.created_at, u.id as user_id, u.name user_name").
    Joins("LEFT JOIN users u ON c.user_id = u.id").
    Where("c.article_id = ?", articleId).
    Limit(limit).
    Scan(&comments)

  return comments
}

/**
 * Get comments total count by article's ID
 */
func (sm *DatabaseStaticMethods) GetCommentsCount(articleId uint) int {
  var commentsCount int
  sm.DB.Model(Comment{}).Where("article_id = ?", articleId).Count(&commentsCount)

  return commentsCount
}
