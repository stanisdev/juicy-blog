package core

import (
  m "github.com/stanisdev/models"
  "github.com/stanisdev/db"
)

func SessionStart()  {
  con := db.GetConnection()
  con.AutoMigrate(&m.SessionCookieKey{}, &m.SessionData{})
  con.Model(&m.SessionData{}).AddForeignKey("session_cookie_key_id", "session_cookie_keys(id)", "CASCADE", "CASCADE")

  session := m.SessionCookieKey{
    CookieName: "cookie-uniq-id",
    SessionDatas: []m.SessionData{{Key: "name", Value: "John"}},
  }
  con.Create(&session)
}

func Set()  {

}

// func (self *Manager) CreateSession(cookieName string, data string)  {
//   connection.Create(&Session{CookieName: cookieName, Data: data})
// }
