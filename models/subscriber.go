package models

type Subscriber struct {
  ID uint
  UserID uint `gorm:"not null"valid:"required"`
  SubscriberID uint `gorm:"not null"valid:"required"`
}

func (sm *DatabaseStaticMethods) SubscriberExists(userId uint, subscriberId uint) bool {
  var record Subscriber
  return !sm.DB.Where("user_id = ? AND subscriber_id = ?", userId, subscriberId).First(&record).RecordNotFound()
}
