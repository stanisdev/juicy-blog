package models

type NewArticlesSubscriber struct {
  ID uint
  SubscriberID uint `gorm:"not null"valid:"required"`
  ArticleID uint `gorm:"not null"valid:"required"`
}

func (sm *DatabaseStaticMethods) AddNotifications(subscriberIds []uint, articleId uint) {
  for _, subscriberId := range subscriberIds {
    sm.DB.Create(&NewArticlesSubscriber {
      SubscriberID: subscriberId,
      ArticleID: articleId,
    })
  }
}
