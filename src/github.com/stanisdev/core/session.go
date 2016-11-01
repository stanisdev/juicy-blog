package core

import (
  m "github.com/stanisdev/models"
  "github.com/jinzhu/gorm"
  "net/http"
  "time"
  "fmt"
)

type SessionManager struct {
  DB *gorm.DB
  Sid string
}

func (s *SessionManager) Start(w http.ResponseWriter, r *http.Request) {
  // s.DB.AutoMigrate(&m.SessionCookieKey{}, &m.SessionData{})
  // s.DB.Model(&m.SessionData{}).AddForeignKey("session_cookie_key_id", "session_cookie_keys(id)", "CASCADE", "CASCADE")
  
  cookie, _ := r.Cookie("sid")
  if len(cookie.String()) < 1 { // Set sid
    sid := GenerateRandomString(40)
    expiration := time.Now().Add(365 * 24 * time.Hour)
    cookie := http.Cookie{Name: "sid", Value: sid, Expires: expiration}
    http.SetCookie(w, &cookie)

    // Save sid to DB
    session := m.SessionCookieKey{
      CookieName: sid,
    }
    s.DB.Create(&session)
    s.Sid = sid
    // @TODO: clean outdated cookie keys from DB  
  } else {
    s.Sid = cookie.Value
  }
}

func (s *SessionManager) Set(key string, value string) {
  var data m.SessionCookieKey
  s.DB.Select("id").First(&data, "cookie_name = ?", s.Sid)
  fmt.Println(data.ID)
}