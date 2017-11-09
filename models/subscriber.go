package models

import (

)

type Subscriber struct {
  ID uint
  UserID uint `gorm:"not null"valid:"required"`
  SubscriberID uint `gorm:"not null"valid:"required"`
}

func (sm *DatabaseStaticMethods) SubscriberExists(userId uint, subscriberId uint) bool {
  var record Subscriber
  return !sm.DB.Where("user_id = ? AND subscriber_id = ?", userId, subscriberId).First(&record).RecordNotFound()
}

func (sm *DatabaseStaticMethods) FindAllSubscriberIds(userId uint) []uint {
  var ids []uint
  var subscribers []Subscriber
  sm.DB.Select("subscriber_id").Where("user_id = ?", userId).Find(&subscribers)

  for _, subscriber := range subscribers {
    ids = append(ids, subscriber.SubscriberID)
  }
  return ids
}
