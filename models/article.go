package models

import (
  "github.com/jinzhu/gorm"
  "time"
)

type Article struct {
  gorm.Model
  Title string `gorm:"not null"valid:"required,length(5|255)"`
  Content string `gorm:"size:3000;not null"valid:"required,length(1|3000)"`
  UserID uint `gorm:"not null"valid:"required"`
}

/**
 * Get articles and related info by filter
 */
func (sm *DatabaseStaticMethods) GetArticles(params map[string]int, filter string) map[string]interface{} {

  var articles []struct{
    Id int;
    Title string;
    Content string;
    CreatedAt time.Time;
    Userid int;
    Username string;
    IsNew int
  }
  var count int
  limit := params["limit"]
  offset := params["offset"]
  currUserId := params["currUserId"]
  result := make(map[string]interface{})
  var caption string
  var newArticleIds []int

  switch filter {
    case "subscribed":
      // Get articles
      sm.DB.Table("articles a").
        Select(`
          a.id,
          a.title,
          SUBSTR(a.content, 1, 190) AS content,
          a.created_at,
          u.name username,
          u.id userid,
          IF (nus.id IS NULL, 0, 1 ) AS is_new`).
        Joins("LEFT JOIN subscribers sb ON a.user_id = sb.user_id").
        Joins("LEFT JOIN users u ON a.user_id = u.id").
        Joins("LEFT JOIN new_articles_subscribers nus ON a.id = nus.article_id AND nus.subscriber_id = ?", currUserId).
        Where("sb.subscriber_id = ? AND a.user_id != ?", currUserId, currUserId).
        Order("a.updated_at, a.created_at DESC").
        Limit(limit).
        Offset(offset).
        Scan(&articles)

      // Count articles
      sm.DB.Table("articles a").
        Joins("LEFT JOIN subscribers sb ON a.user_id = sb.user_id").
        Where("sb.subscriber_id = ? AND a.user_id != ?", currUserId, currUserId).
        Count(&count)

      caption = "subscribed"
    case "common":

      // Get articles
      sm.DB.Table("articles a").
        Select(`
          a.id,
          a.title,
          SUBSTR(a.content, 1, 190) AS content,
          a.created_at,
          u.name username,
          u.id userid,
          IF (nus.id IS NULL, 0, 1 ) AS is_new`).
        Joins("LEFT JOIN users u ON a.user_id = u.id").
        Joins("LEFT JOIN new_articles_subscribers nus ON a.id = nus.article_id AND nus.subscriber_id = ?", currUserId).
        Order("a.created_at DESC").
        Limit(limit).
        Offset(offset).
        Scan(&articles)

      // Count articles
      sm.DB.Model(Article{}).Count(&count)
      caption = "all"
  }
  for _, value := range articles {
    if value.IsNew > 0 {
      newArticleIds = append(newArticleIds, value.Id)
    }
  }
  // Remove notifications
  if len(newArticleIds) > 0 {
    sm.DB.Unscoped().Delete(NewArticlesSubscriber{}, "article_id in (?) AND subscriber_id = ?", newArticleIds, currUserId)
  }

  result["articles"] = articles
  result["count"] = count
  result["caption"] = caption
  return result
}

/**
 * Find article by id
 */
func (sm *DatabaseStaticMethods) FindArticleById(article interface{}, id int) {
  sm.DB.Table("articles a").
    Select("a.id, a.title, a.content, a.created_at, u.name username, u.id userid").
    Joins("LEFT JOIN users u ON a.user_id = u.id").
    Where("a.id = ?", id).
    Limit(1).
    Scan(article)
}
